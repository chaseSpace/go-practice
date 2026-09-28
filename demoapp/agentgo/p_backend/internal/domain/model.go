package domain

import "time"

type TaskStatus string

const (
	TaskPending        TaskStatus = "pending"
	TaskRunning        TaskStatus = "running"
	TaskPausing        TaskStatus = "pausing"
	TaskPaused         TaskStatus = "paused"
	TaskPartialSuccess TaskStatus = "partial_success"
	TaskFailed         TaskStatus = "failed"
	TaskCompleted      TaskStatus = "completed"
	TaskCancelled      TaskStatus = "cancelled"
)

func (s TaskStatus) Terminal() bool {
	return s == TaskPartialSuccess || s == TaskFailed || s == TaskCompleted || s == TaskCancelled
}

type WorkerStatus string

const (
	WorkerPending   WorkerStatus = "pending"
	WorkerRunning   WorkerStatus = "running"
	WorkerPausing   WorkerStatus = "pausing"
	WorkerPaused    WorkerStatus = "paused"
	WorkerCompleted WorkerStatus = "completed"
	WorkerFailed    WorkerStatus = "failed"
	WorkerCancelled WorkerStatus = "cancelled"
)

type Severity string

const (
	SeverityUrgent     Severity = "urgent"
	SeverityNormal     Severity = "normal"
	SeveritySuggestion Severity = "suggestion"
)

type ConnectionConfig struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateTaskRequest struct {
	ConnectionURL     string            `json:"connectionUrl"`
	Connection        *ConnectionConfig `json:"connection"`
	ExcludeTables     []string          `json:"excludeTables"`
	RedactObjectNames bool              `json:"redactObjectNames"`
}

type ResumeTaskRequest struct {
	ConnectionURL string            `json:"connectionUrl"`
	Connection    *ConnectionConfig `json:"connection"`
}

type Target struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

type Task struct {
	ID                 string       `json:"id"`
	Status             TaskStatus   `json:"status"`
	Stage              string       `json:"stage"`
	Target             Target       `json:"target"`
	CreatedAt          time.Time    `json:"createdAt"`
	UpdatedAt          time.Time    `json:"updatedAt"`
	ExpiresAt          time.Time    `json:"expiresAt"`
	PauseRequestedAt   *time.Time   `json:"pauseRequestedAt,omitempty"`
	PausedAt           *time.Time   `json:"pausedAt,omitempty"`
	ConnectionRequired bool         `json:"connectionRequired"`
	FailureSummary     string       `json:"failureSummary,omitempty"`
	Progress           TaskProgress `json:"progress"`
	AccessTokenHash    string       `json:"-"`
	RedactObjectNames  bool         `json:"-"`
}

type TaskProgress struct {
	Stage          string           `json:"stage"`
	TotalUnits     int              `json:"totalUnits"`
	CompletedUnits int              `json:"completedUnits"`
	Percent        int              `json:"percent"`
	CurrentObject  string           `json:"currentObject,omitempty"`
	TokenUsage     TokenUsage       `json:"tokenUsage"`
	WorkerFlow     []WorkerFlowNode `json:"workerFlow"`
	Workers        []WorkerProgress `json:"workers"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

// TokenUsage is the actual usage returned by the model provider through Eino.
// A provider that omits usage returns zeroes; values are never estimated.
type TokenUsage struct {
	InputTokens  int `json:"inputTokens"`
	OutputTokens int `json:"outputTokens"`
	TotalTokens  int `json:"totalTokens"`
}

func (usage TokenUsage) Add(other TokenUsage) TokenUsage {
	return TokenUsage{InputTokens: usage.InputTokens + other.InputTokens, OutputTokens: usage.OutputTokens + other.OutputTokens, TotalTokens: usage.TotalTokens + other.TotalTokens}
}

// WorkerFlowNode is the execution-plan contract returned with every task. The
// live status and counters remain in WorkerProgress; this metadata tells a
// client how to draw ordering, dependencies and parallel branches.
type WorkerFlowNode struct {
	WorkerID      string   `json:"workerId"`
	Name          string   `json:"name"`
	Stage         string   `json:"stage"`
	Order         int      `json:"order"`
	ParallelGroup string   `json:"parallelGroup,omitempty"`
	DependsOn     []string `json:"dependsOn"`
	AIInvolved    bool     `json:"aiInvolved"`
}

type WorkerProgress struct {
	WorkerID       string       `json:"workerId"`
	Name           string       `json:"name"`
	Status         WorkerStatus `json:"status"`
	TotalUnits     int          `json:"totalUnits"`
	CompletedUnits int          `json:"completedUnits"`
	CurrentObject  string       `json:"currentObject,omitempty"`
	FindingCount   int          `json:"findingCount"`
	TokenUsage     TokenUsage   `json:"tokenUsage"`
	Checkpoint     string       `json:"checkpoint,omitempty"`
	ErrorSummary   string       `json:"errorSummary,omitempty"`
	StartedAt      *time.Time   `json:"startedAt,omitempty"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	FinishedAt     *time.Time   `json:"finishedAt,omitempty"`
}

func DefaultWorkerFlow() []WorkerFlowNode {
	return []WorkerFlowNode{
		{WorkerID: "schema-collector", Name: "Schema 采集", Stage: "collecting_schema", Order: 1, DependsOn: []string{}},
		{WorkerID: "naming-convention", Name: "命名规范检查", Stage: "analyzing", Order: 2, ParallelGroup: "rule-analysis", DependsOn: []string{"schema-collector"}},
		{WorkerID: "index-performance", Name: "索引与性能检查(AI参与)", Stage: "analyzing", Order: 2, ParallelGroup: "rule-analysis", DependsOn: []string{"schema-collector"}, AIInvolved: true},
		{WorkerID: "type-constraint", Name: "类型与约束检查", Stage: "analyzing", Order: 2, ParallelGroup: "rule-analysis", DependsOn: []string{"schema-collector"}},
		{WorkerID: "semantic-governance", Name: "语义与治理检查(AI参与)", Stage: "analyzing", Order: 2, ParallelGroup: "rule-analysis", DependsOn: []string{"schema-collector"}, AIInvolved: true},
		{WorkerID: "finding-normalizer", Name: "问题聚合", Stage: "normalizing_findings", Order: 3, DependsOn: []string{"naming-convention", "index-performance", "type-constraint", "semantic-governance"}},
		{WorkerID: "ddl-generator", Name: "DDL 建议生成(AI参与)", Stage: "generating_ddl", Order: 4, DependsOn: []string{"finding-normalizer"}, AIInvolved: true},
		{WorkerID: "report-generator", Name: "报告生成(AI参与)", Stage: "generating_report", Order: 5, DependsOn: []string{"ddl-generator"}, AIInvolved: true},
	}
}

func DefaultWorkerProgress() []WorkerProgress {
	flow := DefaultWorkerFlow()
	workers := make([]WorkerProgress, 0, len(flow))
	for _, node := range flow {
		workers = append(workers, WorkerProgress{WorkerID: node.WorkerID, Name: node.Name, Status: WorkerPending})
	}
	return workers
}

type Column struct {
	Name          string `json:"name"`
	DataType      string `json:"dataType"`
	ColumnType    string `json:"columnType"`
	Nullable      bool   `json:"nullable"`
	DefaultValue  string `json:"defaultValue,omitempty"`
	DefaultIsNull bool   `json:"defaultIsNull"`
	Collation     string `json:"collation,omitempty"`
	Comment       string `json:"comment"`
}

type Index struct {
	Name      string   `json:"name"`
	Columns   []string `json:"columns"`
	NonUnique bool     `json:"nonUnique"`
	Primary   bool     `json:"primary"`
}

type Table struct {
	Name          string       `json:"name"`
	Comment       string       `json:"comment"`
	EstimatedRows int64        `json:"estimatedRows"`
	Collation     string       `json:"collation,omitempty"`
	Columns       []Column     `json:"columns"`
	Indexes       []Index      `json:"indexes"`
	ForeignKeys   []ForeignKey `json:"foreignKeys"`
}

type ForeignKey struct {
	Column           string `json:"column"`
	ReferencedTable  string `json:"referencedTable"`
	ReferencedColumn string `json:"referencedColumn"`
}

type SchemaSnapshot struct {
	Dialect          string    `json:"dialect"`
	Database         string    `json:"database"`
	DefaultCollation string    `json:"defaultCollation,omitempty"`
	Tables           []Table   `json:"tables"`
	UncheckedObjects []string  `json:"uncheckedObjects,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
}

type Finding struct {
	ID                   string    `json:"id"`
	RuleID               string    `json:"ruleId"`
	Severity             Severity  `json:"severity"`
	SeverityReason       string    `json:"severityReason"`
	ObjectType           string    `json:"objectType"`
	ObjectName           string    `json:"objectName"`
	Title                string    `json:"title"`
	Evidence             string    `json:"evidence"`
	Recommendation       string    `json:"recommendation"`
	DMLSuggestion        string    `json:"dmlSuggestion,omitempty"`
	DDLSuggestion        string    `json:"ddlSuggestion,omitempty"`
	PrecheckSQL          string    `json:"precheckSql,omitempty"`
	RollbackPlan         string    `json:"rollbackPlan,omitempty"`
	RequiresManualReview bool      `json:"requiresManualReview"`
	SQLRisk              string    `json:"sqlRisk"`
	CreatedAt            time.Time `json:"createdAt"`
}

type Report struct {
	TaskID          string           `json:"taskId"`
	GeneratedAt     time.Time        `json:"generatedAt"`
	RuleVersion     string           `json:"ruleVersion"`
	Summary         map[Severity]int `json:"summary"`
	Workers         []WorkerProgress `json:"workers"`
	Findings        []Finding        `json:"findings"`
	UncheckedObject []string         `json:"uncheckedObjects"`
	Failures        []string         `json:"failures"`
	Markdown        string           `json:"-"`
	HTML            string           `json:"-"`
}

type ReportRebuildState string

const (
	ReportRebuildIdle      ReportRebuildState = "idle"
	ReportRebuildPending   ReportRebuildState = "pending"
	ReportRebuildRunning   ReportRebuildState = "running"
	ReportRebuildSucceeded ReportRebuildState = "succeeded"
	ReportRebuildFailed    ReportRebuildState = "failed"
)

// ReportRebuildStatus is persisted separately from the original inspection
// task so report rebuilding can continue after its HTTP request has returned.
type ReportRebuildStatus struct {
	TaskID         string             `json:"taskId"`
	State          ReportRebuildState `json:"state"`
	RequestedAt    *time.Time         `json:"requestedAt,omitempty"`
	StartedAt      *time.Time         `json:"startedAt,omitempty"`
	FinishedAt     *time.Time         `json:"finishedAt,omitempty"`
	FailureSummary string             `json:"failureSummary,omitempty"`
}

type ReportPage struct {
	Report
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type BaselineComparison struct {
	BaselineTaskID string   `json:"baselineTaskId,omitempty"`
	NewFindingIDs  []string `json:"newFindingIds"`
}

type TaskCreateResponse struct {
	TaskID          string    `json:"taskId"`
	TaskAccessToken string    `json:"taskAccessToken"`
	ExpiresAt       time.Time `json:"expiresAt"`
}

type AuditEvent struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"taskId,omitempty"`
	Action    string    `json:"action"`
	Outcome   string    `json:"outcome"`
	Detail    string    `json:"detail,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
