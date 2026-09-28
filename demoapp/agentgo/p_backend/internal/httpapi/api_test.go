package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"agentgo/p_backend/internal/domain"
	"agentgo/p_backend/internal/scanner"
	"agentgo/p_backend/internal/service"
	"agentgo/p_backend/internal/store"
	"agentgo/p_backend/internal/workflow"
)

func TestTaskEndpointsRequireTaskAccessToken(t *testing.T) {
	taskStore, err := store.Open(filepath.Join(t.TempDir(), "tasks.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer taskStore.Close()
	runner := workflow.New(taskStore, scanner.MultiCollector{MySQL: scanner.MySQLCollector{ConnectTimeout: 10 * time.Millisecond}})
	handler := New(taskStore, service.New(taskStore, runner, time.Hour))

	request := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(`{"connectionUrl":"mysql://readonly:secret@127.0.0.1:1/app"}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", response.Code, response.Body.String())
	}
	var created domain.TaskCreateResponse
	if err := json.NewDecoder(response.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.TaskID == "" || created.TaskAccessToken == "" {
		t.Fatalf("unexpected create response: %#v", created)
	}

	unauthorized := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+created.TaskID, nil)
	unauthorizedResponse := httptest.NewRecorder()
	handler.ServeHTTP(unauthorizedResponse, unauthorized)
	if unauthorizedResponse.Code != http.StatusNotFound {
		t.Fatalf("unauthorized status=%d", unauthorizedResponse.Code)
	}

	authorized := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+created.TaskID, nil)
	authorized.Header.Set("X-Task-Access-Token", created.TaskAccessToken)
	authorizedResponse := httptest.NewRecorder()
	handler.ServeHTTP(authorizedResponse, authorized)
	if authorizedResponse.Code != http.StatusOK {
		t.Fatalf("authorized status=%d body=%s", authorizedResponse.Code, authorizedResponse.Body.String())
	}
	var task domain.Task
	if err := json.NewDecoder(authorizedResponse.Body).Decode(&task); err != nil {
		t.Fatal(err)
	}
	if len(task.Progress.WorkerFlow) != 8 || task.Progress.WorkerFlow[1].ParallelGroup != "rule-analysis" || task.Progress.WorkerFlow[7].WorkerID != "report-generator" {
		t.Fatalf("unexpected worker flow: %#v", task.Progress.WorkerFlow)
	}
}

func TestReportEndpointPaginatesFindings(t *testing.T) {
	taskStore, err := store.Open(filepath.Join(t.TempDir(), "tasks.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer taskStore.Close()
	runner := workflow.New(taskStore, scanner.MultiCollector{MySQL: scanner.MySQLCollector{}})
	handler := New(taskStore, service.New(taskStore, runner, time.Hour))
	create := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(`{"connectionUrl":"mysql://readonly:secret@127.0.0.1:1/app"}`))
	createdResponse := httptest.NewRecorder()
	handler.ServeHTTP(createdResponse, create)
	var created domain.TaskCreateResponse
	if err := json.NewDecoder(createdResponse.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	report := domain.Report{TaskID: created.TaskID, GeneratedAt: time.Now().UTC(), RuleVersion: "test", Summary: map[domain.Severity]int{domain.SeverityUrgent: 3}, Findings: []domain.Finding{
		{ID: "one", Severity: domain.SeverityUrgent, Title: "one"}, {ID: "two", Severity: domain.SeverityUrgent, Title: "two"}, {ID: "three", Severity: domain.SeverityUrgent, Title: "three"},
	}, Markdown: "# report", HTML: "<h1>report</h1>"}
	if err := taskStore.SaveReport(context.Background(), report); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/reports/"+created.TaskID+"?page=2&page_size=1", nil)
	request.Header.Set("X-Task-Access-Token", created.TaskAccessToken)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var page domain.ReportPage
	if err := json.NewDecoder(response.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if page.Total != 3 || page.Page != 2 || len(page.Findings) != 1 || page.Findings[0].ID != "two" {
		t.Fatalf("unexpected page: %#v", page)
	}
}
