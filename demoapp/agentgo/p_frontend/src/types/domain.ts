export type TaskStatus =
  | 'pending'
  | 'running'
  | 'pausing'
  | 'paused'
  | 'partial_success'
  | 'failed'
  | 'completed'
  | 'cancelled'

export type WorkerStatus = 'pending' | 'running' | 'pausing' | 'paused' | 'completed' | 'failed' | 'cancelled'
export type Severity = 'urgent' | 'normal' | 'suggestion'
export type DatabaseDialect = 'mysql' | 'postgres'

export interface TokenUsage {
  inputTokens: number
  outputTokens: number
  totalTokens: number
}

export interface ConnectionConfig {
  dialect: DatabaseDialect
  host: string
  port: number
  database: string
  username: string
  password: string
}

export interface CreateTaskPayload {
  connectionUrl?: string
  connection?: ConnectionConfig
  excludeTables?: string[]
  redactObjectNames?: boolean
}

export interface WorkerProgress {
  workerId: string
  name: string
  status: WorkerStatus
  totalUnits: number
  completedUnits: number
  currentObject?: string
  findingCount: number
  tokenUsage: TokenUsage
  checkpoint?: string
  errorSummary?: string
  startedAt?: string
  updatedAt: string
  finishedAt?: string
}

export interface WorkerFlowNode {
  workerId: string
  name: string
  stage: string
  order: number
  parallelGroup?: string
  dependsOn: string[]
  aiInvolved: boolean
}

export interface Task {
  id: string
  status: TaskStatus
  stage: string
  target: { dialect: string; host: string; port: number; database: string }
  createdAt: string
  updatedAt: string
  expiresAt: string
  connectionRequired: boolean
  failureSummary?: string
  progress: {
    stage: string
    totalUnits: number
    completedUnits: number
    percent: number
    currentObject?: string
    tokenUsage: TokenUsage
    workerFlow: WorkerFlowNode[]
    workers: WorkerProgress[]
    updatedAt: string
  }
}

export interface Finding {
  id: string
  ruleId: string
  severity: Severity
  severityReason: string
  objectType: string
  objectName: string
  title: string
  evidence: string
  recommendation: string
  dmlSuggestion?: string
  ddlSuggestion?: string
  precheckSql?: string
  rollbackPlan?: string
  requiresManualReview: boolean
  sqlRisk: string
  createdAt: string
}

export interface Report {
  taskId: string
  generatedAt: string
  ruleVersion: string
  summary: Record<Severity, number>
  workers: WorkerProgress[]
  findings: Finding[]
  uncheckedObjects: string[]
  failures: string[]
}

export interface ReportPage extends Report {
	total: number
	page: number
	pageSize: number
}

export type ReportRebuildState = 'idle' | 'pending' | 'running' | 'succeeded' | 'failed'

export interface ReportRebuildStatus {
	taskId: string
	state: ReportRebuildState
	requestedAt?: string
	startedAt?: string
	finishedAt?: string
	failureSummary?: string
}

export interface TaskCreateResponse {
  taskId: string
  taskAccessToken: string
  expiresAt: string
}

export interface AuditEvent {
  id: string
  taskId?: string
  action: string
  outcome: string
  detail?: string
  createdAt: string
}

export interface LLMExplanation {
  content: string
  model: string
  tokenUsage: TokenUsage
  generatedAt: string
}
