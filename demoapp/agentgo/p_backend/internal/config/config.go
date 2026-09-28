package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"agentgo/p_backend/internal/llm"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Security  SecurityConfig  `yaml:"security"`
	TaskStore TaskStoreConfig `yaml:"task_store"`
	Scanner   ScannerConfig   `yaml:"scanner"`
	Execution ExecutionConfig `yaml:"execution"`
	LLM       LLMConfig       `yaml:"llm"`
}

type ServerConfig struct {
	Address                     string `yaml:"address"`
	ReadHeaderTimeoutSeconds    int    `yaml:"read_header_timeout_seconds"`
	ReportRebuildTimeoutSeconds int    `yaml:"report_rebuild_timeout_seconds"`
}
type SecurityConfig struct {
	AllowedOrigins      []string `yaml:"allowed_origins"`
	RequestBodyMaxBytes int64    `yaml:"request_body_max_bytes"`
}
type TaskStoreConfig struct {
	Path         string `yaml:"path"`
	TaskTTLHours int    `yaml:"task_ttl_hours"`
}
type ScannerConfig struct {
	ConnectTimeoutSeconds int `yaml:"connect_timeout_seconds"`
	MaxConcurrentWorkers  int `yaml:"max_concurrent_workers"`
}
type ExecutionConfig struct {
	Enabled      bool     `yaml:"enabled"`
	AllowedHosts []string `yaml:"allowed_hosts"`
}

type LLMConfig struct {
	Enabled                    bool                `yaml:"enabled"`
	Endpoint                   string              `yaml:"endpoint"`
	APIKeyEnv                  string              `yaml:"api_key_env"`
	Model                      string              `yaml:"model"`
	StartupCheckTimeoutSeconds int                 `yaml:"startup_check_timeout_seconds"`
	RequestTimeoutSeconds      int                 `yaml:"request_timeout_seconds"`
	IndexAdvice                IndexAdviceConfig   `yaml:"index_advice"`
	ReportSummary              ReportSummaryConfig `yaml:"report_summary"`
}

type IndexAdviceConfig struct {
	Enabled          bool  `yaml:"enabled"`
	MinEstimatedRows int64 `yaml:"min_estimated_rows"`
}

type ReportSummaryConfig struct {
	Enabled bool `yaml:"enabled"`
}

type RuleAndPrompt struct {
	Version string                         `yaml:"version"`
	Workers map[string]WorkerRuleAndPrompt `yaml:"workers"`
	AIRules map[string]AIRule              `yaml:"ai_rules"`
}

type WorkerRuleAndPrompt struct {
	Enabled            bool     `yaml:"enabled"`
	Rules              []string `yaml:"rules"`
	BatchMaxColumns    int      `yaml:"batch_max_columns"`
	MaxParallelBatches int      `yaml:"max_parallel_batches"`
	MaxInputFindings   int      `yaml:"max_input_findings"`
}

type AIRule struct {
	Title  string `yaml:"title"`
	Prompt string `yaml:"prompt"`
}

type AIRuleSpec struct {
	ID     string
	Title  string
	Prompt string
}

var expectedWorkerRules = map[string]map[string]struct{}{
	"schema-collector":    {},
	"naming-convention":   {"NAME-001": {}, "NAME-002": {}, "NAME-003": {}},
	"index-performance":   {"INDEX-001": {}, "INDEX-002": {}, "INDEX-003": {}, "INDEX-004": {}, "INDEX-005": {}, "AI-INDEX-001": {}},
	"type-constraint":     {"TYPE-001": {}, "TYPE-002": {}, "TYPE-003": {}, "TYPE-004": {}},
	"semantic-governance": {"AI-SEMANTIC-001": {}},
	"finding-normalizer":  {},
	"ddl-generator":       {"AI-DDL-001": {}},
	"report-generator":    {"AI-REPORT-001": {}},
}

var expectedAIRules = map[string]struct{}{"AI-INDEX-001": {}, "AI-SEMANTIC-001": {}, "AI-DDL-001": {}, "AI-REPORT-001": {}}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	config := Default()
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	if config.LLM.IndexAdvice.MinEstimatedRows < 1 {
		config.LLM.IndexAdvice.MinEstimatedRows = 100_000
	}
	if config.LLM.StartupCheckTimeoutSeconds < 1 {
		config.LLM.StartupCheckTimeoutSeconds = 15
	}
	if config.LLM.RequestTimeoutSeconds < 1 {
		config.LLM.RequestTimeoutSeconds = 60
	}
	if config.Server.Address == "" {
		config.Server.Address = ":8080"
	}
	if config.Server.ReadHeaderTimeoutSeconds < 1 {
		config.Server.ReadHeaderTimeoutSeconds = 5
	}
	if config.Server.ReportRebuildTimeoutSeconds < 1 {
		config.Server.ReportRebuildTimeoutSeconds = 210
	}
	if config.Security.RequestBodyMaxBytes < 1 {
		config.Security.RequestBodyMaxBytes = 1 << 20
	}
	if config.TaskStore.Path == "" {
		config.TaskStore.Path = "./data/agentgo.db"
	}
	if config.TaskStore.TaskTTLHours < 1 {
		config.TaskStore.TaskTTLHours = 24
	}
	if config.Scanner.ConnectTimeoutSeconds < 1 {
		config.Scanner.ConnectTimeoutSeconds = 10
	}
	if config.Scanner.MaxConcurrentWorkers < 1 {
		config.Scanner.MaxConcurrentWorkers = 3
	}
	if len(config.Security.AllowedOrigins) == 0 {
		config.Security.AllowedOrigins = []string{"http://localhost:5173"}
	}
	if err := config.Validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func Default() Config {
	return Config{Server: ServerConfig{Address: ":8080", ReadHeaderTimeoutSeconds: 5, ReportRebuildTimeoutSeconds: 210}, Security: SecurityConfig{AllowedOrigins: []string{"http://localhost:5173"}, RequestBodyMaxBytes: 1 << 20}, TaskStore: TaskStoreConfig{Path: "./data/agentgo.db", TaskTTLHours: 24}, Scanner: ScannerConfig{ConnectTimeoutSeconds: 10, MaxConcurrentWorkers: 3}, LLM: LLMConfig{StartupCheckTimeoutSeconds: 15, RequestTimeoutSeconds: 60, IndexAdvice: IndexAdviceConfig{MinEstimatedRows: 100_000}}}
}

func (c LLMConfig) Explainer() (llm.Explainer, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	apiKey := os.Getenv(c.APIKeyEnv)
	return llm.NewOpenAICompatible(llm.OpenAICompatibleConfig{Endpoint: c.Endpoint, APIKey: apiKey, Model: c.Model, RequestTimeout: time.Duration(c.RequestTimeoutSeconds) * time.Second})
}

// Validate performs the startup self-check for every configuration domain. It
// intentionally verifies LLM credentials locally but never logs their values.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Server.Address) == "" {
		return errors.New("server.address is required")
	}
	if c.Server.ReadHeaderTimeoutSeconds < 1 {
		return errors.New("server.read_header_timeout_seconds must be positive")
	}
	if c.Server.ReportRebuildTimeoutSeconds < 1 {
		return errors.New("server.report_rebuild_timeout_seconds must be positive")
	}
	if c.Security.RequestBodyMaxBytes < 1 {
		return errors.New("security.request_body_max_bytes must be positive")
	}
	for _, origin := range c.Security.AllowedOrigins {
		parsed, err := url.ParseRequestURI(origin)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return fmt.Errorf("security.allowed_origins contains an invalid origin")
		}
	}
	if strings.TrimSpace(c.TaskStore.Path) == "" || c.TaskStore.TaskTTLHours < 1 {
		return errors.New("task_store.path and task_store.task_ttl_hours must be configured")
	}
	if c.Scanner.ConnectTimeoutSeconds < 1 || c.Scanner.MaxConcurrentWorkers < 1 {
		return errors.New("scanner timeouts and worker concurrency must be positive")
	}
	if c.Execution.Enabled && len(c.Execution.AllowedHosts) == 0 {
		return errors.New("execution.allowed_hosts is required when execution is enabled")
	}
	for _, host := range c.Execution.AllowedHosts {
		if strings.TrimSpace(host) == "" {
			return errors.New("execution.allowed_hosts must not contain empty values")
		}
	}
	return c.LLM.Validate()
}

func (c LLMConfig) Validate() error {
	if !c.Enabled {
		return errors.New("llm.enabled must be true: LLM is required")
	}
	endpoint, err := url.ParseRequestURI(strings.TrimSpace(c.Endpoint))
	if err != nil || (endpoint.Scheme != "http" && endpoint.Scheme != "https") || endpoint.Host == "" {
		return errors.New("llm.endpoint must be an absolute http(s) URL")
	}
	if strings.TrimSpace(c.APIKeyEnv) == "" {
		return errors.New("llm.api_key_env is required")
	}
	if strings.TrimSpace(os.Getenv(c.APIKeyEnv)) == "" {
		return fmt.Errorf("llm api key environment variable %q is empty", c.APIKeyEnv)
	}
	if strings.TrimSpace(c.Model) == "" {
		return errors.New("llm.model is required")
	}
	if c.StartupCheckTimeoutSeconds < 1 || c.RequestTimeoutSeconds < 1 {
		return errors.New("llm startup_check_timeout_seconds and request_timeout_seconds must be positive")
	}
	if !c.IndexAdvice.Enabled || c.IndexAdvice.MinEstimatedRows < 1 {
		return errors.New("llm.index_advice must be enabled with a positive min_estimated_rows")
	}
	if !c.ReportSummary.Enabled {
		return errors.New("llm.report_summary must be enabled")
	}
	return nil
}

// LoadWithRuleAndPrompt loads and cross-validates the main configuration and
// the Worker rule/prompt catalog used by the workflow.
func LoadWithRuleAndPrompt(path string) (Config, RuleAndPrompt, error) {
	config, err := Load(path)
	if err != nil {
		return Config{}, RuleAndPrompt{}, err
	}
	rulePath := filepath.Join(filepath.Dir(path), "rule_and_prompt.yaml")
	rules, err := LoadRuleAndPrompt(rulePath)
	if err != nil {
		return Config{}, RuleAndPrompt{}, err
	}
	return config, rules, nil
}

func LoadRuleAndPrompt(path string) (RuleAndPrompt, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RuleAndPrompt{}, err
	}
	var rules RuleAndPrompt
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return RuleAndPrompt{}, err
	}
	if err := rules.Validate(); err != nil {
		return RuleAndPrompt{}, err
	}
	return rules, nil
}

func (r RuleAndPrompt) Validate() error {
	if strings.TrimSpace(r.Version) == "" {
		return errors.New("rule_and_prompt.version is required")
	}
	for workerID, supportedRules := range expectedWorkerRules {
		worker, ok := r.Workers[workerID]
		if !ok {
			return fmt.Errorf("rule_and_prompt.workers.%s is required", workerID)
		}
		if !worker.Enabled {
			return fmt.Errorf("rule_and_prompt.workers.%s must be enabled", workerID)
		}
		seen := make(map[string]struct{}, len(worker.Rules))
		for _, ruleID := range worker.Rules {
			if _, duplicate := seen[ruleID]; duplicate {
				return fmt.Errorf("rule_and_prompt.workers.%s has duplicate rule %q", workerID, ruleID)
			}
			seen[ruleID] = struct{}{}
			if _, supported := supportedRules[ruleID]; !supported {
				return fmt.Errorf("rule_and_prompt.workers.%s contains unsupported rule %q", workerID, ruleID)
			}
		}
		if len(seen) != len(supportedRules) {
			return fmt.Errorf("rule_and_prompt.workers.%s must declare its complete rule set", workerID)
		}
	}
	for workerID := range r.Workers {
		if _, expected := expectedWorkerRules[workerID]; !expected {
			return fmt.Errorf("rule_and_prompt contains unknown worker %q", workerID)
		}
	}
	for ruleID := range expectedAIRules {
		rule, ok := r.AIRules[ruleID]
		if !ok || strings.TrimSpace(rule.Prompt) == "" || strings.TrimSpace(rule.Title) == "" {
			return fmt.Errorf("rule_and_prompt.ai_rules.%s title and prompt are required", ruleID)
		}
	}
	for ruleID := range r.AIRules {
		if _, expected := expectedAIRules[ruleID]; !expected {
			return fmt.Errorf("rule_and_prompt contains unknown AI rule %q", ruleID)
		}
	}
	semantic := r.Workers["semantic-governance"]
	if semantic.BatchMaxColumns < 1 || semantic.MaxParallelBatches < 1 {
		return errors.New("rule_and_prompt.workers.semantic-governance batch limits must be positive")
	}
	if r.Workers["ddl-generator"].MaxInputFindings < 1 {
		return errors.New("rule_and_prompt.workers.ddl-generator.max_input_findings must be positive")
	}
	return nil
}

func (r RuleAndPrompt) RuleSets() map[string][]string {
	sets := make(map[string][]string, len(r.Workers))
	for workerID, worker := range r.Workers {
		sets[workerID] = append([]string(nil), worker.Rules...)
	}
	return sets
}

func (r RuleAndPrompt) PromptForAIRule(ruleID string) string {
	return r.AIRules[ruleID].Prompt
}

func (r RuleAndPrompt) AIRulesForWorker(workerID string) []AIRuleSpec {
	worker := r.Workers[workerID]
	rules := make([]AIRuleSpec, 0)
	for _, ruleID := range worker.Rules {
		rule, ok := r.AIRules[ruleID]
		if ok {
			rules = append(rules, AIRuleSpec{ID: ruleID, Title: rule.Title, Prompt: rule.Prompt})
		}
	}
	return rules
}

func (r RuleAndPrompt) BatchMaxColumnsFor(workerID string) int {
	return r.Workers[workerID].BatchMaxColumns
}

func (r RuleAndPrompt) MaxParallelBatchesFor(workerID string) int {
	return r.Workers[workerID].MaxParallelBatches
}

func (r RuleAndPrompt) MaxInputFindingsFor(workerID string) int {
	return r.Workers[workerID].MaxInputFindings
}
