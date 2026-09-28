import axios from 'axios'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080/api/v1',
  timeout: 15_000,
  headers: { 'Content-Type': 'application/json' },
})

export const tokenKey = (taskId: string) => `agentgo.task-token.${taskId}`

export function saveTaskToken(taskId: string, token: string) {
  sessionStorage.setItem(tokenKey(taskId), token)
}

export function clearTaskToken(taskId: string) {
  sessionStorage.removeItem(tokenKey(taskId))
}

api.interceptors.request.use((config) => {
  const match = config.url?.match(/\/(?:tasks|reports)\/([^/?]+)/)
  if (match) {
    const token = sessionStorage.getItem(tokenKey(match[1]))
    if (token) config.headers.set('X-Task-Access-Token', token)
  }
  return config
})
