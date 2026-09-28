package workflow

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agentgo/p_backend/internal/domain"
	"agentgo/p_backend/internal/llm"
	"agentgo/p_backend/internal/store"
)

type reportTestExplainer struct{}

func (reportTestExplainer) CheckAvailability(context.Context) error { return nil }
func (reportTestExplainer) Explain(context.Context, domain.Finding) (llm.Explanation, error) {
	return llm.Explanation{Content: "ok"}, nil
}
func (reportTestExplainer) AdviseIndex(context.Context, domain.Table, string) (llm.Explanation, error) {
	return llm.Explanation{Content: "人工复核"}, nil
}
func (reportTestExplainer) AdviseSchema(context.Context, domain.Table, string) (llm.Explanation, error) {
	return llm.Explanation{Content: "无需要人工复核的语义与治理风险"}, nil
}
func (reportTestExplainer) AdviseSchemaBatch(context.Context, []domain.Table, string) (llm.Explanation, error) {
	return llm.Explanation{Content: "无需要人工复核的语义与治理风险"}, nil
}
func (reportTestExplainer) AdviseDDL(context.Context, []domain.Finding, string) (llm.Explanation, error) {
	return llm.Explanation{Content: "无需要补充的 DDL 复核建议"}, nil
}
func (reportTestExplainer) SummarizeReport(context.Context, domain.Report, string) (llm.Explanation, error) {
	return llm.Explanation{Content: "# 修复后的报告\n\n**已恢复 Finding。**"}, nil
}

func TestReportEscapesFindingContentForHTMLPreview(t *testing.T) {
	report := buildReport("task-1", []domain.Finding{{ID: "finding-1", Severity: domain.SeverityUrgent, Title: "<script>alert(1)</script>", ObjectName: "orders", Evidence: "<img src=x>", Recommendation: "review", CreatedAt: time.Now().UTC()}}, nil)
	if strings.Contains(report.HTML, "<script>alert(1)</script>") || strings.Contains(report.HTML, "<img src=x>") {
		t.Fatalf("unsafe report preview: %s", report.HTML)
	}
	if !strings.Contains(report.HTML, "&lt;script&gt;") {
		t.Fatalf("expected escaped html: %s", report.HTML)
	}
}

func TestAppendDDLAppendixIncludesOnlyActionableDDL(t *testing.T) {
	report := appendDDLAppendix("# LLM 总结", []domain.Finding{
		{RuleID: "INDEX-003", Title: "缺少索引", DDLSuggestion: "CREATE INDEX idx_orders_user_id ON orders (user_id);", PrecheckSQL: "SHOW INDEX FROM orders;", RollbackPlan: "DROP INDEX idx_orders_user_id ON orders;"},
		{RuleID: "AI-INDEX-001", Title: "AI 评估", Recommendation: "人工复核"},
	})
	for _, expected := range []string{"## DDL 变更建议（需人工复核）", "INDEX-003 · 缺少索引", "```sql", "CREATE INDEX idx_orders_user_id ON orders (user_id);", "SHOW INDEX FROM orders;", "DROP INDEX idx_orders_user_id ON orders;"} {
		if !strings.Contains(report, expected) {
			t.Fatalf("report is missing %q: %s", expected, report)
		}
	}
	if strings.Contains(report, "AI 评估") {
		t.Fatalf("non-DDL finding must not be rendered in DDL appendix: %s", report)
	}
}

func TestRebuildReportRestoresFindingsFromPersistedSnapshot(t *testing.T) {
	ctx := context.Background()
	taskStore, err := store.Open(filepath.Join(t.TempDir(), "tasks.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer taskStore.Close()
	now := time.Now().UTC()
	task := domain.Task{ID: "task-rebuild", Status: domain.TaskCompleted, Stage: "completed", Target: domain.Target{Dialect: "mysql", Host: "db", Port: 3306, Database: "app"}, CreatedAt: now, UpdatedAt: now, ExpiresAt: now.Add(time.Hour), AccessTokenHash: "hash"}
	if err := taskStore.CreateTask(ctx, task, domain.DefaultWorkerProgress()); err != nil {
		t.Fatal(err)
	}
	snapshot := domain.SchemaSnapshot{Dialect: "mysql", Database: "app", CreatedAt: now, Tables: []domain.Table{{Name: "Orders-Log"}}}
	if err := taskStore.SaveSnapshot(ctx, task.ID, snapshot); err != nil {
		t.Fatal(err)
	}
	runner := New(taskStore, nil, WithReportSummarizer(reportTestExplainer{}), WithReportPrompt("生成 Markdown"))
	report, err := runner.RebuildReport(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Findings) != 3 || report.Summary[domain.SeverityUrgent] != 1 || report.Summary[domain.SeveritySuggestion] != 2 || !strings.Contains(report.Markdown, "修复后的报告") {
		t.Fatalf("unexpected rebuilt report: %#v", report)
	}
}

func TestSemanticBatchesRespectColumnBudget(t *testing.T) {
	runner := New(nil, nil, WithSemanticBatchMaxColumns(80))
	tables := []domain.Table{
		{Name: "accounts", Columns: semanticColumns(50, "user_id")},
		{Name: "orders", Columns: semanticColumns(40, "status")},
		{Name: "events", Columns: semanticColumns(20, "tenant_id")},
	}
	batches := runner.semanticBatches(tables)
	if len(batches) != 2 || batches[0].ColumnCount != 70 || batches[1].ColumnCount != 40 {
		t.Fatalf("unexpected semantic batches: %#v", batches)
	}
	for _, batch := range batches {
		if batch.ColumnCount > 80 {
			t.Fatalf("batch exceeds field budget: %#v", batch)
		}
	}
}

func semanticColumns(count int, firstName string) []domain.Column {
	columns := make([]domain.Column, count)
	for index := range columns {
		columns[index].Name = "field"
	}
	columns[0].Name = firstName
	return columns
}
