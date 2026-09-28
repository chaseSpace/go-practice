package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"agentgo/p_backend/internal/domain"
	"agentgo/p_backend/internal/execution"
	"agentgo/p_backend/internal/llm"
	"agentgo/p_backend/internal/service"
	"agentgo/p_backend/internal/store"
)

type API struct {
	store                *store.Store
	service              *service.Service
	explainer            llm.Explainer
	executor             *execution.Controller
	allowedOrigins       map[string]struct{}
	reportRebuildTimeout time.Duration
	rebuildMu            sync.Mutex
}

type ExecuteRequest struct {
	Confirmation  string                   `json:"confirmation"`
	ConnectionURL string                   `json:"connectionUrl"`
	Connection    *domain.ConnectionConfig `json:"connection"`
}
type Option func(*API)

func New(store *store.Store, service *service.Service, options ...Option) http.Handler {
	api := &API{store: store, service: service, explainer: llm.DisabledExplainer{}, executor: execution.New(execution.Config{}), allowedOrigins: map[string]struct{}{"http://localhost:5173": {}}, reportRebuildTimeout: 90 * time.Second}
	for _, option := range options {
		option(api)
	}
	return api.withCORS(http.HandlerFunc(api.routes))
}

func WithExplainer(explainer llm.Explainer) Option {
	return func(api *API) {
		if explainer != nil {
			api.explainer = explainer
		}
	}
}

func WithExecutor(executor *execution.Controller) Option {
	return func(api *API) {
		if executor != nil {
			api.executor = executor
		}
	}
}

func WithAllowedOrigins(origins []string) Option {
	return func(api *API) {
		api.allowedOrigins = map[string]struct{}{}
		for _, origin := range origins {
			if origin != "" {
				api.allowedOrigins[origin] = struct{}{}
			}
		}
	}
}

// WithReportRebuildTimeout configures the bounded background context used by
// report rebuilding. The work survives a browser-side request cancellation so
// its LLM summary and persisted report are not abandoned halfway through.
func WithReportRebuildTimeout(timeout time.Duration) Option {
	return func(api *API) {
		if timeout > 0 {
			api.reportRebuildTimeout = timeout
		}
	}
}

func (api *API) routes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/healthz" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1")
	if path == r.URL.Path {
		writeError(w, http.StatusNotFound, "not_found", "接口不存在")
		return
	}
	if path == "/tasks" && r.Method == http.MethodPost {
		api.createTask(w, r)
		return
	}
	segments := splitPath(path)
	if len(segments) < 2 {
		writeError(w, http.StatusNotFound, "not_found", "接口不存在")
		return
	}
	resource, taskID := segments[0], segments[1]
	if resource == "tasks" {
		api.taskRoutes(w, r, taskID, segments[2:])
		return
	}
	if resource == "reports" {
		api.reportRoutes(w, r, taskID, segments[2:])
		return
	}
	writeError(w, http.StatusNotFound, "not_found", "接口不存在")
}

func (api *API) createTask(w http.ResponseWriter, r *http.Request) {
	var request domain.CreateTaskRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	response, err := api.service.Create(r.Context(), request)
	if err != nil {
		if errors.Is(err, service.ErrInvalidConnection) {
			writeError(w, http.StatusBadRequest, "invalid_connection", "请填写有效的 MySQL 连接信息")
			return
		}
		writeError(w, http.StatusInternalServerError, "create_task_failed", "创建任务失败")
		return
	}
	_ = api.store.AddAuditEvent(r.Context(), domain.AuditEvent{TaskID: response.TaskID, Action: "task.create", Outcome: "success"})
	writeJSON(w, http.StatusCreated, response)
}

func (api *API) taskRoutes(w http.ResponseWriter, r *http.Request, taskID string, tail []string) {
	if !api.authorize(w, r, taskID) {
		return
	}
	if len(tail) == 0 && r.Method == http.MethodGet {
		task, err := api.store.GetTask(r.Context(), taskID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "task_read_failed", "读取任务失败")
			return
		}
		writeJSON(w, http.StatusOK, task)
		return
	}
	if len(tail) == 1 && tail[0] == "audit" && r.Method == http.MethodGet {
		events, err := api.store.ListAuditEvents(r.Context(), taskID, positiveInt(r.URL.Query().Get("limit"), 50, 1, 200))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "audit_read_failed", "读取审计事件失败")
			return
		}
		writeJSON(w, http.StatusOK, events)
		return
	}
	if len(tail) != 1 || r.Method != http.MethodPost {
		writeError(w, http.StatusNotFound, "not_found", "接口不存在")
		return
	}
	var task domain.Task
	var err error
	switch tail[0] {
	case "pause":
		task, err = api.service.Pause(r.Context(), taskID)
	case "resume":
		var request domain.ResumeTaskRequest
		if !decodeJSON(w, r, &request) {
			return
		}
		task, err = api.service.Resume(r.Context(), taskID, request)
	case "cancel":
		task, err = api.service.Cancel(r.Context(), taskID)
	default:
		writeError(w, http.StatusNotFound, "not_found", "接口不存在")
		return
	}
	if err != nil {
		if errors.Is(err, service.ErrInvalidState) {
			writeError(w, http.StatusConflict, "invalid_task_state", "任务当前状态不支持该操作")
			return
		}
		if strings.Contains(err.Error(), "connection required") {
			writeJSON(w, http.StatusConflict, map[string]any{"code": "connection_required", "message": "请重新提交连接信息后继续任务", "connectionRequired": true})
			return
		}
		if errors.Is(err, service.ErrInvalidConnection) {
			writeError(w, http.StatusBadRequest, "invalid_connection", "请填写有效的 MySQL 连接信息")
			return
		}
		writeError(w, http.StatusInternalServerError, "task_action_failed", "任务操作失败")
		return
	}
	_ = api.store.AddAuditEvent(r.Context(), domain.AuditEvent{TaskID: taskID, Action: "task." + tail[0], Outcome: "success"})
	writeJSON(w, http.StatusOK, task)
}

func (api *API) reportRoutes(w http.ResponseWriter, r *http.Request, taskID string, tail []string) {
	if !api.authorize(w, r, taskID) {
		return
	}
	report, err := api.store.GetReport(r.Context(), taskID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusConflict, "report_not_ready", "报告尚未生成")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "report_read_failed", "读取报告失败")
		return
	}
	if len(tail) == 1 && (tail[0] == "rebuild" || tail[0] == "reanalyze") {
		reanalyze := tail[0] == "reanalyze"
		auditAction := "report.rebuild"
		if reanalyze {
			auditAction = "report.reanalyze"
		}
		switch r.Method {
		case http.MethodGet:
			status, err := api.store.GetReportRebuildStatus(r.Context(), taskID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "report_rebuild_status_failed", "读取报告重建状态失败")
				return
			}
			writeJSON(w, http.StatusOK, status)
			return
		case http.MethodPost:
			status, scheduled, err := api.store.RequestReportRebuild(r.Context(), taskID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "report_rebuild_request_failed", "提交报告重建失败")
				return
			}
			if scheduled {
				_ = api.store.AddAuditEvent(r.Context(), domain.AuditEvent{TaskID: taskID, Action: auditAction, Outcome: "requested"})
				go api.runReportRebuild(taskID, reanalyze)
			}
			writeJSON(w, http.StatusAccepted, status)
			return
		default:
			writeError(w, http.StatusNotFound, "not_found", "接口不存在")
			return
		}
	}
	if len(tail) == 0 && r.Method == http.MethodGet {
		severity := domain.Severity(r.URL.Query().Get("severity"))
		keyword := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("keyword")))
		filtered := make([]domain.Finding, 0, len(report.Findings))
		for _, finding := range report.Findings {
			if severity != "" && finding.Severity != severity {
				continue
			}
			if keyword != "" && !strings.Contains(strings.ToLower(finding.Title+finding.ObjectName+finding.RuleID), keyword) {
				continue
			}
			filtered = append(filtered, finding)
		}
		page := positiveInt(r.URL.Query().Get("page"), 1, 1, 10_000)
		pageSize := positiveInt(r.URL.Query().Get("page_size"), 20, 1, 100)
		total := len(filtered)
		start := (page - 1) * pageSize
		if start > total {
			start = total
		}
		end := start + pageSize
		if end > total {
			end = total
		}
		report.Findings = filtered[start:end]
		comparison := domain.BaselineComparison{NewFindingIDs: []string{}}
		if baselineName := r.URL.Query().Get("baseline"); baselineName != "" {
			if baselineTaskID, baselineErr := api.store.GetBaseline(r.Context(), baselineName); baselineErr == nil {
				comparison.BaselineTaskID = baselineTaskID
				if baselineReport, reportErr := api.store.GetReport(r.Context(), baselineTaskID); reportErr == nil {
					known := map[string]struct{}{}
					for _, finding := range baselineReport.Findings {
						known[finding.RuleID+"|"+finding.ObjectName] = struct{}{}
					}
					for _, finding := range report.Findings {
						if _, ok := known[finding.RuleID+"|"+finding.ObjectName]; !ok {
							comparison.NewFindingIDs = append(comparison.NewFindingIDs, finding.ID)
						}
					}
				}
			}
		}
		writeJSON(w, http.StatusOK, struct {
			domain.ReportPage
			Baseline domain.BaselineComparison `json:"baseline"`
		}{ReportPage: domain.ReportPage{Report: report, Total: total, Page: page, PageSize: pageSize}, Baseline: comparison})
		return
	}
	if len(tail) == 1 && tail[0] == "baseline" && r.Method == http.MethodPost {
		var request struct {
			Name string `json:"name"`
		}
		if !decodeJSON(w, r, &request) {
			return
		}
		if request.Name == "" {
			request.Name = "default"
		}
		if err := api.store.SetBaseline(r.Context(), request.Name, taskID); err != nil {
			writeError(w, http.StatusInternalServerError, "baseline_save_failed", "保存基线失败")
			return
		}
		_ = api.store.AddAuditEvent(r.Context(), domain.AuditEvent{TaskID: taskID, Action: "report.baseline", Outcome: "success", Detail: request.Name})
		writeJSON(w, http.StatusOK, map[string]string{"name": request.Name, "taskId": taskID})
		return
	}
	if len(tail) == 1 && tail[0] == "preview" && r.Method == http.MethodGet {
		api.writePreview(w, r, report)
		return
	}
	if len(tail) == 1 && tail[0] == "sql" && r.Method == http.MethodGet {
		api.writeSQL(w, report)
		return
	}
	if len(tail) == 1 && tail[0] == "export" && r.Method == http.MethodGet {
		api.writeExport(w, r, report)
		return
	}
	if len(tail) == 2 && tail[0] == "findings" && r.Method == http.MethodGet {
		for _, finding := range report.Findings {
			if finding.ID == tail[1] {
				writeJSON(w, http.StatusOK, finding)
				return
			}
		}
		writeError(w, http.StatusNotFound, "finding_not_found", "问题不存在")
		return
	}
	if len(tail) == 3 && tail[0] == "findings" && tail[2] == "explain" && r.Method == http.MethodPost {
		for _, finding := range report.Findings {
			if finding.ID != tail[1] {
				continue
			}
			explanation, err := api.explainer.Explain(r.Context(), finding)
			if errors.Is(err, llm.ErrDisabled) {
				writeError(w, http.StatusConflict, "llm_disabled", "LLM 辅助解释未配置或未启用")
				return
			}
			if err != nil {
				writeError(w, http.StatusBadGateway, "llm_failed", "LLM 辅助解释暂时不可用")
				return
			}
			if err := api.store.AddWorkerTokenUsage(r.Context(), taskID, "report-generator", explanation.TokenUsage); err != nil {
				writeError(w, http.StatusInternalServerError, "token_usage_persist_failed", "记录模型用量失败")
				return
			}
			_ = api.store.AddAuditEvent(r.Context(), domain.AuditEvent{TaskID: taskID, Action: "finding.explain", Outcome: "success", Detail: finding.RuleID})
			writeJSON(w, http.StatusOK, explanation)
			return
		}
		writeError(w, http.StatusNotFound, "finding_not_found", "问题不存在")
		return
	}
	if len(tail) == 3 && tail[0] == "findings" && tail[2] == "execute" && r.Method == http.MethodPost {
		var request ExecuteRequest
		if !decodeJSON(w, r, &request) {
			return
		}
		for _, finding := range report.Findings {
			if finding.ID != tail[1] {
				continue
			}
			result, err := api.executor.Execute(r.Context(), finding, request.Confirmation, request.ConnectionURL, request.Connection)
			if errors.Is(err, execution.ErrDisabled) {
				writeError(w, http.StatusConflict, "execution_disabled", "受控执行未启用")
				return
			}
			if errors.Is(err, execution.ErrConfirmation) || errors.Is(err, execution.ErrTargetNotAllowed) || errors.Is(err, execution.ErrUnsafeStatement) {
				writeError(w, http.StatusBadRequest, "execution_rejected", "执行请求未通过安全校验")
				return
			}
			if err != nil {
				writeError(w, http.StatusBadGateway, "execution_failed", "执行失败，目标数据库未被自动重试")
				return
			}
			_ = api.store.AddAuditEvent(r.Context(), domain.AuditEvent{TaskID: taskID, Action: "finding.execute", Outcome: "success", Detail: finding.RuleID})
			writeJSON(w, http.StatusOK, result)
			return
		}
		writeError(w, http.StatusNotFound, "finding_not_found", "问题不存在")
		return
	}
	writeError(w, http.StatusNotFound, "not_found", "接口不存在")
}

func (api *API) runReportRebuild(taskID string, reanalyze bool) {
	// This is intentionally detached from the request context. Its state and
	// output are persisted, while the configured deadline still bounds LLM work.
	api.rebuildMu.Lock()
	defer api.rebuildMu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), api.reportRebuildTimeout)
	defer cancel()
	if err := api.store.MarkReportRebuildRunning(ctx, taskID); err != nil {
		slog.Error("report rebuild could not start", "task_id", taskID, "error", err)
		return
	}
	mode, auditAction := "saved_findings", "report.rebuild"
	if reanalyze {
		mode, auditAction = "schema_reanalysis", "report.reanalyze"
	}
	slog.Info("report rebuild started", "task_id", taskID, "mode", mode, "timeout", api.reportRebuildTimeout)
	var err error
	if reanalyze {
		_, err = api.service.ReanalyzeReport(ctx, taskID)
	} else {
		_, err = api.service.RegenerateReport(ctx, taskID)
	}
	if err != nil {
		failure := "报告重建失败，请检查 LLM 服务与后端日志后重试"
		if errors.Is(err, context.DeadlineExceeded) {
			failure = "报告重建超时，请检查模型服务响应后重试"
		}
		if saveErr := api.store.CompleteReportRebuild(context.Background(), taskID, failure); saveErr != nil {
			slog.Error("report rebuild failure state could not be saved", "task_id", taskID, "error", saveErr)
		}
		_ = api.store.AddAuditEvent(context.Background(), domain.AuditEvent{TaskID: taskID, Action: auditAction, Outcome: "failed", Detail: failure})
		slog.Error("report rebuild failed", "task_id", taskID, "mode", mode, "error", err)
		return
	}
	if err := api.store.CompleteReportRebuild(context.Background(), taskID, ""); err != nil {
		slog.Error("report rebuild success state could not be saved", "task_id", taskID, "error", err)
		return
	}
	_ = api.store.AddAuditEvent(context.Background(), domain.AuditEvent{TaskID: taskID, Action: auditAction, Outcome: "success"})
	slog.Info("report rebuild completed", "task_id", taskID, "mode", mode)
}

func (api *API) writePreview(w http.ResponseWriter, r *http.Request, report domain.Report) {
	format := r.URL.Query().Get("format")
	if format == "markdown" {
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Header().Set("Content-Disposition", "inline; filename=report.md")
		_, _ = w.Write([]byte(report.Markdown))
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'")
	w.Header().Set("Content-Disposition", "inline; filename=report.html")
	_, _ = w.Write([]byte(report.HTML))
}

func (api *API) writeSQL(w http.ResponseWriter, report domain.Report) {
	var sql strings.Builder
	for _, finding := range report.Findings {
		if finding.DMLSuggestion == "" && finding.DDLSuggestion == "" {
			continue
		}
		sql.WriteString("-- ")
		sql.WriteString(finding.Title)
		sql.WriteString(" (需要人工复核)\n")
		if finding.PrecheckSQL != "" {
			sql.WriteString(finding.PrecheckSQL)
			sql.WriteString("\n")
		}
		if finding.DMLSuggestion != "" {
			sql.WriteString(finding.DMLSuggestion)
			sql.WriteString("\n")
		}
		if finding.DDLSuggestion != "" {
			sql.WriteString(finding.DDLSuggestion)
		}
		sql.WriteString("\n\n")
	}
	w.Header().Set("Content-Type", "text/sql; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=database-suggestions.sql")
	_, _ = w.Write([]byte(sql.String()))
}

func (api *API) writeExport(w http.ResponseWriter, r *http.Request, report domain.Report) {
	switch r.URL.Query().Get("format") {
	case "json":
		w.Header().Set("Content-Disposition", "attachment; filename=report.json")
		writeJSON(w, http.StatusOK, report)
	case "markdown":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=report.md")
		_, _ = w.Write([]byte(report.Markdown))
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=report.html")
		_, _ = w.Write([]byte(report.HTML))
	}
}

func (api *API) authorize(w http.ResponseWriter, r *http.Request, taskID string) bool {
	err := api.store.Authorize(r.Context(), taskID, r.Header.Get("X-Task-Access-Token"))
	if err == nil {
		return true
	}
	if errors.Is(err, store.ErrNotFound) || errors.Is(err, store.ErrUnauthorized) {
		writeError(w, http.StatusNotFound, "task_unavailable", "任务已过期或不可访问")
		return false
	}
	writeError(w, http.StatusInternalServerError, "authorization_failed", "任务访问失败")
	return false
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "请求参数格式错误")
		return false
	}
	return true
}

func splitPath(value string) []string {
	parts := strings.Split(strings.Trim(value, "/"), "/")
	if len(parts) == 1 && parts[0] == "" {
		return nil
	}
	return parts
}

func positiveInt(value string, fallback, min, max int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < min {
		return fallback
	}
	if parsed > max {
		return max
	}
	return parsed
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"code": code, "message": message})
}

func (api *API) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if origin := r.Header.Get("Origin"); origin != "" {
			if _, ok := api.allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Task-Access-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
