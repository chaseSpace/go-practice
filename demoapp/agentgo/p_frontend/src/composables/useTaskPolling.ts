import { onBeforeUnmount, ref } from 'vue'
import { AxiosError } from 'axios'
import { getTask } from '@/api/tasks'
import { clearTaskToken } from '@/api/client'
import type { Task } from '@/types/domain'

export function useTaskPolling(taskId: string, onUnavailable?: () => void) {
  const task = ref<Task>()
  const loading = ref(false)
  const error = ref('')
  let timer: number | undefined
  let stopped = false

  const stop = () => {
    stopped = true
    if (timer) window.clearTimeout(timer)
  }

  const refresh = async () => {
    if (stopped) return
    loading.value = !task.value
    try {
      task.value = await getTask(taskId)
      error.value = ''
      if (!task.value.status.match(/^(completed|partial_success|failed|cancelled)$/)) {
        timer = window.setTimeout(refresh, task.value.status === 'paused' ? 3_000 : 1_000)
      }
    } catch (cause) {
      const status = cause instanceof AxiosError ? cause.response?.status : undefined
      if (status === 404) {
        clearTaskToken(taskId)
        onUnavailable?.()
      } else {
        error.value = '无法获取任务进度，请稍后重试。'
        timer = window.setTimeout(refresh, 3_000)
      }
    } finally {
      loading.value = false
    }
  }

  onBeforeUnmount(stop)
  return { task, loading, error, refresh, stop }
}
