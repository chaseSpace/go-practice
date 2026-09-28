import { api } from './client'
import type { AuditEvent, CreateTaskPayload, LLMExplanation, ReportRebuildStatus, ReportPage, Task, TaskCreateResponse } from '@/types/domain'

export async function createTask(payload: CreateTaskPayload) {
  return (await api.post<TaskCreateResponse>('/tasks', payload)).data
}

export async function getTask(taskId: string) {
  return (await api.get<Task>(`/tasks/${taskId}`)).data
}

export async function pauseTask(taskId: string) {
  return (await api.post<Task>(`/tasks/${taskId}/pause`)).data
}

export async function resumeTask(taskId: string, payload: Partial<CreateTaskPayload> = {}) {
  return (await api.post<Task>(`/tasks/${taskId}/resume`, payload)).data
}

export async function cancelTask(taskId: string) {
  return (await api.post<Task>(`/tasks/${taskId}/cancel`)).data
}

export async function getAudit(taskId: string) {
  return (await api.get<AuditEvent[]>(`/tasks/${taskId}/audit`)).data
}

export async function explainFinding(taskId: string, findingId: string) {
  return (await api.post<LLMExplanation>(`/reports/${taskId}/findings/${findingId}/explain`)).data
}

export async function getReport(taskId: string, params: { severity?: string; keyword?: string; page?: number; page_size?: number } = {}) {
	return (await api.get<ReportPage>(`/reports/${taskId}`, { params })).data
}

export async function rebuildReport(taskId: string) {
  return (await api.post<ReportRebuildStatus>(`/reports/${taskId}/rebuild`)).data
}

export async function reanalyzeReport(taskId: string) {
  return (await api.post<ReportRebuildStatus>(`/reports/${taskId}/reanalyze`)).data
}

export async function getReportRebuildStatus(taskId: string) {
  return (await api.get<ReportRebuildStatus>(`/reports/${taskId}/rebuild`)).data
}

export async function getPreview(taskId: string, format: 'html' | 'markdown') {
  return (await api.get<string>(`/reports/${taskId}/preview`, { params: { format }, responseType: 'text' })).data
}

export function download(taskId: string, resource: 'sql' | 'export', format?: 'html' | 'markdown' | 'json') {
  const token = sessionStorage.getItem(`agentgo.task-token.${taskId}`)
  const query = format ? `?format=${format}` : ''
  return api.get(`/reports/${taskId}/${resource}${query}`, {
    responseType: 'blob',
    headers: token ? { 'X-Task-Access-Token': token } : {},
  })
}
