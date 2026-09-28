import { defineStore } from 'pinia'

export interface RecentTask {
  id: string
  database: string
  dialect?: 'mysql' | 'postgres'
  maskedUrl?: string
  createdAt: string
  status: string
}

const storageKey = 'agentgo.recent-tasks'

type TaskTarget = { database: string; dialect: string }

function enrichTask(task: RecentTask): RecentTask {
  if (!task.maskedUrl) return task
  try {
    const url = new URL(task.maskedUrl)
    const database = decodeURIComponent(url.pathname.replace(/^\//, ''))
    const dialect = url.protocol === 'postgres:' || url.protocol === 'postgresql:' ? 'postgres' : 'mysql'
    return {
      ...task,
      database: !task.database || task.database === '数据库' || task.database === 'MySQL 数据库' ? database || '—' : task.database,
      dialect: task.dialect ?? dialect,
    }
  } catch {
    return task
  }
}

function readHistory(): RecentTask[] {
  try {
    const value = JSON.parse(localStorage.getItem(storageKey) ?? '[]')
    return Array.isArray(value) ? value.map(enrichTask) : []
  } catch {
    return []
  }
}

export const useHistoryStore = defineStore('history', {
  state: () => ({ tasks: readHistory() as RecentTask[] }),
  actions: {
    remember(task: RecentTask) {
      this.tasks = [task, ...this.tasks.filter((item) => item.id !== task.id)].slice(0, 12)
      localStorage.setItem(storageKey, JSON.stringify(this.tasks))
    },
    update(id: string, status: string, target?: TaskTarget) {
      const task = this.tasks.find((item) => item.id === id)
      if (task) {
        task.status = status
        if (target) {
          task.database = target.database || task.database
          task.dialect = target.dialect === 'postgres' || target.dialect === 'postgresql' ? 'postgres' : 'mysql'
        }
        localStorage.setItem(storageKey, JSON.stringify(this.tasks))
      }
    },
    remove(id: string) {
      this.tasks = this.tasks.filter((item) => item.id !== id)
      localStorage.setItem(storageKey, JSON.stringify(this.tasks))
    },
  },
})
