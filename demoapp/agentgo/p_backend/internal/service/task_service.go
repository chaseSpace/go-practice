package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"agentgo/p_backend/internal/domain"
	"agentgo/p_backend/internal/workflow"
)

var (
	ErrInvalidConnection = errors.New("请填写有效的 MySQL 连接信息")
	ErrInvalidState      = errors.New("任务当前状态不支持该操作")
)

type TaskRepository interface {
	CreateTask(context.Context, domain.Task, []domain.WorkerProgress) error
	GetTask(context.Context, string) (domain.Task, error)
	SaveTask(context.Context, domain.Task) error
}

type Service struct {
	store  TaskRepository
	runner *workflow.Runner
	ttl    time.Duration
	mu     sync.Mutex
}

func New(store TaskRepository, runner *workflow.Runner, ttl time.Duration) *Service {
	return &Service{store: store, runner: runner, ttl: ttl}
}

func (s *Service) Create(ctx context.Context, request domain.CreateTaskRequest) (domain.TaskCreateResponse, error) {
	connection, err := ParseConnection(request.ConnectionURL, request.Connection)
	if err != nil {
		return domain.TaskCreateResponse{}, err
	}
	now := time.Now().UTC()
	token, err := randomSecret(32)
	if err != nil {
		return domain.TaskCreateResponse{}, err
	}
	taskID, err := randomSecret(16)
	if err != nil {
		return domain.TaskCreateResponse{}, err
	}
	task := domain.Task{
		ID: taskID, Status: domain.TaskPending, Stage: "pending", CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(s.ttl),
		Target:          domain.Target{Dialect: connection.Dialect, Host: connection.Host, Port: connection.Port, Database: connection.Database},
		AccessTokenHash: hashToken(token), RedactObjectNames: request.RedactObjectNames,
	}
	workers := domain.DefaultWorkerProgress()
	if err := s.store.CreateTask(ctx, task, workers); err != nil {
		return domain.TaskCreateResponse{}, err
	}
	s.runner.SetRuntime(task.ID, workflow.Runtime{Connection: connection, ExcludedTables: request.ExcludeTables})
	s.runner.Start(task.ID)
	return domain.TaskCreateResponse{TaskID: task.ID, TaskAccessToken: token, ExpiresAt: task.ExpiresAt}, nil
}

func (s *Service) Pause(ctx context.Context, taskID string) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, err := s.store.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, err
	}
	if task.Status == domain.TaskPausing || task.Status == domain.TaskPaused {
		return task, nil
	}
	if task.Status != domain.TaskRunning {
		return domain.Task{}, ErrInvalidState
	}
	now := time.Now().UTC()
	task.Status, task.PauseRequestedAt, task.UpdatedAt = domain.TaskPausing, &now, now
	if err := s.store.SaveTask(ctx, task); err != nil {
		return domain.Task{}, err
	}
	return s.store.GetTask(ctx, taskID)
}

func (s *Service) Resume(ctx context.Context, taskID string, request domain.ResumeTaskRequest) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, err := s.store.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, err
	}
	if task.Status != domain.TaskPaused && task.Status != domain.TaskFailed {
		return domain.Task{}, ErrInvalidState
	}
	if !s.runner.HasRuntime(taskID) {
		connection, err := ParseConnection(request.ConnectionURL, request.Connection)
		if err != nil {
			return task, fmt.Errorf("connection required: %w", err)
		}
		if connection.Dialect != task.Target.Dialect || connection.Host != task.Target.Host || connection.Port != task.Target.Port || connection.Database != task.Target.Database {
			return task, errors.New("恢复任务时连接目标必须与原任务一致")
		}
		s.runner.SetRuntime(taskID, workflow.Runtime{Connection: connection})
	}
	if task.Status == domain.TaskFailed {
		if err := s.runner.RecoverFailed(ctx, taskID); err != nil {
			return domain.Task{}, err
		}
	}
	now := time.Now().UTC()
	task.Status, task.Stage, task.PauseRequestedAt, task.PausedAt, task.ConnectionRequired, task.FailureSummary, task.UpdatedAt = domain.TaskRunning, "resuming", nil, nil, false, "", now
	if err := s.store.SaveTask(ctx, task); err != nil {
		return domain.Task{}, err
	}
	s.runner.Start(taskID)
	return s.store.GetTask(ctx, taskID)
}

func (s *Service) Cancel(ctx context.Context, taskID string) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, err := s.store.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, err
	}
	if task.Status.Terminal() {
		return task, nil
	}
	now := time.Now().UTC()
	task.Status, task.Stage, task.UpdatedAt = domain.TaskCancelled, "cancelled", now
	if err := s.store.SaveTask(ctx, task); err != nil {
		return domain.Task{}, err
	}
	s.runner.ClearRuntime(taskID)
	return s.store.GetTask(ctx, taskID)
}

// RegenerateReport only renders a new report from the persisted Finding set.
func (s *Service) RegenerateReport(ctx context.Context, taskID string) (domain.Report, error) {
	return s.runner.RegenerateReport(ctx, taskID)
}

// ReanalyzeReport recomputes workers from the persisted Schema snapshot before
// rendering a report. It never reconnects to the target database.
func (s *Service) ReanalyzeReport(ctx context.Context, taskID string) (domain.Report, error) {
	return s.runner.ReanalyzeReport(ctx, taskID)
}

func ParseConnection(rawURL string, structured *domain.ConnectionConfig) (domain.ConnectionConfig, error) {
	if strings.TrimSpace(rawURL) != "" && structured != nil {
		return domain.ConnectionConfig{}, ErrInvalidConnection
	}
	if strings.TrimSpace(rawURL) == "" && structured == nil {
		return domain.ConnectionConfig{}, ErrInvalidConnection
	}
	if strings.TrimSpace(rawURL) != "" {
		parsed, err := url.Parse(rawURL)
		dialect := strings.ToLower(parsed.Scheme)
		if err != nil || (dialect != "mysql" && dialect != "postgres" && dialect != "postgresql") || parsed.User == nil {
			return domain.ConnectionConfig{}, ErrInvalidConnection
		}
		password, _ := parsed.User.Password()
		port := defaultPort(dialect)
		if rawPort := parsed.Port(); rawPort != "" {
			port, err = strconv.Atoi(rawPort)
			if err != nil || port < 1 || port > 65535 {
				return domain.ConnectionConfig{}, ErrInvalidConnection
			}
		}
		config := domain.ConnectionConfig{Dialect: dialect, Host: parsed.Hostname(), Port: port, Database: strings.TrimPrefix(parsed.EscapedPath(), "/"), Username: parsed.User.Username(), Password: password}
		if config.Database, err = url.PathUnescape(config.Database); err != nil {
			return domain.ConnectionConfig{}, ErrInvalidConnection
		}
		return validateConnection(config)
	}
	return validateConnection(*structured)
}

func validateConnection(config domain.ConnectionConfig) (domain.ConnectionConfig, error) {
	config.Dialect = strings.ToLower(strings.TrimSpace(config.Dialect))
	config.Host = strings.TrimSpace(config.Host)
	config.Database = strings.TrimSpace(config.Database)
	config.Username = strings.TrimSpace(config.Username)
	if (config.Dialect != "mysql" && config.Dialect != "postgres" && config.Dialect != "postgresql") || config.Host == "" || config.Database == "" || config.Username == "" {
		return domain.ConnectionConfig{}, ErrInvalidConnection
	}
	if config.Port == 0 {
		config.Port = defaultPort(config.Dialect)
	}
	if config.Port < 1 || config.Port > 65535 {
		return domain.ConnectionConfig{}, ErrInvalidConnection
	}
	if forbiddenHost(config.Host) {
		return domain.ConnectionConfig{}, ErrInvalidConnection
	}
	return config, nil
}

func defaultPort(dialect string) int {
	if dialect == "postgres" || dialect == "postgresql" {
		return 5432
	}
	return 3306
}

func forbiddenHost(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	if ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	return ip.String() == "169.254.169.254"
}

func randomSecret(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
