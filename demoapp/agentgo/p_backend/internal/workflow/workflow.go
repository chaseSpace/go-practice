package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"

	"agentgo/p_backend/internal/domain"
	"agentgo/p_backend/internal/llm"
	"agentgo/p_backend/internal/rules"
	"agentgo/p_backend/internal/scanner"
	"agentgo/p_backend/internal/store"

	"github.com/cloudwego/eino/compose"
)

var (
	errPaused    = errors.New("task paused")
	errCancelled = errors.New("task cancelled")
)

type Runtime struct {
	Connection     domain.ConnectionConfig
	ExcludedTables []string
}

type SemanticRule struct {
	ID     string
	Title  string
	Prompt string
}

type Runner struct {
	store                      *store.Store
	collector                  scanner.Collector
	schemaTool                 scanner.SchemaCollectorTool
	rules                      *rules.Engine
	indexAdvisor               llm.Explainer
	indexAdviceMinRows         int64
	indexAdvicePrompt          string
	semanticAdvisor            llm.Explainer
	semanticRules              []SemanticRule
	semanticBatchMaxColumns    int
	semanticMaxParallelBatches int
	ddlAdvisor                 llm.Explainer
	ddlPrompt                  string
	ddlMaxInputFindings        int
	reportSummarizer           llm.Explainer
	reportPrompt               string
	ruleVersion                string
	mu                         sync.RWMutex
	runtime                    map[string]Runtime
	active                     map[string]bool
}

type Option func(*Runner)

func WithRuleEngine(engine *rules.Engine) Option {
	return func(runner *Runner) {
		if engine != nil {
			runner.rules = engine
		}
	}
}

func WithIndexAdvisor(advisor llm.Explainer) Option {
	return func(runner *Runner) {
		if advisor != nil {
			runner.indexAdvisor = advisor
		}
	}
}

func WithIndexAdviceMinRows(rows int64) Option {
	return func(runner *Runner) {
		if rows > 0 {
			runner.indexAdviceMinRows = rows
		}
	}
}

func WithIndexAdvicePrompt(prompt string) Option {
	return func(runner *Runner) { runner.indexAdvicePrompt = prompt }
}

func WithSemanticAdvisor(advisor llm.Explainer) Option {
	return func(runner *Runner) {
		if advisor != nil {
			runner.semanticAdvisor = advisor
		}
	}
}

func WithSemanticRules(rules []SemanticRule) Option {
	return func(runner *Runner) {
		runner.semanticRules = append([]SemanticRule(nil), rules...)
	}
}

func WithSemanticBatchMaxColumns(columns int) Option {
	return func(runner *Runner) {
		if columns > 0 {
			runner.semanticBatchMaxColumns = columns
		}
	}
}

func WithSemanticMaxParallelBatches(parallel int) Option {
	return func(runner *Runner) {
		if parallel > 0 {
			runner.semanticMaxParallelBatches = parallel
		}
	}
}

func WithDDLAdvisor(advisor llm.Explainer) Option {
	return func(runner *Runner) {
		if advisor != nil {
			runner.ddlAdvisor = advisor
		}
	}
}

func WithDDLPrompt(prompt string) Option {
	return func(runner *Runner) { runner.ddlPrompt = prompt }
}

func WithDDLMaxInputFindings(max int) Option {
	return func(runner *Runner) {
		if max > 0 {
			runner.ddlMaxInputFindings = max
		}
	}
}

func WithReportSummarizer(summarizer llm.Explainer) Option {
	return func(runner *Runner) {
		if summarizer != nil {
			runner.reportSummarizer = summarizer
		}
	}
}

func WithReportPrompt(prompt string) Option {
	return func(runner *Runner) { runner.reportPrompt = prompt }
}

func WithRuleVersion(version string) Option {
	return func(runner *Runner) {
		if strings.TrimSpace(version) != "" {
			runner.ruleVersion = version
		}
	}
}

func New(store *store.Store, collector scanner.Collector, options ...Option) *Runner {
	runner := &Runner{store: store, collector: collector, schemaTool: scanner.SchemaCollectorTool{Collector: collector}, rules: rules.NewDefault(), indexAdvisor: llm.DisabledExplainer{}, indexAdviceMinRows: 100_000, semanticAdvisor: llm.DisabledExplainer{}, semanticBatchMaxColumns: 80, semanticMaxParallelBatches: 3, ddlAdvisor: llm.DisabledExplainer{}, ddlMaxInputFindings: 40, reportSummarizer: llm.DisabledExplainer{}, ruleVersion: "builtin-0.1", runtime: make(map[string]Runtime), active: make(map[string]bool)}
	for _, option := range options {
		option(runner)
	}
	return runner
}

func (r *Runner) SetRuntime(taskID string, runtime Runtime) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.runtime[taskID] = runtime
}

func (r *Runner) HasRuntime(taskID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.runtime[taskID]
	return ok
}

func (r *Runner) ClearRuntime(taskID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.runtime, taskID)
}

func (r *Runner) runtimeFor(taskID string) (Runtime, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	runtime, ok := r.runtime[taskID]
	return runtime, ok
}

func (r *Runner) Start(taskID string) {
	r.mu.Lock()
	if r.active[taskID] {
		r.mu.Unlock()
		return
	}
	r.active[taskID] = true
	r.mu.Unlock()
	slog.Info("task execution started", "task_id", taskID)
	go func() {
		defer func() {
			r.mu.Lock()
			delete(r.active, taskID)
			r.mu.Unlock()
		}()
		r.run(context.Background(), taskID)
	}()
}

func (r *Runner) run(ctx context.Context, taskID string) {
	chain := compose.NewChain[string, string]()
	chain.AppendLambda(compose.InvokableLambda(r.collectStage), compose.WithNodeName("SchemaCollector"))
	chain.AppendLambda(compose.InvokableLambda(r.analyzeStage), compose.WithNodeName("RuleWorkers"))
	chain.AppendLambda(compose.InvokableLambda(r.reportStage), compose.WithNodeName("ReportGenerator"))
	runnable, err := chain.Compile(ctx)
	if err == nil {
		_, err = runnable.Invoke(ctx, taskID)
	}
	r.finish(ctx, taskID, err)
}

func (r *Runner) finish(ctx context.Context, taskID string, flowErr error) {
	task, err := r.store.GetTask(ctx, taskID)
	if err != nil || task.Status == domain.TaskCancelled || task.Status.Terminal() {
		if task.Status.Terminal() {
			r.ClearRuntime(taskID)
		}
		return
	}
	now := time.Now().UTC()
	switch {
	case errors.Is(flowErr, errPaused):
		task.Status, task.UpdatedAt, task.PausedAt = domain.TaskPaused, now, &now
		task.PauseRequestedAt = &now
		_ = r.store.SaveTask(ctx, task)
		_ = r.store.SetAllRunningWorkers(ctx, taskID, domain.WorkerPaused)
		slog.Info("task execution paused", "task_id", taskID, "stage", task.Stage)
	case errors.Is(flowErr, errCancelled):
		task.Status, task.Stage, task.UpdatedAt = domain.TaskCancelled, "cancelled", now
		_ = r.store.SaveTask(ctx, task)
		r.ClearRuntime(taskID)
		slog.Info("task execution cancelled", "task_id", taskID)
	case flowErr != nil:
		task.Status, task.Stage, task.UpdatedAt = domain.TaskFailed, "failed", now
		task.FailureSummary = safeError(flowErr)
		task.ConnectionRequired = true
		_ = r.store.SaveTask(ctx, task)
		_ = r.store.SetAllRunningWorkers(ctx, taskID, domain.WorkerFailed)
		r.ClearRuntime(taskID)
		slog.Warn("task execution failed", "task_id", taskID, "stage", task.Stage)
	default:
		task.Status, task.Stage, task.UpdatedAt = domain.TaskCompleted, "completed", now
		_ = r.store.SaveTask(ctx, task)
		r.ClearRuntime(taskID)
		slog.Info("task execution completed", "task_id", taskID)
	}
}

func (r *Runner) RecoverFailed(ctx context.Context, taskID string) error {
	if err := r.store.RecoverFailedWork(ctx, taskID); err != nil {
		return err
	}
	if err := r.rebatchLegacySemanticWork(ctx, taskID); err != nil {
		return err
	}
	slog.Info("failed task work recovered for resume", "task_id", taskID)
	return nil
}

func (r *Runner) rebatchLegacySemanticWork(ctx context.Context, taskID string) error {
	const workerID = "semantic-governance"
	pending, err := r.store.ListPendingWorkUnitKeys(ctx, taskID, workerID)
	if err != nil || len(pending) == 0 {
		return err
	}
	for _, key := range pending {
		if strings.HasPrefix(strings.TrimSpace(key), "{") {
			return nil
		}
	}
	snapshot, err := r.store.LoadSnapshot(ctx, taskID)
	if err != nil {
		return err
	}
	byName := make(map[string]domain.Table, len(snapshot.Tables))
	for _, table := range snapshot.Tables {
		byName[table.Name] = table
	}
	tables := make([]domain.Table, 0, len(pending))
	for _, key := range pending {
		if table, ok := byName[key]; ok && shouldRunSemanticReview(table) {
			tables = append(tables, table)
		}
	}
	batches := r.semanticBatches(tables)
	keys := make([]string, 0, len(batches))
	for _, batch := range batches {
		key, err := batch.key()
		if err != nil {
			return err
		}
		keys = append(keys, key)
	}
	if err := r.store.ReplacePendingWorkUnits(ctx, taskID, workerID, keys); err != nil {
		return err
	}
	slog.Info("legacy semantic work rebatch completed", "task_id", taskID, "old_pending_units", len(pending), "new_batches", len(keys), "batch_max_columns", r.semanticBatchMaxColumns)
	return nil
}

func (r *Runner) collectStage(ctx context.Context, taskID string) (string, error) {
	if err := r.checkControl(ctx, taskID); err != nil {
		return taskID, err
	}
	worker, err := r.store.GetWorker(ctx, taskID, "schema-collector")
	if err != nil {
		return taskID, err
	}
	if worker.Status == domain.WorkerCompleted {
		return taskID, nil
	}
	if err := r.setStage(ctx, taskID, "collecting_schema"); err != nil {
		return taskID, err
	}
	worker.Status, worker.TotalUnits, worker.CurrentObject = domain.WorkerRunning, 1, "information_schema"
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return taskID, err
	}
	slog.Info("worker execution started", "task_id", taskID, "worker_id", worker.WorkerID, "worker_name", worker.Name, "total_units", worker.TotalUnits)
	runtime, ok := r.runtimeFor(taskID)
	if !ok {
		return taskID, r.pauseForConnection(ctx, taskID)
	}
	snapshot, err := r.schemaTool.Collect(ctx, runtime.Connection, runtime.ExcludedTables)
	if err != nil {
		return taskID, err
	}
	if err := r.checkControl(ctx, taskID); err != nil {
		return taskID, err
	}
	if err := r.store.SaveSnapshot(ctx, taskID, snapshot); err != nil {
		return taskID, err
	}
	slog.Info("schema collection completed", "task_id", taskID, "table_count", len(snapshot.Tables), "unchecked_count", len(snapshot.UncheckedObjects))
	keys := make([]string, 0, len(snapshot.Tables))
	for _, table := range snapshot.Tables {
		keys = append(keys, table.Name)
	}
	for _, workerID := range []string{"naming-convention", "index-performance", "type-constraint"} {
		if err := r.store.AddWorkUnits(ctx, taskID, workerID, keys); err != nil {
			return taskID, err
		}
	}
	semanticBatches := r.semanticBatches(snapshot.Tables)
	semanticKeys := make([]string, 0, len(semanticBatches))
	for _, batch := range semanticBatches {
		key, err := batch.key()
		if err != nil {
			return taskID, err
		}
		semanticKeys = append(semanticKeys, key)
	}
	if err := r.store.AddWorkUnits(ctx, taskID, "semantic-governance", semanticKeys); err != nil {
		return taskID, err
	}
	if err := r.store.AddWorkUnits(ctx, taskID, "ddl-generator", []string{"ddl-review"}); err != nil {
		return taskID, err
	}
	slog.Info("semantic governance batches created", "task_id", taskID, "batch_count", len(semanticKeys), "batch_max_columns", r.semanticBatchMaxColumns, "max_parallel_batches", r.semanticMaxParallelBatches)
	worker, err = r.store.GetWorker(ctx, taskID, "schema-collector")
	if err != nil {
		return taskID, err
	}
	worker.Status, worker.CompletedUnits, worker.CurrentObject, worker.Checkpoint = domain.WorkerCompleted, 1, "", "schema_collected"
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return taskID, err
	}
	slog.Info("worker completed", "task_id", taskID, "worker_id", worker.WorkerID, "completed_units", worker.CompletedUnits, "total_units", worker.TotalUnits)
	return taskID, nil
}

func (r *Runner) analyzeStage(ctx context.Context, taskID string) (string, error) {
	if err := r.checkControl(ctx, taskID); err != nil {
		return taskID, err
	}
	if err := r.setStage(ctx, taskID, "analyzing"); err != nil {
		return taskID, err
	}
	snapshot, err := r.store.LoadSnapshot(ctx, taskID)
	if err != nil {
		return taskID, err
	}
	tables := make(map[string]domain.Table, len(snapshot.Tables))
	for _, table := range snapshot.Tables {
		tables[table.Name] = table
	}

	workerIDs := []string{"naming-convention", "index-performance", "type-constraint", "semantic-governance"}
	errCh := make(chan error, len(workerIDs))
	var group sync.WaitGroup
	for _, workerID := range workerIDs {
		group.Add(1)
		go func(workerID string) {
			defer group.Done()
			errCh <- r.runEngineWorker(ctx, taskID, workerID, snapshot, tables)
		}(workerID)
	}
	group.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return taskID, err
		}
	}
	if err := r.completeSimpleWorker(ctx, taskID, "finding-normalizer", "findings_normalized"); err != nil {
		return taskID, err
	}
	if err := r.runDDLWorker(ctx, taskID); err != nil {
		return taskID, err
	}
	return taskID, nil
}

func (r *Runner) runDDLWorker(ctx context.Context, taskID string) error {
	const workerID = "ddl-generator"
	worker, err := r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	if worker.Status == domain.WorkerCompleted {
		return nil
	}
	worker.Status, worker.CurrentObject, worker.Checkpoint = domain.WorkerRunning, "ddl-review", "reviewing_ddl"
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	key, found, err := r.store.ClaimNextWorkUnit(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	if !found {
		worker, err = r.store.GetWorker(ctx, taskID, workerID)
		if err != nil {
			return err
		}
		worker.Status, worker.CurrentObject, worker.Checkpoint = domain.WorkerCompleted, "", "ddl_review_completed"
		return r.store.SetWorker(ctx, taskID, worker)
	}
	allFindings, err := r.store.ListFindings(ctx, taskID)
	if err != nil {
		return err
	}
	inputs := make([]domain.Finding, 0)
	for _, finding := range allFindings {
		if strings.HasPrefix(finding.RuleID, "INDEX-") || finding.RuleID == "AI-INDEX-001" {
			inputs = append(inputs, finding)
		}
	}
	if len(inputs) > r.ddlMaxInputFindings {
		inputs = inputs[:r.ddlMaxInputFindings]
	}
	generated := make([]domain.Finding, 0, 1)
	if len(inputs) > 0 {
		explanation, err := r.ddlAdvisor.AdviseDDL(ctx, inputs, r.ddlPrompt)
		if err != nil || strings.TrimSpace(explanation.Content) == "" {
			_ = r.failWorker(ctx, taskID, workerID, "LLM DDL 建议复核失败，请检查模型服务配置与可用性")
			slog.Warn("DDL AI worker failed", "task_id", taskID, "worker_id", workerID, "failure_code", "ddl_llm_call_failed")
			return errors.New("LLM DDL review failed")
		}
		if err := r.store.AddWorkerTokenUsage(ctx, taskID, workerID, explanation.TokenUsage); err != nil {
			return err
		}
		if !strings.Contains(explanation.Content, "无需要补充的 DDL 复核建议") {
			generated = append(generated, domain.Finding{ID: "ai-ddl-review", RuleID: "AI-DDL-001", Severity: domain.SeveritySuggestion, SeverityReason: "AI DDL 建议人工复核", ObjectType: "task", ObjectName: "DDL 建议", Title: "AI DDL 建议人工复核", Evidence: fmt.Sprintf("基于 %d 条确定性 DDL 变更候选复核，不包含业务数据", len(inputs)), Recommendation: explanation.Content, RequiresManualReview: true, SQLRisk: "medium", CreatedAt: time.Now().UTC()})
		}
	}
	if err := r.store.CompleteWorkUnitWithProgress(ctx, taskID, workerID, key, "ddl-review", generated); err != nil {
		return err
	}
	worker, err = r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	worker.Status, worker.CurrentObject, worker.Checkpoint = domain.WorkerCompleted, "", "ddl_review_completed"
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("DDL AI worker completed", "task_id", taskID, "worker_id", workerID, "reviewed_findings", len(inputs), "finding_count", worker.FindingCount, "token_total", worker.TokenUsage.TotalTokens)
	return nil
}

func (r *Runner) reportStage(ctx context.Context, taskID string) (string, error) {
	if err := r.checkControl(ctx, taskID); err != nil {
		return taskID, err
	}
	if err := r.setStage(ctx, taskID, "generating_report"); err != nil {
		return taskID, err
	}
	if err := r.startSimpleWorker(ctx, taskID, "report-generator", "summarizing_worker_results"); err != nil {
		return taskID, err
	}
	slog.Info("LLM report generation started", "task_id", taskID)
	findings, err := r.store.ListFindings(ctx, taskID)
	if err != nil {
		return taskID, err
	}
	task, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return taskID, err
	}
	if task.RedactObjectNames {
		for index := range findings {
			findings[index].ObjectName = redact(findings[index].ObjectName)
		}
	}
	snapshot, err := r.store.LoadSnapshot(ctx, taskID)
	if err != nil {
		return taskID, err
	}
	report := buildReport(taskID, findings, snapshot.UncheckedObjects)
	report.RuleVersion = r.ruleVersion
	report.Workers = projectedCompletedWorker(task.Progress.Workers, "report-generator", "report_saved")
	report, usage, err := r.summarizeReport(ctx, report)
	if err != nil {
		_ = r.failWorker(ctx, taskID, "report-generator", "LLM 报告生成失败，请检查模型服务配置与可用性")
		return taskID, errors.New("LLM report generation failed")
	}
	if err := r.store.AddWorkerTokenUsage(ctx, taskID, "report-generator", usage); err != nil {
		return taskID, err
	}
	updatedTask, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return taskID, err
	}
	report.Workers = projectedCompletedWorker(updatedTask.Progress.Workers, "report-generator", "report_saved")
	if err := r.store.SaveReport(ctx, report); err != nil {
		return taskID, err
	}
	slog.Info("LLM Markdown report saved", "task_id", taskID, "finding_count", len(report.Findings), "worker_count", len(report.Workers), "token_total", usage.TotalTokens)
	if err := r.completeSimpleWorker(ctx, taskID, "report-generator", "report_saved"); err != nil {
		return taskID, err
	}
	if len(snapshot.UncheckedObjects) > 0 {
		task.Status = domain.TaskPartialSuccess
		task.FailureSummary = "部分对象因元数据采集失败未检查"
		task.UpdatedAt = time.Now().UTC()
		if err := r.store.SaveTask(ctx, task); err != nil {
			return taskID, err
		}
		slog.Warn("task completed with unchecked objects", "task_id", taskID, "unchecked_count", len(snapshot.UncheckedObjects))
	}
	return taskID, nil
}

// ReanalyzeReport re-runs deterministic and AI analysis from the saved Schema
// snapshot. It is reserved for repairing legacy/missing findings and is not
// the normal report regeneration path.
func (r *Runner) ReanalyzeReport(ctx context.Context, taskID string) (domain.Report, error) {
	ctx = llm.WithInputLogging(ctx)
	snapshot, err := r.store.LoadSnapshot(ctx, taskID)
	if err != nil {
		return domain.Report{}, err
	}
	recovered := make([]domain.Finding, 0)
	for _, workerID := range []string{"naming-convention", "index-performance", "type-constraint"} {
		for _, table := range snapshot.Tables {
			findings, usage, err := r.findingsForWorker(ctx, workerID, snapshot, table)
			if err != nil {
				return domain.Report{}, err
			}
			if err := r.store.AddWorkerTokenUsage(ctx, taskID, workerID, usage); err != nil {
				return domain.Report{}, err
			}
			recovered = append(recovered, findings...)
		}
	}
	byName := make(map[string]domain.Table, len(snapshot.Tables))
	for _, table := range snapshot.Tables {
		byName[table.Name] = table
	}
	for _, batch := range r.semanticBatches(snapshot.Tables) {
		batchTables := make([]domain.Table, 0, len(batch.Tables))
		for _, name := range batch.Tables {
			batchTables = append(batchTables, byName[name])
		}
		findings, usage, err := r.findingsForSemanticBatch(ctx, batch, batchTables)
		if err != nil {
			return domain.Report{}, err
		}
		if err := r.store.AddWorkerTokenUsage(ctx, taskID, "semantic-governance", usage); err != nil {
			return domain.Report{}, err
		}
		recovered = append(recovered, findings...)
	}
	if err := r.store.SaveFindings(ctx, taskID, recovered); err != nil {
		return domain.Report{}, err
	}
	findings, err := r.store.ListFindings(ctx, taskID)
	if err != nil {
		return domain.Report{}, err
	}
	return r.generateReportFromFindings(ctx, taskID, snapshot, findings, "report reanalyzed from persisted schema snapshot")
}

// RegenerateReport only renders a new reader-friendly report from saved
// findings. It deliberately avoids re-running any worker or Schema analysis.
func (r *Runner) RegenerateReport(ctx context.Context, taskID string) (domain.Report, error) {
	ctx = llm.WithInputLogging(ctx)
	snapshot, err := r.store.LoadSnapshot(ctx, taskID)
	if err != nil {
		return domain.Report{}, err
	}
	findings, err := r.store.ListFindings(ctx, taskID)
	if err != nil {
		return domain.Report{}, err
	}
	return r.generateReportFromFindings(ctx, taskID, snapshot, findings, "report regenerated from saved findings")
}

// RebuildReport is retained for callers compiled against the old workflow API.
// New callers should choose RegenerateReport or ReanalyzeReport explicitly.
func (r *Runner) RebuildReport(ctx context.Context, taskID string) (domain.Report, error) {
	return r.ReanalyzeReport(ctx, taskID)
}

func (r *Runner) generateReportFromFindings(ctx context.Context, taskID string, snapshot domain.SchemaSnapshot, findings []domain.Finding, completionLog string) (domain.Report, error) {
	task, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return domain.Report{}, err
	}
	if task.RedactObjectNames {
		for index := range findings {
			findings[index].ObjectName = redact(findings[index].ObjectName)
		}
	}
	report := buildReport(taskID, findings, snapshot.UncheckedObjects)
	report.RuleVersion = r.ruleVersion
	report.Workers = task.Progress.Workers
	report, usage, err := r.summarizeReport(ctx, report)
	if err != nil {
		return domain.Report{}, err
	}
	if err := r.store.AddWorkerTokenUsage(ctx, taskID, "report-generator", usage); err != nil {
		return domain.Report{}, err
	}
	updatedTask, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return domain.Report{}, err
	}
	report.Workers = updatedTask.Progress.Workers
	if err := r.store.SaveReport(ctx, report); err != nil {
		return domain.Report{}, err
	}
	slog.Info(completionLog, "task_id", taskID, "finding_count", len(report.Findings))
	return report, nil
}

func (r *Runner) summarizeReport(ctx context.Context, report domain.Report) (domain.Report, domain.TokenUsage, error) {
	explanation, err := r.reportSummarizer.SummarizeReport(ctx, report, r.reportPrompt)
	if err != nil || strings.TrimSpace(explanation.Content) == "" {
		return domain.Report{}, domain.TokenUsage{}, errors.New("LLM report generation failed")
	}
	// The LLM writes the reader-friendly summary, but executable DDL is kept
	// deterministic and is appended verbatim. This prevents the summary prompt
	// (which deliberately forbids model-generated SQL) from hiding a rule's
	// already-generated DDL suggestion.
	report.Markdown = appendDDLAppendix(explanation.Content, report.Findings)
	return report, explanation.TokenUsage, nil
}

func appendDDLAppendix(markdown string, findings []domain.Finding) string {
	ddlFindings := make([]domain.Finding, 0)
	for _, finding := range findings {
		if strings.TrimSpace(finding.DDLSuggestion) != "" {
			ddlFindings = append(ddlFindings, finding)
		}
	}
	if len(ddlFindings) == 0 {
		return strings.TrimSpace(markdown)
	}

	var result strings.Builder
	result.WriteString(strings.TrimSpace(markdown))
	result.WriteString("\n\n## DDL 变更建议（需人工复核）\n\n")
	result.WriteString("以下语句由确定性规则生成。执行前请在目标环境核验，并按回滚方案操作。\n")
	for _, finding := range ddlFindings {
		fmt.Fprintf(&result, "\n### %s · %s\n", finding.RuleID, finding.Title)
		if precheck := strings.TrimSpace(finding.PrecheckSQL); precheck != "" {
			fmt.Fprintf(&result, "\n执行前核验：\n\n```sql\n%s\n```\n", precheck)
		}
		fmt.Fprintf(&result, "\n```sql\n%s\n```\n", strings.TrimSpace(finding.DDLSuggestion))
		if rollback := strings.TrimSpace(finding.RollbackPlan); rollback != "" {
			fmt.Fprintf(&result, "\n回滚说明：%s\n", rollback)
		}
	}
	return result.String()
}

type semanticBatch struct {
	Index       int      `json:"index"`
	Tables      []string `json:"tables"`
	ColumnCount int      `json:"columnCount"`
}

func (b semanticBatch) key() (string, error) {
	payload, err := json.Marshal(b)
	return string(payload), err
}

func (r *Runner) semanticBatches(tables []domain.Table) []semanticBatch {
	candidates := make([]domain.Table, 0)
	for _, table := range tables {
		if shouldRunSemanticReview(table) {
			candidates = append(candidates, table)
		}
	}
	sort.Slice(candidates, func(left, right int) bool { return candidates[left].Name < candidates[right].Name })
	batches := make([]semanticBatch, 0)
	current := semanticBatch{Index: 1}
	for _, table := range candidates {
		columns := len(table.Columns)
		if len(current.Tables) > 0 && current.ColumnCount+columns > r.semanticBatchMaxColumns {
			batches = append(batches, current)
			current = semanticBatch{Index: len(batches) + 1}
		}
		current.Tables = append(current.Tables, table.Name)
		current.ColumnCount += columns
	}
	if len(current.Tables) > 0 {
		batches = append(batches, current)
	}
	return batches
}

func (r *Runner) runEngineWorker(ctx context.Context, taskID, workerID string, snapshot domain.SchemaSnapshot, tables map[string]domain.Table) error {
	if workerID == "semantic-governance" {
		return r.runSemanticWorker(ctx, taskID, snapshot, tables)
	}
	return r.runRuleWorker(ctx, taskID, workerID, tables, func(table domain.Table) ([]domain.Finding, error) {
		findings, usage, err := r.findingsForWorker(ctx, workerID, snapshot, table)
		if err != nil {
			return nil, err
		}
		if err := r.store.AddWorkerTokenUsage(ctx, taskID, workerID, usage); err != nil {
			return nil, err
		}
		return findings, nil
	})
}

func (r *Runner) runSemanticWorker(ctx context.Context, taskID string, snapshot domain.SchemaSnapshot, tables map[string]domain.Table) error {
	const workerID = "semantic-governance"
	worker, err := r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	if worker.Status == domain.WorkerCompleted {
		return nil
	}
	worker.Status = domain.WorkerRunning
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("semantic batch worker started", "task_id", taskID, "worker_id", workerID, "total_batches", worker.TotalUnits, "max_parallel_batches", r.semanticMaxParallelBatches, "batch_max_columns", r.semanticBatchMaxColumns)

	parallel := r.semanticMaxParallelBatches
	if worker.TotalUnits > 0 && parallel > worker.TotalUnits {
		parallel = worker.TotalUnits
	}
	if parallel < 1 {
		parallel = 1
	}
	errCh := make(chan error, parallel)
	var group sync.WaitGroup
	for index := 0; index < parallel; index++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for {
				if err := r.checkControl(ctx, taskID); err != nil {
					errCh <- err
					return
				}
				key, found, err := r.store.ClaimNextWorkUnit(ctx, taskID, workerID)
				if err != nil {
					errCh <- err
					return
				}
				if !found {
					return
				}
				var batch semanticBatch
				if err := json.Unmarshal([]byte(key), &batch); err != nil || len(batch.Tables) == 0 {
					// Tasks created before semantic batching used a table name as
					// the WorkUnit key. Keep them resumable as one-table batches.
					legacyTable, ok := tables[key]
					if !ok {
						errCh <- errors.New("invalid semantic batch checkpoint")
						return
					}
					batch = semanticBatch{Index: 0, Tables: []string{legacyTable.Name}, ColumnCount: len(legacyTable.Columns)}
				}
				batchTables := make([]domain.Table, 0, len(batch.Tables))
				for _, name := range batch.Tables {
					table, ok := tables[name]
					if !ok {
						errCh <- errors.New("semantic batch table is missing")
						return
					}
					batchTables = append(batchTables, table)
				}
				findings, usage, err := r.findingsForSemanticBatch(ctx, batch, batchTables)
				if err != nil {
					errCh <- err
					return
				}
				if err := r.store.AddWorkerTokenUsage(ctx, taskID, workerID, usage); err != nil {
					errCh <- err
					return
				}
				progressObject := fmt.Sprintf("语义批次 %d（%d 表 / %d 字段）", batch.Index, len(batch.Tables), batch.ColumnCount)
				if err := r.store.CompleteWorkUnitWithProgress(ctx, taskID, workerID, key, progressObject, findings); err != nil {
					errCh <- err
					return
				}
				progress, err := r.store.GetWorker(ctx, taskID, workerID)
				if err != nil {
					errCh <- err
					return
				}
				slog.Info("semantic batch committed", "task_id", taskID, "worker_id", workerID, "batch", batch.Index, "table_count", len(batch.Tables), "column_count", batch.ColumnCount, "completed_batches", progress.CompletedUnits, "total_batches", progress.TotalUnits, "finding_count", progress.FindingCount, "token_total", progress.TokenUsage.TotalTokens)
			}
		}()
	}
	group.Wait()
	close(errCh)
	for workerErr := range errCh {
		if workerErr == nil {
			continue
		}
		if errors.Is(workerErr, errPaused) || errors.Is(workerErr, errCancelled) {
			return workerErr
		}
		_ = r.failWorker(ctx, taskID, workerID, workerFailureSummary(workerID))
		slog.Warn("semantic batch worker failed", "task_id", taskID, "worker_id", workerID, "failure_code", workerFailureCode(workerID))
		return workerErr
	}
	worker, err = r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	worker.Status, worker.CurrentObject, worker.Checkpoint = domain.WorkerCompleted, "", "all_batches_completed"
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("semantic batch worker completed", "task_id", taskID, "worker_id", workerID, "completed_batches", worker.CompletedUnits, "total_batches", worker.TotalUnits, "finding_count", worker.FindingCount, "token_total", worker.TokenUsage.TotalTokens)
	return nil
}

func (r *Runner) findingsForWorker(ctx context.Context, workerID string, snapshot domain.SchemaSnapshot, table domain.Table) ([]domain.Finding, domain.TokenUsage, error) {
	findings := r.rules.Inspect(workerID, snapshot, table)
	if workerID == "index-performance" && table.EstimatedRows >= r.indexAdviceMinRows {
		explanation, err := r.indexAdvisor.AdviseIndex(ctx, table, r.indexAdvicePrompt)
		if err != nil || strings.TrimSpace(explanation.Content) == "" {
			return nil, domain.TokenUsage{}, errors.New("LLM index advice failed")
		}
		findings = append(findings, domain.Finding{ID: fmt.Sprintf("ai-index-%s", table.Name), RuleID: "AI-INDEX-001", Severity: domain.SeveritySuggestion, SeverityReason: "AI 大表索引静态建议", ObjectType: "table", ObjectName: table.Name, Title: "AI 大表索引建议", Evidence: fmt.Sprintf("表估算行数为 %d；建议仅基于静态元数据，不包含查询负载", table.EstimatedRows), Recommendation: explanation.Content, RequiresManualReview: true, SQLRisk: "medium", CreatedAt: time.Now().UTC()})
		return findings, explanation.TokenUsage, nil
	}
	return findings, domain.TokenUsage{}, nil
}

func (r *Runner) findingsForSemanticBatch(ctx context.Context, batch semanticBatch, tables []domain.Table) ([]domain.Finding, domain.TokenUsage, error) {
	findings := make([]domain.Finding, 0)
	usage := domain.TokenUsage{}
	batchID, objectName := fmt.Sprintf("batch-%03d", batch.Index), fmt.Sprintf("语义批次 %d", batch.Index)
	if batch.Index == 0 && len(batch.Tables) == 1 {
		batchID, objectName = "legacy-"+batch.Tables[0], batch.Tables[0]
	}
	for _, rule := range r.semanticRules {
		explanation, err := r.semanticAdvisor.AdviseSchemaBatch(ctx, tables, rule.Prompt)
		if err != nil || strings.TrimSpace(explanation.Content) == "" {
			return nil, usage, errors.New("LLM semantic review failed")
		}
		usage = usage.Add(explanation.TokenUsage)
		if strings.Contains(explanation.Content, "无需要人工复核的语义与治理风险") {
			continue
		}
		findings = append(findings, domain.Finding{ID: fmt.Sprintf("semantic-ai-%s-%s", rule.ID, batchID), RuleID: rule.ID, Severity: domain.SeveritySuggestion, SeverityReason: "AI " + rule.Title, ObjectType: "schema_batch", ObjectName: objectName, Title: "AI " + rule.Title, Evidence: fmt.Sprintf("基于 %d 个表、%d 个字段的静态元数据分析，不包含业务数据", len(batch.Tables), batch.ColumnCount), Recommendation: explanation.Content, RequiresManualReview: true, SQLRisk: "low", CreatedAt: time.Now().UTC()})
	}
	return findings, usage, nil
}

func shouldRunSemanticReview(table domain.Table) bool {
	keywords := []string{"_id", "status", "state", "type", "flag", "email", "mail", "phone", "mobile", "tel", "id_card", "identity", "passport", "bank", "account", "address", "password", "secret", "token", "tenant", "created", "updated", "deleted"}
	for _, column := range table.Columns {
		name := strings.ToLower(column.Name)
		for _, keyword := range keywords {
			if strings.Contains(name, keyword) {
				return true
			}
		}
	}
	return false
}

func (r *Runner) runRuleWorker(ctx context.Context, taskID, workerID string, tables map[string]domain.Table, rule func(domain.Table) ([]domain.Finding, error)) error {
	worker, err := r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	if worker.Status == domain.WorkerCompleted {
		return nil
	}
	worker.Status = domain.WorkerRunning
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("worker execution started", "task_id", taskID, "worker_id", workerID, "worker_name", worker.Name, "total_units", worker.TotalUnits)
	for {
		if err := r.checkControl(ctx, taskID); err != nil {
			return err
		}
		key, found, err := r.store.ClaimNextWorkUnit(ctx, taskID, workerID)
		if err != nil {
			return err
		}
		if !found {
			break
		}
		worker, err = r.store.GetWorker(ctx, taskID, workerID)
		if err != nil {
			return err
		}
		worker.CurrentObject, worker.Status = key, domain.WorkerRunning
		if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
			return err
		}
		table, ok := tables[key]
		if !ok {
			return fmt.Errorf("schema table %s is missing", key)
		}
		findings, err := rule(table)
		if err != nil {
			_ = r.failWorker(ctx, taskID, workerID, workerFailureSummary(workerID))
			slog.Warn("worker rule execution failed", "task_id", taskID, "worker_id", workerID, "failure_code", workerFailureCode(workerID))
			return err
		}
		if err := r.store.CompleteWorkUnit(ctx, taskID, workerID, key, findings); err != nil {
			return err
		}
		progress, err := r.store.GetWorker(ctx, taskID, workerID)
		if err != nil {
			return err
		}
		slog.Info("worker progress committed", "task_id", taskID, "worker_id", workerID, "completed_units", progress.CompletedUnits, "total_units", progress.TotalUnits, "finding_count", progress.FindingCount, "token_total", progress.TokenUsage.TotalTokens, "current_object", redact(key))
	}
	worker, err = r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	worker.Status, worker.CurrentObject, worker.Checkpoint = domain.WorkerCompleted, "", "all_units_completed"
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("worker completed", "task_id", taskID, "worker_id", workerID, "completed_units", worker.CompletedUnits, "total_units", worker.TotalUnits, "finding_count", worker.FindingCount, "token_total", worker.TokenUsage.TotalTokens)
	return nil
}

func (r *Runner) completeSimpleWorker(ctx context.Context, taskID, workerID, checkpoint string) error {
	if err := r.checkControl(ctx, taskID); err != nil {
		return err
	}
	worker, err := r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	if worker.Status == domain.WorkerCompleted {
		return nil
	}
	worker.Status, worker.TotalUnits, worker.CompletedUnits, worker.Checkpoint = domain.WorkerCompleted, 1, 1, checkpoint
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("worker completed", "task_id", taskID, "worker_id", workerID, "completed_units", worker.CompletedUnits, "total_units", worker.TotalUnits)
	return nil
}

func (r *Runner) startSimpleWorker(ctx context.Context, taskID, workerID, checkpoint string) error {
	worker, err := r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	if worker.Status == domain.WorkerCompleted {
		return nil
	}
	worker.Status, worker.TotalUnits, worker.CurrentObject, worker.Checkpoint = domain.WorkerRunning, 1, "worker_results", checkpoint
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Info("worker execution started", "task_id", taskID, "worker_id", workerID, "worker_name", worker.Name, "total_units", worker.TotalUnits)
	return nil
}

func (r *Runner) failWorker(ctx context.Context, taskID, workerID, message string) error {
	worker, err := r.store.GetWorker(ctx, taskID, workerID)
	if err != nil {
		return err
	}
	worker.Status, worker.ErrorSummary, worker.CurrentObject = domain.WorkerFailed, message, ""
	if err := r.store.SetWorker(ctx, taskID, worker); err != nil {
		return err
	}
	slog.Warn("worker execution failed", "task_id", taskID, "worker_id", workerID, "failure_summary", message)
	return nil
}

func workerFailureSummary(workerID string) string {
	switch workerID {
	case "index-performance":
		return "LLM 索引建议生成失败，请检查模型服务配置与可用性"
	case "semantic-governance":
		return "LLM 语义与治理分析失败，请检查模型服务配置与可用性"
	default:
		return "Worker 规则执行失败，请检查服务日志"
	}
}

func workerFailureCode(workerID string) string {
	switch workerID {
	case "index-performance":
		return "index_llm_call_failed"
	case "semantic-governance":
		return "semantic_llm_call_failed"
	default:
		return "rule_evaluation_failed"
	}
}

func (r *Runner) setStage(ctx context.Context, taskID, stage string) error {
	task, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	if task.Status == domain.TaskPending {
		task.Status = domain.TaskRunning
	}
	task.Stage, task.UpdatedAt = stage, time.Now().UTC()
	if err := r.store.SaveTask(ctx, task); err != nil {
		return err
	}
	slog.Info("task stage changed", "task_id", taskID, "stage", stage)
	return nil
}

func (r *Runner) checkControl(ctx context.Context, taskID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	task, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	switch task.Status {
	case domain.TaskCancelled:
		return errCancelled
	case domain.TaskPausing, domain.TaskPaused:
		return errPaused
	}
	return nil
}

func (r *Runner) pauseForConnection(ctx context.Context, taskID string) error {
	task, err := r.store.GetTask(ctx, taskID)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	task.Status, task.ConnectionRequired, task.PausedAt, task.UpdatedAt = domain.TaskPaused, true, &now, now
	task.FailureSummary = "服务重启后需要重新提交连接信息"
	if err := r.store.SaveTask(ctx, task); err != nil {
		return err
	}
	return errPaused
}

func buildReport(taskID string, findings []domain.Finding, unchecked []string) domain.Report {
	summary := map[domain.Severity]int{domain.SeverityUrgent: 0, domain.SeverityNormal: 0, domain.SeveritySuggestion: 0}
	for _, finding := range findings {
		summary[finding.Severity]++
	}
	now := time.Now().UTC()
	report := domain.Report{TaskID: taskID, GeneratedAt: now, RuleVersion: "builtin-0.1", Summary: summary, Findings: findings, UncheckedObject: unchecked}
	var markdown strings.Builder
	fmt.Fprintf(&markdown, "# 数据库体检报告\n\n生成时间：%s\n\n", now.Format(time.RFC3339))
	fmt.Fprintf(&markdown, "- 紧急：%d\n- 一般：%d\n- 建议：%d\n\n", summary[domain.SeverityUrgent], summary[domain.SeverityNormal], summary[domain.SeveritySuggestion])
	for _, finding := range findings {
		fmt.Fprintf(&markdown, "## [%s] %s\n\n- 对象：`%s`\n- 证据：%s\n- 建议：%s\n", finding.Severity, finding.Title, finding.ObjectName, finding.Evidence, finding.Recommendation)
		if finding.DMLSuggestion != "" {
			fmt.Fprintf(&markdown, "\n### DML 修复建议（人工复核）\n\n```sql\n%s\n```\n", finding.DMLSuggestion)
		}
		if finding.DDLSuggestion != "" {
			fmt.Fprintf(&markdown, "\n### DDL 变更建议（人工复核）\n\n```sql\n%s\n```\n", finding.DDLSuggestion)
		}
	}
	if len(unchecked) > 0 {
		fmt.Fprintf(&markdown, "\n## 未检查对象\n\n%s\n", strings.Join(unchecked, "\n"))
	}
	report.Markdown = markdown.String()
	var body strings.Builder
	fmt.Fprintf(&body, "<h1>数据库体检报告</h1><p>生成时间：%s</p><ul><li>紧急：%d</li><li>一般：%d</li><li>建议：%d</li></ul>", html.EscapeString(now.Format(time.RFC3339)), summary[domain.SeverityUrgent], summary[domain.SeverityNormal], summary[domain.SeveritySuggestion])
	for _, finding := range findings {
		fmt.Fprintf(&body, "<article><h2>[%s] %s</h2><p><strong>对象：</strong>%s</p><p><strong>证据：</strong>%s</p><p><strong>建议：</strong>%s</p>", html.EscapeString(string(finding.Severity)), html.EscapeString(finding.Title), html.EscapeString(finding.ObjectName), html.EscapeString(finding.Evidence), html.EscapeString(finding.Recommendation))
		if finding.DMLSuggestion != "" {
			fmt.Fprintf(&body, "<pre><code>%s</code></pre>", html.EscapeString(finding.DMLSuggestion))
		}
		if finding.DDLSuggestion != "" {
			fmt.Fprintf(&body, "<pre><code>%s</code></pre>", html.EscapeString(finding.DDLSuggestion))
		}
		body.WriteString("</article>")
	}
	if len(unchecked) > 0 {
		fmt.Fprintf(&body, "<h2>未检查对象</h2><ul><li>%s</li></ul>", html.EscapeString(strings.Join(unchecked, "</li><li>")))
	}
	report.HTML = "<!doctype html><html><head><meta charset=\"utf-8\"><meta http-equiv=\"Content-Security-Policy\" content=\"default-src 'none'; style-src 'unsafe-inline'\"><title>数据库体检报告</title></head><body>" + body.String() + "</body></html>"
	return report
}

func projectedCompletedWorker(workers []domain.WorkerProgress, workerID, checkpoint string) []domain.WorkerProgress {
	projected := append([]domain.WorkerProgress(nil), workers...)
	for index := range projected {
		if projected[index].WorkerID == workerID {
			projected[index].Status = domain.WorkerCompleted
			projected[index].TotalUnits = 1
			projected[index].CompletedUnits = 1
			projected[index].CurrentObject = ""
			projected[index].Checkpoint = checkpoint
		}
	}
	return projected
}

func redact(value string) string {
	if value == "" {
		return value
	}
	parts := strings.Split(value, ".")
	for index, part := range parts {
		if len(part) <= 2 {
			parts[index] = "**"
		} else {
			parts[index] = part[:1] + strings.Repeat("*", len(part)-2) + part[len(part)-1:]
		}
	}
	return strings.Join(parts, ".")
}

func safeError(err error) string {
	if err == nil {
		return ""
	}
	return "任务执行失败，请检查数据库可达性、只读权限和连接参数。"
}
