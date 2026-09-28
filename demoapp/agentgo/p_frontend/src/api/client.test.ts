import type { InternalAxiosRequestConfig } from 'axios'
import { afterEach, expect, it } from 'vitest'
import { api, clearTaskToken, saveTaskToken } from './client'

afterEach(() => sessionStorage.clear())

it('injects a task access token only for task-scoped requests', async () => {
  let captured: InternalAxiosRequestConfig | undefined
  api.defaults.adapter = async (config) => {
    captured = config
    return { data: {}, status: 200, statusText: 'OK', headers: {}, config }
  }
  saveTaskToken('task-1', 'short-lived-token')
  await api.get('/tasks/task-1')
  expect(captured?.headers.get('X-Task-Access-Token')).toBe('short-lived-token')
  clearTaskToken('task-1')
  await api.get('/tasks/task-1')
  expect(captured?.headers.get('X-Task-Access-Token')).toBeUndefined()
})
