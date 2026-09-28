package store

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"agentgo/p_backend/internal/domain"

	_ "modernc.org/sqlite"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := s.RecoverInterrupted(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS tasks (
  id TEXT PRIMARY KEY,
  access_token_hash TEXT NOT NULL,
  status TEXT NOT NULL,
  stage TEXT NOT NULL,
  dialect TEXT NOT NULL,
  host TEXT NOT NULL,
  port INTEGER NOT NULL,
  database_name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  pause_requested_at TEXT,
  paused_at TEXT,
  connection_required INTEGER NOT NULL DEFAULT 0,
  failure_summary TEXT NOT NULL DEFAULT '',
  redact_object_names INTEGER NOT NULL DEFAULT 0,
  schema_json TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS workers (
  task_id TEXT NOT NULL,
  worker_id TEXT NOT NULL,
  name TEXT NOT NULL,
  status TEXT NOT NULL,
  total_units INTEGER NOT NULL DEFAULT 0,
  completed_units INTEGER NOT NULL DEFAULT 0,
  current_object TEXT NOT NULL DEFAULT '',
  finding_count INTEGER NOT NULL DEFAULT 0,
  input_tokens INTEGER NOT NULL DEFAULT 0,
  output_tokens INTEGER NOT NULL DEFAULT 0,
  total_tokens INTEGER NOT NULL DEFAULT 0,
  checkpoint TEXT NOT NULL DEFAULT '',
  error_summary TEXT NOT NULL DEFAULT '',
  started_at TEXT,
  updated_at TEXT NOT NULL,
  finished_at TEXT,
  PRIMARY KEY (task_id, worker_id),
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS work_units (
  task_id TEXT NOT NULL,
  worker_id TEXT NOT NULL,
  unit_key TEXT NOT NULL,
  status TEXT NOT NULL,
  checkpoint TEXT NOT NULL DEFAULT '',
  attempts INTEGER NOT NULL DEFAULT 0,
  lease_until TEXT,
  completed_at TEXT,
  PRIMARY KEY (task_id, worker_id, unit_key),
  FOREIGN KEY (task_id, worker_id) REFERENCES workers(task_id, worker_id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS findings (
  task_id TEXT NOT NULL,
  id TEXT NOT NULL,
  payload_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  PRIMARY KEY (task_id, id),
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS reports (
  task_id TEXT PRIMARY KEY,
  payload_json TEXT NOT NULL,
  markdown TEXT NOT NULL,
  html TEXT NOT NULL,
  generated_at TEXT NOT NULL,
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS report_rebuilds (
  task_id TEXT PRIMARY KEY,
  state TEXT NOT NULL,
  requested_at TEXT NOT NULL,
  started_at TEXT,
  finished_at TEXT,
  failure_summary TEXT NOT NULL DEFAULT '',
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS audit_events (
  id TEXT PRIMARY KEY,
  task_id TEXT NOT NULL DEFAULT '',
  action TEXT NOT NULL,
  outcome TEXT NOT NULL,
  detail TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS report_baselines (
  name TEXT PRIMARY KEY,
  task_id TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_work_units_pending ON work_units(task_id, worker_id, status);
CREATE INDEX IF NOT EXISTS idx_findings_task ON findings(task_id, created_at);`)
	if err != nil {
		return err
	}
	if err := s.migrateFindingsPrimaryKey(ctx); err != nil {
		return err
	}
	if err := s.migrateWorkerTokenColumns(ctx); err != nil {
		return err
	}
	return s.migrateDMLWorkerToDDL(ctx)
}

func (s *Store) migrateDMLWorkerToDDL(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	if _, err = tx.ExecContext(ctx, `INSERT INTO workers (
 task_id, worker_id, name, status, total_units, completed_units, current_object, finding_count,
 input_tokens, output_tokens, total_tokens, checkpoint, error_summary, started_at, updated_at, finished_at
)
SELECT task_id, 'ddl-generator', 'DDL 建议生成(AI参与)', status, total_units, completed_units, current_object, finding_count,
 input_tokens, output_tokens, total_tokens, checkpoint, error_summary, started_at, updated_at, finished_at
FROM workers old
WHERE old.worker_id='dml-generator'
  AND NOT EXISTS (SELECT 1 FROM workers current WHERE current.task_id=old.task_id AND current.worker_id='ddl-generator')`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE work_units SET worker_id='ddl-generator' WHERE worker_id='dml-generator'`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM workers WHERE worker_id='dml-generator'`); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) migrateWorkerTokenColumns(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `PRAGMA table_info(workers)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	columns := map[string]bool{}
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, dataType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			return err
		}
		columns[name] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, column := range []string{"input_tokens", "output_tokens", "total_tokens"} {
		if columns[column] {
			continue
		}
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE workers ADD COLUMN `+column+` INTEGER NOT NULL DEFAULT 0`); err != nil {
			return err
		}
	}
	return nil
}

// Earlier versions used findings.id as a global primary key. Rule finding IDs
// are deterministic, so identical checks in a later task were silently ignored
// while the Worker counter still increased. Findings must be unique per task.
func (s *Store) migrateFindingsPrimaryKey(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `PRAGMA table_info(findings)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	primaryKey := make(map[int]string)
	for rows.Next() {
		var cid, notNull, position int
		var name, dataType string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &position); err != nil {
			return err
		}
		if position > 0 {
			primaryKey[position] = name
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if primaryKey[1] == "task_id" && primaryKey[2] == "id" && len(primaryKey) == 2 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	if _, err = tx.ExecContext(ctx, `CREATE TABLE findings_next (
  task_id TEXT NOT NULL,
  id TEXT NOT NULL,
  payload_json TEXT NOT NULL,
  created_at TEXT NOT NULL,
  PRIMARY KEY (task_id, id),
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
)`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO findings_next (task_id, id, payload_json, created_at) SELECT task_id, id, payload_json, created_at FROM findings`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DROP TABLE findings`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `ALTER TABLE findings_next RENAME TO findings`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `CREATE INDEX idx_findings_task ON findings(task_id, created_at)`); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) RecoverInterrupted(ctx context.Context) error {
	now := stamp(time.Now().UTC())
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	if _, err = tx.ExecContext(ctx, `UPDATE work_units SET status = 'pending', lease_until = NULL WHERE status = 'running'`); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE workers SET status = 'paused', updated_at = ? WHERE status IN ('running', 'pausing')`, now); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE tasks SET status = 'paused', connection_required = 1, failure_summary = '服务重启后需要重新提交连接信息', updated_at = ? WHERE status IN ('running', 'pausing')`, now); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE report_rebuilds SET state='failed', finished_at=?, failure_summary='服务重启导致报告重建中断，请重新生成报告' WHERE state IN ('pending', 'running')`, now); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) RecoverFailedWork(ctx context.Context, taskID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	if _, err = tx.ExecContext(ctx, `UPDATE work_units SET status='pending', lease_until=NULL WHERE task_id=? AND status='running'`, taskID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE workers SET status='pending', current_object='', checkpoint='', error_summary='', finished_at=NULL, updated_at=? WHERE task_id=? AND status='failed'`, stamp(time.Now().UTC()), taskID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) ListPendingWorkUnitKeys(ctx context.Context, taskID, workerID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT unit_key FROM work_units WHERE task_id=? AND worker_id=? AND status='pending' ORDER BY unit_key`, taskID, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := make([]string, 0)
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

func (s *Store) ReplacePendingWorkUnits(ctx context.Context, taskID, workerID string, keys []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	if _, err = tx.ExecContext(ctx, `DELETE FROM work_units WHERE task_id=? AND worker_id=? AND status='pending'`, taskID, workerID); err != nil {
		return err
	}
	for _, key := range keys {
		if _, err = tx.ExecContext(ctx, `INSERT INTO work_units (task_id, worker_id, unit_key, status) VALUES (?, ?, ?, 'pending')`, taskID, workerID, key); err != nil {
			return err
		}
	}
	var completed, total int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FILTER (WHERE status='completed'), COUNT(*) FROM work_units WHERE task_id=? AND worker_id=?`, taskID, workerID).Scan(&completed, &total); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE workers SET completed_units=?, total_units=?, current_object='', checkpoint='', updated_at=? WHERE task_id=? AND worker_id=?`, completed, total, stamp(time.Now().UTC()), taskID, workerID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) CreateTask(ctx context.Context, task domain.Task, workers []domain.WorkerProgress) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	_, err = tx.ExecContext(ctx, `INSERT INTO tasks (
 id, access_token_hash, status, stage, dialect, host, port, database_name, created_at, updated_at, expires_at, redact_object_names
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.AccessTokenHash, task.Status, task.Stage, task.Target.Dialect, task.Target.Host, task.Target.Port, task.Target.Database,
		stamp(task.CreatedAt), stamp(task.UpdatedAt), stamp(task.ExpiresAt), boolInt(task.RedactObjectNames))
	if err != nil {
		return err
	}
	for _, worker := range workers {
		if _, err = tx.ExecContext(ctx, `INSERT INTO workers (task_id, worker_id, name, status, updated_at) VALUES (?, ?, ?, ?, ?)`,
			task.ID, worker.WorkerID, worker.Name, worker.Status, stamp(task.UpdatedAt)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) Authorize(ctx context.Context, taskID, token string) error {
	var hash, expires string
	err := s.db.QueryRowContext(ctx, `SELECT access_token_hash, expires_at FROM tasks WHERE id = ?`, taskID).Scan(&hash, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	expiresAt, err := parseStamp(expires)
	if err != nil || time.Now().UTC().After(expiresAt) || token == "" || subtle.ConstantTimeCompare([]byte(hash), []byte(tokenHash(token))) != 1 {
		return ErrUnauthorized
	}
	return nil
}

func (s *Store) GetTask(ctx context.Context, id string) (domain.Task, error) {
	var t domain.Task
	var created, updated, expires string
	var pauseRequested, paused sql.NullString
	var connectionRequired, redact int
	err := s.db.QueryRowContext(ctx, `SELECT id, access_token_hash, status, stage, dialect, host, port, database_name,
 created_at, updated_at, expires_at, pause_requested_at, paused_at, connection_required, failure_summary, redact_object_names
 FROM tasks WHERE id = ?`, id).Scan(&t.ID, &t.AccessTokenHash, &t.Status, &t.Stage, &t.Target.Dialect, &t.Target.Host, &t.Target.Port, &t.Target.Database,
		&created, &updated, &expires, &pauseRequested, &paused, &connectionRequired, &t.FailureSummary, &redact)
	if errors.Is(err, sql.ErrNoRows) {
		return t, ErrNotFound
	}
	if err != nil {
		return t, err
	}
	var parseErr error
	if t.CreatedAt, parseErr = parseStamp(created); parseErr != nil {
		return t, parseErr
	}
	if t.UpdatedAt, parseErr = parseStamp(updated); parseErr != nil {
		return t, parseErr
	}
	if t.ExpiresAt, parseErr = parseStamp(expires); parseErr != nil {
		return t, parseErr
	}
	if pauseRequested.Valid {
		if value, err := parseStamp(pauseRequested.String); err == nil {
			t.PauseRequestedAt = &value
		}
	}
	if paused.Valid {
		if value, err := parseStamp(paused.String); err == nil {
			t.PausedAt = &value
		}
	}
	t.ConnectionRequired = connectionRequired == 1
	t.RedactObjectNames = redact == 1
	workers, err := s.listWorkers(ctx, id)
	if err != nil {
		return t, err
	}
	t.Progress.Workers = workers
	t.Progress.WorkerFlow = domain.DefaultWorkerFlow()
	t.Progress.Stage = t.Stage
	t.Progress.UpdatedAt = t.UpdatedAt
	for _, worker := range workers {
		t.Progress.TotalUnits += worker.TotalUnits
		t.Progress.CompletedUnits += worker.CompletedUnits
		t.Progress.TokenUsage = t.Progress.TokenUsage.Add(worker.TokenUsage)
		if worker.CurrentObject != "" && worker.Status == domain.WorkerRunning {
			t.Progress.CurrentObject = worker.CurrentObject
		}
	}
	if t.Progress.TotalUnits > 0 {
		t.Progress.Percent = t.Progress.CompletedUnits * 100 / t.Progress.TotalUnits
	}
	return t, nil
}

func (s *Store) SaveTask(ctx context.Context, t domain.Task) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tasks SET status=?, stage=?, updated_at=?, pause_requested_at=?, paused_at=?,
 connection_required=?, failure_summary=? WHERE id=?`, t.Status, t.Stage, stamp(t.UpdatedAt), nullableStamp(t.PauseRequestedAt), nullableStamp(t.PausedAt), boolInt(t.ConnectionRequired), t.FailureSummary, t.ID)
	return err
}

func (s *Store) SaveSnapshot(ctx context.Context, taskID string, snapshot domain.SchemaSnapshot) error {
	payload, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `UPDATE tasks SET schema_json=?, updated_at=? WHERE id=?`, string(payload), stamp(time.Now().UTC()), taskID)
	return err
}

func (s *Store) LoadSnapshot(ctx context.Context, taskID string) (domain.SchemaSnapshot, error) {
	var payload string
	err := s.db.QueryRowContext(ctx, `SELECT schema_json FROM tasks WHERE id=?`, taskID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.SchemaSnapshot{}, ErrNotFound
	}
	if err != nil {
		return domain.SchemaSnapshot{}, err
	}
	if payload == "" {
		return domain.SchemaSnapshot{}, ErrNotFound
	}
	var snapshot domain.SchemaSnapshot
	if err := json.Unmarshal([]byte(payload), &snapshot); err != nil {
		return domain.SchemaSnapshot{}, err
	}
	return snapshot, nil
}

func (s *Store) listWorkers(ctx context.Context, taskID string) ([]domain.WorkerProgress, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT worker_id, name, status, total_units, completed_units, current_object, finding_count,
 input_tokens, output_tokens, total_tokens, checkpoint, error_summary, started_at, updated_at, finished_at FROM workers WHERE task_id=? ORDER BY worker_id`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	workers := make([]domain.WorkerProgress, 0)
	for rows.Next() {
		var worker domain.WorkerProgress
		var started, finished sql.NullString
		var updated string
		if err := rows.Scan(&worker.WorkerID, &worker.Name, &worker.Status, &worker.TotalUnits, &worker.CompletedUnits, &worker.CurrentObject,
			&worker.FindingCount, &worker.TokenUsage.InputTokens, &worker.TokenUsage.OutputTokens, &worker.TokenUsage.TotalTokens, &worker.Checkpoint, &worker.ErrorSummary, &started, &updated, &finished); err != nil {
			return nil, err
		}
		if worker.UpdatedAt, err = parseStamp(updated); err != nil {
			return nil, err
		}
		if started.Valid {
			if value, err := parseStamp(started.String); err == nil {
				worker.StartedAt = &value
			}
		}
		if finished.Valid {
			if value, err := parseStamp(finished.String); err == nil {
				worker.FinishedAt = &value
			}
		}
		workers = append(workers, worker)
	}
	return workers, rows.Err()
}

func (s *Store) GetWorker(ctx context.Context, taskID, workerID string) (domain.WorkerProgress, error) {
	workers, err := s.listWorkers(ctx, taskID)
	if err != nil {
		return domain.WorkerProgress{}, err
	}
	for _, worker := range workers {
		if worker.WorkerID == workerID {
			return worker, nil
		}
	}
	return domain.WorkerProgress{}, ErrNotFound
}

func (s *Store) SetWorker(ctx context.Context, taskID string, worker domain.WorkerProgress) error {
	now := time.Now().UTC()
	if worker.UpdatedAt.IsZero() {
		worker.UpdatedAt = now
	}
	if worker.Status == domain.WorkerRunning && worker.StartedAt == nil {
		worker.StartedAt = &now
	}
	if (worker.Status == domain.WorkerCompleted || worker.Status == domain.WorkerFailed || worker.Status == domain.WorkerCancelled) && worker.FinishedAt == nil {
		worker.FinishedAt = &now
	}
	_, err := s.db.ExecContext(ctx, `UPDATE workers SET status=?, total_units=?, completed_units=?, current_object=?, finding_count=?, input_tokens=?, output_tokens=?, total_tokens=?, checkpoint=?, error_summary=?, started_at=?, updated_at=?, finished_at=? WHERE task_id=? AND worker_id=?`,
		worker.Status, worker.TotalUnits, worker.CompletedUnits, worker.CurrentObject, worker.FindingCount, worker.TokenUsage.InputTokens, worker.TokenUsage.OutputTokens, worker.TokenUsage.TotalTokens, worker.Checkpoint, worker.ErrorSummary,
		nullableStamp(worker.StartedAt), stamp(worker.UpdatedAt), nullableStamp(worker.FinishedAt), taskID, worker.WorkerID)
	return err
}

func (s *Store) AddWorkerTokenUsage(ctx context.Context, taskID, workerID string, usage domain.TokenUsage) error {
	if usage.InputTokens == 0 && usage.OutputTokens == 0 && usage.TotalTokens == 0 {
		return nil
	}
	_, err := s.db.ExecContext(ctx, `UPDATE workers SET input_tokens=input_tokens+?, output_tokens=output_tokens+?, total_tokens=total_tokens+?, updated_at=? WHERE task_id=? AND worker_id=?`, usage.InputTokens, usage.OutputTokens, usage.TotalTokens, stamp(time.Now().UTC()), taskID, workerID)
	return err
}

func (s *Store) SetAllRunningWorkers(ctx context.Context, taskID string, status domain.WorkerStatus) error {
	_, err := s.db.ExecContext(ctx, `UPDATE workers SET status=?, updated_at=? WHERE task_id=? AND status IN ('running','pausing')`, status, stamp(time.Now().UTC()), taskID)
	return err
}

func (s *Store) AddWorkUnits(ctx context.Context, taskID, workerID string, keys []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	for _, key := range keys {
		if _, err = tx.ExecContext(ctx, `INSERT OR IGNORE INTO work_units (task_id, worker_id, unit_key, status) VALUES (?, ?, ?, 'pending')`, taskID, workerID, key); err != nil {
			return err
		}
	}
	var total int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM work_units WHERE task_id=? AND worker_id=?`, taskID, workerID).Scan(&total); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE workers SET total_units=?, updated_at=? WHERE task_id=? AND worker_id=?`, total, stamp(time.Now().UTC()), taskID, workerID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) ClaimNextWorkUnit(ctx context.Context, taskID, workerID string) (string, bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", false, err
	}
	defer rollback(tx)
	var key string
	err = tx.QueryRowContext(ctx, `SELECT unit_key FROM work_units WHERE task_id=? AND worker_id=? AND status='pending' ORDER BY unit_key LIMIT 1`, taskID, workerID).Scan(&key)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, tx.Commit()
	}
	if err != nil {
		return "", false, err
	}
	leaseUntil := stamp(time.Now().UTC().Add(2 * time.Minute))
	if _, err = tx.ExecContext(ctx, `UPDATE work_units SET status='running', attempts=attempts+1, lease_until=? WHERE task_id=? AND worker_id=? AND unit_key=? AND status='pending'`, leaseUntil, taskID, workerID, key); err != nil {
		return "", false, err
	}
	if err = tx.Commit(); err != nil {
		return "", false, err
	}
	return key, true, nil
}

func (s *Store) ReturnWorkUnit(ctx context.Context, taskID, workerID, key string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE work_units SET status='pending', lease_until=NULL WHERE task_id=? AND worker_id=? AND unit_key=? AND status='running'`, taskID, workerID, key)
	return err
}

func (s *Store) CompleteWorkUnit(ctx context.Context, taskID, workerID, key string, findings []domain.Finding) error {
	return s.completeWorkUnit(ctx, taskID, workerID, key, key, findings)
}

func (s *Store) CompleteWorkUnitWithProgress(ctx context.Context, taskID, workerID, key, progressObject string, findings []domain.Finding) error {
	return s.completeWorkUnit(ctx, taskID, workerID, key, progressObject, findings)
}

func (s *Store) completeWorkUnit(ctx context.Context, taskID, workerID, key, progressObject string, findings []domain.Finding) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	for _, finding := range findings {
		payload, err := json.Marshal(finding)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO findings (id, task_id, payload_json, created_at) VALUES (?, ?, ?, ?)
ON CONFLICT(task_id, id) DO UPDATE SET payload_json=excluded.payload_json, created_at=excluded.created_at`, finding.ID, taskID, string(payload), stamp(finding.CreatedAt)); err != nil {
			return err
		}
	}
	result, err := tx.ExecContext(ctx, `UPDATE work_units SET status='completed', checkpoint=?, lease_until=NULL, completed_at=? WHERE task_id=? AND worker_id=? AND unit_key=? AND status IN ('running','pending')`, key, stamp(time.Now().UTC()), taskID, workerID, key)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return fmt.Errorf("work unit %s is not claimable", key)
	}
	var completed, total int
	if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FILTER (WHERE status='completed'), COUNT(*) FROM work_units WHERE task_id=? AND worker_id=?`, taskID, workerID).Scan(&completed, &total); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE workers SET completed_units=?, total_units=?, current_object=?, checkpoint=?, finding_count=finding_count+?, updated_at=? WHERE task_id=? AND worker_id=?`, completed, total, progressObject, key, len(findings), stamp(time.Now().UTC()), taskID, workerID); err != nil {
		return err
	}
	return tx.Commit()
}

// SaveFindings persists recovered findings without changing WorkUnit or Worker
// progress. It is used only to repair old reports affected by legacy storage.
func (s *Store) SaveFindings(ctx context.Context, taskID string, findings []domain.Finding) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer rollback(tx)
	for _, finding := range findings {
		payload, err := json.Marshal(finding)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO findings (id, task_id, payload_json, created_at) VALUES (?, ?, ?, ?)
ON CONFLICT(task_id, id) DO UPDATE SET payload_json=excluded.payload_json, created_at=excluded.created_at`, finding.ID, taskID, string(payload), stamp(finding.CreatedAt)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListFindings(ctx context.Context, taskID string) ([]domain.Finding, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT payload_json FROM findings WHERE task_id=? ORDER BY created_at, id`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	findings := make([]domain.Finding, 0)
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}
		var finding domain.Finding
		if err := json.Unmarshal([]byte(payload), &finding); err != nil {
			return nil, err
		}
		findings = append(findings, finding)
	}
	return findings, rows.Err()
}

func (s *Store) SaveReport(ctx context.Context, report domain.Report) error {
	payload, err := json.Marshal(report)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO reports (task_id, payload_json, markdown, html, generated_at) VALUES (?, ?, ?, ?, ?)
ON CONFLICT(task_id) DO UPDATE SET payload_json=excluded.payload_json, markdown=excluded.markdown, html=excluded.html, generated_at=excluded.generated_at`,
		report.TaskID, string(payload), report.Markdown, report.HTML, stamp(report.GeneratedAt))
	return err
}

func (s *Store) GetReport(ctx context.Context, taskID string) (domain.Report, error) {
	var payload, markdown, html string
	err := s.db.QueryRowContext(ctx, `SELECT payload_json, markdown, html FROM reports WHERE task_id=?`, taskID).Scan(&payload, &markdown, &html)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Report{}, ErrNotFound
	}
	if err != nil {
		return domain.Report{}, err
	}
	var report domain.Report
	if err := json.Unmarshal([]byte(payload), &report); err != nil {
		return domain.Report{}, err
	}
	report.Markdown, report.HTML = markdown, html
	return report, nil
}

// RequestReportRebuild persists a request and reports whether this caller must
// start a background worker. Pending/running requests are joined rather than
// starting duplicate LLM calls.
func (s *Store) RequestReportRebuild(ctx context.Context, taskID string) (domain.ReportRebuildStatus, bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ReportRebuildStatus{}, false, err
	}
	defer rollback(tx)
	status, exists, err := getReportRebuildStatus(ctx, tx, taskID)
	if err != nil {
		return domain.ReportRebuildStatus{}, false, err
	}
	if exists && (status.State == domain.ReportRebuildPending || status.State == domain.ReportRebuildRunning) {
		if err := tx.Commit(); err != nil {
			return domain.ReportRebuildStatus{}, false, err
		}
		return status, false, nil
	}
	now := time.Now().UTC()
	if _, err = tx.ExecContext(ctx, `INSERT INTO report_rebuilds (task_id, state, requested_at, started_at, finished_at, failure_summary) VALUES (?, 'pending', ?, NULL, NULL, '')
ON CONFLICT(task_id) DO UPDATE SET state='pending', requested_at=excluded.requested_at, started_at=NULL, finished_at=NULL, failure_summary=''`, taskID, stamp(now)); err != nil {
		return domain.ReportRebuildStatus{}, false, err
	}
	if err := tx.Commit(); err != nil {
		return domain.ReportRebuildStatus{}, false, err
	}
	return domain.ReportRebuildStatus{TaskID: taskID, State: domain.ReportRebuildPending, RequestedAt: &now}, true, nil
}

func (s *Store) GetReportRebuildStatus(ctx context.Context, taskID string) (domain.ReportRebuildStatus, error) {
	status, exists, err := getReportRebuildStatus(ctx, s.db, taskID)
	if err != nil {
		return domain.ReportRebuildStatus{}, err
	}
	if exists {
		return status, nil
	}
	// A completed inspection already has a generated report even before anyone
	// manually requests a rebuild. Surface that fact and its completion time to
	// the task-list UI instead of incorrectly showing "not generated".
	var generatedAt string
	err = s.db.QueryRowContext(ctx, `SELECT generated_at FROM reports WHERE task_id=?`, taskID).Scan(&generatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ReportRebuildStatus{TaskID: taskID, State: domain.ReportRebuildIdle}, nil
	}
	if err != nil {
		return domain.ReportRebuildStatus{}, err
	}
	generated, err := parseStamp(generatedAt)
	if err != nil {
		return domain.ReportRebuildStatus{}, err
	}
	return domain.ReportRebuildStatus{TaskID: taskID, State: domain.ReportRebuildSucceeded, FinishedAt: &generated}, nil
}

func (s *Store) MarkReportRebuildRunning(ctx context.Context, taskID string) error {
	now := time.Now().UTC()
	result, err := s.db.ExecContext(ctx, `UPDATE report_rebuilds SET state='running', started_at=?, finished_at=NULL, failure_summary='' WHERE task_id=? AND state='pending'`, stamp(now), taskID)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		return errors.New("report rebuild is not pending")
	}
	return nil
}

func (s *Store) CompleteReportRebuild(ctx context.Context, taskID string, failureSummary string) error {
	now := time.Now().UTC()
	state := domain.ReportRebuildSucceeded
	if strings.TrimSpace(failureSummary) != "" {
		state = domain.ReportRebuildFailed
	}
	_, err := s.db.ExecContext(ctx, `UPDATE report_rebuilds SET state=?, finished_at=?, failure_summary=? WHERE task_id=?`, state, stamp(now), failureSummary, taskID)
	return err
}

type reportRebuildStatusReader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func getReportRebuildStatus(ctx context.Context, reader reportRebuildStatusReader, taskID string) (domain.ReportRebuildStatus, bool, error) {
	status := domain.ReportRebuildStatus{TaskID: taskID}
	var requestedAt string
	var startedAt, finishedAt sql.NullString
	err := reader.QueryRowContext(ctx, `SELECT state, requested_at, started_at, finished_at, failure_summary FROM report_rebuilds WHERE task_id=?`, taskID).Scan(&status.State, &requestedAt, &startedAt, &finishedAt, &status.FailureSummary)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ReportRebuildStatus{}, false, nil
	}
	if err != nil {
		return domain.ReportRebuildStatus{}, false, err
	}
	parsedRequested, err := parseStamp(requestedAt)
	if err != nil {
		return domain.ReportRebuildStatus{}, false, err
	}
	status.RequestedAt = &parsedRequested
	if startedAt.Valid {
		parsedStarted, err := parseStamp(startedAt.String)
		if err != nil {
			return domain.ReportRebuildStatus{}, false, err
		}
		status.StartedAt = &parsedStarted
	}
	if finishedAt.Valid {
		parsedFinished, err := parseStamp(finishedAt.String)
		if err != nil {
			return domain.ReportRebuildStatus{}, false, err
		}
		status.FinishedAt = &parsedFinished
	}
	return status, true, nil
}

func (s *Store) RedactObjectNames(ctx context.Context, taskID string) (bool, error) {
	var value int
	err := s.db.QueryRowContext(ctx, `SELECT redact_object_names FROM tasks WHERE id=?`, taskID).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrNotFound
	}
	return value == 1, err
}

func (s *Store) AddAuditEvent(ctx context.Context, event domain.AuditEvent) error {
	if event.ID == "" {
		event.ID = fmt.Sprintf("audit-%d", time.Now().UnixNano())
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_events (id, task_id, action, outcome, detail, created_at) VALUES (?, ?, ?, ?, ?, ?)`, event.ID, event.TaskID, event.Action, event.Outcome, event.Detail, stamp(event.CreatedAt))
	return err
}

func (s *Store) ListAuditEvents(ctx context.Context, taskID string, limit int) ([]domain.AuditEvent, error) {
	if limit < 1 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, task_id, action, outcome, detail, created_at FROM audit_events WHERE task_id = ? ORDER BY created_at DESC LIMIT ?`, taskID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]domain.AuditEvent, 0)
	for rows.Next() {
		var event domain.AuditEvent
		var created string
		if err := rows.Scan(&event.ID, &event.TaskID, &event.Action, &event.Outcome, &event.Detail, &created); err != nil {
			return nil, err
		}
		if event.CreatedAt, err = parseStamp(created); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (s *Store) SetBaseline(ctx context.Context, name, taskID string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO report_baselines (name, task_id, updated_at) VALUES (?, ?, ?) ON CONFLICT(name) DO UPDATE SET task_id=excluded.task_id, updated_at=excluded.updated_at`, name, taskID, stamp(time.Now().UTC()))
	return err
}

func (s *Store) GetBaseline(ctx context.Context, name string) (string, error) {
	var taskID string
	err := s.db.QueryRowContext(ctx, `SELECT task_id FROM report_baselines WHERE name=?`, name).Scan(&taskID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return taskID, err
}

func rollback(tx *sql.Tx) { _ = tx.Rollback() }
func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
func stamp(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }
func nullableStamp(value *time.Time) any {
	if value == nil {
		return nil
	}
	return stamp(*value)
}
func parseStamp(value string) (time.Time, error) { return time.Parse(time.RFC3339Nano, value) }
func tokenHash(value string) string              { return fmt.Sprintf("%x", sha256sum(value)) }
func sha256sum(value string) [32]byte            { return sha256.Sum256([]byte(strings.TrimSpace(value))) }
