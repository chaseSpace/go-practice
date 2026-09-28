package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"agentgo/p_backend/internal/config"
	"agentgo/p_backend/internal/execution"
	"agentgo/p_backend/internal/httpapi"
	"agentgo/p_backend/internal/rules"
	"agentgo/p_backend/internal/scanner"
	"agentgo/p_backend/internal/service"
	"agentgo/p_backend/internal/store"
	"agentgo/p_backend/internal/workflow"
)

func main() {
	// Include the originating Go file and line in every structured log, including
	// configuration/LLM startup failures that occur before the HTTP server runs.
	slog.SetDefault(newLogger())
	configPath := valueOrDefault(os.Getenv("AGENTGO_CONFIG"), "./config/config.yaml")
	appConfig, ruleAndPrompt, err := config.LoadWithRuleAndPrompt(configPath)
	if err != nil {
		slog.Error("backend configuration self-check failed", "path", configPath, "error", err)
		os.Exit(1)
	}
	address := appConfig.Server.Address
	databasePath := appConfig.TaskStore.Path
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o750); err != nil {
		slog.Error("create data directory", "error", err)
		os.Exit(1)
	}
	taskStore, err := store.Open(databasePath)
	if err != nil {
		slog.Error("open task store", "error", err)
		os.Exit(1)
	}
	defer taskStore.Close()
	explainer, err := appConfig.LLM.Explainer()
	if err != nil {
		slog.Error("LLM startup self-check failed", "error", err)
		os.Exit(1)
	}
	llmCheckContext, cancelLLMCheck := context.WithTimeout(context.Background(), time.Duration(appConfig.LLM.StartupCheckTimeoutSeconds)*time.Second)
	err = explainer.CheckAvailability(llmCheckContext)
	cancelLLMCheck()
	if err != nil {
		slog.Error("LLM startup reachability check failed", "error", err)
		os.Exit(1)
	}
	llmEndpoint, _ := url.Parse(appConfig.LLM.Endpoint)
	slog.Info("LLM connection verified", "model", appConfig.LLM.Model, "endpoint_host", llmEndpoint.Host)
	semanticRules := make([]workflow.SemanticRule, 0)
	for _, rule := range ruleAndPrompt.AIRulesForWorker("semantic-governance") {
		semanticRules = append(semanticRules, workflow.SemanticRule{ID: rule.ID, Title: rule.Title, Prompt: rule.Prompt})
	}
	runnerOptions := []workflow.Option{
		workflow.WithRuleEngine(rules.NewDefault(rules.WithWorkerRuleSets(ruleAndPrompt.RuleSets()))),
		workflow.WithIndexAdvisor(explainer),
		workflow.WithIndexAdviceMinRows(appConfig.LLM.IndexAdvice.MinEstimatedRows),
		workflow.WithIndexAdvicePrompt(ruleAndPrompt.PromptForAIRule("AI-INDEX-001")),
		workflow.WithSemanticAdvisor(explainer),
		workflow.WithSemanticRules(semanticRules),
		workflow.WithSemanticBatchMaxColumns(ruleAndPrompt.BatchMaxColumnsFor("semantic-governance")),
		workflow.WithSemanticMaxParallelBatches(ruleAndPrompt.MaxParallelBatchesFor("semantic-governance")),
		workflow.WithDDLAdvisor(explainer),
		workflow.WithDDLPrompt(ruleAndPrompt.PromptForAIRule("AI-DDL-001")),
		workflow.WithDDLMaxInputFindings(ruleAndPrompt.MaxInputFindingsFor("ddl-generator")),
		workflow.WithReportSummarizer(explainer),
		workflow.WithReportPrompt(ruleAndPrompt.PromptForAIRule("AI-REPORT-001")),
		workflow.WithRuleVersion(ruleAndPrompt.Version),
	}
	connectTimeout := time.Duration(appConfig.Scanner.ConnectTimeoutSeconds) * time.Second
	runner := workflow.New(taskStore, scanner.MultiCollector{MySQL: scanner.MySQLCollector{ConnectTimeout: connectTimeout}, PostgreSQL: scanner.PostgreSQLCollector{ConnectTimeout: connectTimeout}}, runnerOptions...)
	tasks := service.New(taskStore, runner, time.Duration(appConfig.TaskStore.TaskTTLHours)*time.Hour)
	allowedHosts := map[string]struct{}{}
	for _, host := range appConfig.Execution.AllowedHosts {
		if host = strings.TrimSpace(host); host != "" {
			allowedHosts[host] = struct{}{}
		}
	}
	executor := execution.New(execution.Config{Enabled: appConfig.Execution.Enabled, AllowedHosts: allowedHosts})
	server := &http.Server{Addr: address, Handler: httpapi.New(taskStore, tasks, httpapi.WithExplainer(explainer), httpapi.WithExecutor(executor), httpapi.WithAllowedOrigins(appConfig.Security.AllowedOrigins), httpapi.WithReportRebuildTimeout(time.Duration(appConfig.Server.ReportRebuildTimeoutSeconds)*time.Second)), ReadHeaderTimeout: time.Duration(appConfig.Server.ReadHeaderTimeoutSeconds) * time.Second}

	go func() {
		slog.Info("agentgo backend listening", "address", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
	}
}

func newLogger() *slog.Logger {
	// main.go is under <project-root>/p_backend/cmd/server. Resolve the root from
	// the compiled source path so log locations stay portable across machines.
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
	replaceAttr := func(_ []string, attr slog.Attr) slog.Attr {
		if attr.Key != slog.SourceKey {
			return attr
		}
		source, ok := attr.Value.Any().(*slog.Source)
		if !ok || source == nil {
			return attr
		}
		relativePath, err := filepath.Rel(projectRoot, source.File)
		if err != nil || relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
			relativePath = filepath.Base(source.File)
		}
		return slog.String(slog.SourceKey, filepath.ToSlash(relativePath)+":"+strconv.Itoa(source.Line))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{AddSource: true, ReplaceAttr: replaceAttr}))
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
