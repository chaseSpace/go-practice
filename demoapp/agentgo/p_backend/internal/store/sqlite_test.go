package store

import (
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"agentgo/p_backend/internal/domain"
)

func TestTaskProgressAndAccessTokenPersistence(t *testing.T) {
	ctx := context.Background()
	db, err := Open(filepath.Join(t.TempDir(), "tasks.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := time.Now().UTC()
	sum := sha256.Sum256([]byte("access-token"))
	task := domain.Task{ID: "task-1", Status: domain.TaskRunning, Stage: "analyzing", Target: domain.Target{Dialect: "mysql", Host: "db", Port: 3306, Database: "app"}, CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(time.Hour), AccessTokenHash: fmt.Sprintf("%x", sum)}
	worker := domain.WorkerProgress{WorkerID: "naming-convention", Name: "命名规范检查", Status: domain.WorkerPending}
	if err := db.CreateTask(ctx, task, []domain.WorkerProgress{worker}); err != nil {
		t.Fatal(err)
	}
	if err := db.Authorize(ctx, task.ID, "access-token"); err != nil {
		t.Fatalf("authorize task: %v", err)
	}
	if err := db.Authorize(ctx, task.ID, "wrong"); err != ErrUnauthorized {
		t.Fatalf("expected unauthorized, got %v", err)
	}
	if err := db.AddWorkUnits(ctx, task.ID, worker.WorkerID, []string{"orders"}); err != nil {
		t.Fatal(err)
	}
	if _, claimed, err := db.ClaimNextWorkUnit(ctx, task.ID, worker.WorkerID); err != nil || !claimed {
		t.Fatalf("claim work unit: claimed=%v err=%v", claimed, err)
	}
	if err := db.CompleteWorkUnit(ctx, task.ID, worker.WorkerID, "orders", nil); err != nil {
		t.Fatal(err)
	}
	if err := db.AddWorkerTokenUsage(ctx, task.ID, worker.WorkerID, domain.TokenUsage{InputTokens: 11, OutputTokens: 7, TotalTokens: 18}); err != nil {
		t.Fatal(err)
	}
	loaded, err := db.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Progress.TotalUnits != 1 || loaded.Progress.CompletedUnits != 1 || loaded.Progress.Percent != 100 {
		t.Fatalf("unexpected progress: %#v", loaded.Progress)
	}
	if loaded.Progress.TokenUsage.TotalTokens != 18 || loaded.Progress.Workers[0].TokenUsage.InputTokens != 11 {
		t.Fatalf("unexpected token usage: %#v", loaded.Progress)
	}
}

func TestRecoverInterruptedTaskPreservesWorkUnitForResume(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "restart.db")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	task := domain.Task{ID: "task-restart", Status: domain.TaskRunning, Stage: "analyzing", Target: domain.Target{Dialect: "mysql", Host: "db", Port: 3306, Database: "app"}, CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(time.Hour), AccessTokenHash: "hash"}
	worker := domain.WorkerProgress{WorkerID: "naming-convention", Name: "命名规范检查", Status: domain.WorkerRunning}
	if err := db.CreateTask(ctx, task, []domain.WorkerProgress{worker}); err != nil {
		t.Fatal(err)
	}
	if err := db.AddWorkUnits(ctx, task.ID, worker.WorkerID, []string{"orders"}); err != nil {
		t.Fatal(err)
	}
	if _, claimed, err := db.ClaimNextWorkUnit(ctx, task.ID, worker.WorkerID); err != nil || !claimed {
		t.Fatalf("claim unit: %v %v", claimed, err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	recovered, err := reopened.GetTask(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Status != domain.TaskPaused || !recovered.ConnectionRequired {
		t.Fatalf("unexpected recovered task: %#v", recovered)
	}
	if _, claimed, err := reopened.ClaimNextWorkUnit(ctx, task.ID, worker.WorkerID); err != nil || !claimed {
		t.Fatalf("expected recovered work unit: %v %v", claimed, err)
	}
}

func TestFindingsWithSameRuleIDAreStoredPerTask(t *testing.T) {
	ctx := context.Background()
	db, err := Open(filepath.Join(t.TempDir(), "tasks.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := time.Now().UTC()
	for _, taskID := range []string{"task-one", "task-two"} {
		task := domain.Task{ID: taskID, Status: domain.TaskRunning, Stage: "analyzing", Target: domain.Target{Dialect: "mysql", Host: "db", Port: 3306, Database: "app"}, CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(time.Hour), AccessTokenHash: "hash"}
		worker := domain.WorkerProgress{WorkerID: "naming-convention", Name: "命名规范检查", Status: domain.WorkerPending}
		if err := db.CreateTask(ctx, task, []domain.WorkerProgress{worker}); err != nil {
			t.Fatal(err)
		}
		if err := db.AddWorkUnits(ctx, taskID, worker.WorkerID, []string{"orders"}); err != nil {
			t.Fatal(err)
		}
		if _, claimed, err := db.ClaimNextWorkUnit(ctx, taskID, worker.WorkerID); err != nil || !claimed {
			t.Fatalf("claim work unit: claimed=%v err=%v", claimed, err)
		}
		finding := domain.Finding{ID: "name-001-same", RuleID: "NAME-001", Severity: domain.SeveritySuggestion, ObjectName: "orders", CreatedAt: now}
		if err := db.CompleteWorkUnit(ctx, taskID, worker.WorkerID, "orders", []domain.Finding{finding}); err != nil {
			t.Fatal(err)
		}
		findings, err := db.ListFindings(ctx, taskID)
		if err != nil || len(findings) != 1 {
			t.Fatalf("task %s findings=%#v err=%v", taskID, findings, err)
		}
	}
}
