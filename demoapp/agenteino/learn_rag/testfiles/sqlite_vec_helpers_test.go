//go:build sqlitevec

package testfiles

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"testing"

	sqlitevec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"
)

var sqliteVecAutoOnce sync.Once

type sqliteVecStore struct {
	db  *sql.DB
	dim int
}

func newSQLiteVecStore(t *testing.T, path string, dim int) *sqliteVecStore {
	t.Helper()
	if dim != 3 {
		t.Fatalf("this learning table is intentionally fixed to 3 dimensions, got %d", dim)
	}
	sqliteVecAutoOnce.Do(sqlitevec.Auto)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	db.SetMaxOpenConns(1)
	store := &sqliteVecStore{db: db, dim: dim}
	if _, err := db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS rag_vectors USING vec0(
			chunk_id text primary key,
			tenant_id text partition key,
			document_id text,
			embedding float[3],
			+source_uri text,
			+chunk_text text
		)
	`); err != nil {
		db.Close()
		t.Fatalf("create vec0 table: %v", err)
	}
	return store
}

func (s *sqliteVecStore) Close() error {
	return s.db.Close()
}

func (s *sqliteVecStore) VecVersion() (string, error) {
	var version string
	err := s.db.QueryRow(`SELECT vec_version()`).Scan(&version)
	return version, err
}

// sqlite-vec 当前不将 INSERT OR REPLACE 视为向量主键的安全 upsert。
// 先删除再插入能让更新生命周期保持显式、可验证。
func (s *sqliteVecStore) Upsert(chunk Chunk, vector []float32) error {
	if len(vector) != s.dim {
		return fmt.Errorf("vector dimension = %d, want %d", len(vector), s.dim)
	}
	if chunk.ID == "" || chunk.DocumentID == "" || chunk.TenantID == "" {
		return errors.New("chunk id, document id, and tenant id are required")
	}
	serialized, err := sqlitevec.SerializeFloat32(vector)
	if err != nil {
		return fmt.Errorf("serialize vector: %w", err)
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM rag_vectors WHERE chunk_id = ?`, chunk.ID); err != nil {
		return fmt.Errorf("delete previous vector: %w", err)
	}
	if _, err := tx.Exec(`
		INSERT INTO rag_vectors(chunk_id, tenant_id, document_id, embedding, source_uri, chunk_text)
		VALUES (?, ?, ?, ?, ?, ?)
	`, chunk.ID, chunk.TenantID, chunk.DocumentID, serialized, chunk.URI, chunk.Text); err != nil {
		return fmt.Errorf("insert vector: %w", err)
	}
	return tx.Commit()
}

func (s *sqliteVecStore) Search(vector []float32, tenantID string, topK int) ([]RetrievedChunk, error) {
	if len(vector) != s.dim {
		return nil, fmt.Errorf("query vector dimension = %d, want %d", len(vector), s.dim)
	}
	if tenantID == "" || topK <= 0 {
		return nil, errors.New("tenant id and positive top-k are required")
	}
	serialized, err := sqlitevec.SerializeFloat32(vector)
	if err != nil {
		return nil, fmt.Errorf("serialize query: %w", err)
	}
	rows, err := s.db.Query(`
		SELECT chunk_id, document_id, source_uri, chunk_text, distance
		FROM rag_vectors
		WHERE embedding MATCH ? AND tenant_id = ? AND k = ?
		ORDER BY distance ASC
	`, serialized, tenantID, topK)
	if err != nil {
		return nil, fmt.Errorf("search vec0: %w", err)
	}
	defer rows.Close()
	var results []RetrievedChunk
	for rows.Next() {
		var result RetrievedChunk
		var distance float32
		if err := rows.Scan(&result.ID, &result.DocumentID, &result.URI, &result.Text, &distance); err != nil {
			return nil, fmt.Errorf("scan vector result: %w", err)
		}
		result.TenantID = tenantID
		// sqlite-vec 返回的是距离，不是可跨模型比较的相关性分数。
		// 此处仅为展示保留它；调用方不能拿不同模型的分数直接比较。
		result.Score = -distance
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *sqliteVecStore) DeleteByDocumentID(documentID string) error {
	if documentID == "" {
		return errors.New("document id is required")
	}
	_, err := s.db.Exec(`DELETE FROM rag_vectors WHERE document_id = ?`, documentID)
	return err
}
