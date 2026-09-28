<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { QuestionFilled } from '@element-plus/icons-vue'
import { getReportRebuildStatus } from '@/api/tasks'
import { useHistoryStore } from '@/stores/history'
import type { ReportRebuildState, ReportRebuildStatus } from '@/types/domain'

const router = useRouter()
const history = useHistoryStore()
const rebuildStatusByTask = ref<Record<string, ReportRebuildStatus>>({})
let refreshTimer: number | undefined

const rebuildStateLabel: Record<ReportRebuildState, string> = { idle: '未生成', pending: '排队中', running: '生成中', succeeded: '已完成', failed: '生成失败' }
const rebuildStateType: Record<ReportRebuildState, 'info' | 'warning' | 'success' | 'danger'> = { idle: 'info', pending: 'warning', running: 'warning', succeeded: 'success', failed: 'danger' }

function dialectLabel(value?: string) {
  if (value === 'postgres') return 'PostgreSQL'
  if (value === 'mysql') return 'MySQL'
  return '—'
}

function reportStatus(taskId: string): ReportRebuildStatus {
  return rebuildStatusByTask.value[taskId] ?? { taskId, state: 'idle' }
}

function formatTime(value?: string) {
  return value ? new Date(value).toLocaleString() : '—'
}

function stopStatusRefresh() {
  if (refreshTimer !== undefined) window.clearTimeout(refreshTimer)
  refreshTimer = undefined
}

async function refreshReportStatuses() {
  const statuses = await Promise.all(history.tasks.map(async (task) => {
    try { return await getReportRebuildStatus(task.id) } catch { return undefined }
  }))
  const next: Record<string, ReportRebuildStatus> = {}
  for (const status of statuses) {
    if (status) next[status.taskId] = status
  }
  rebuildStatusByTask.value = next
  stopStatusRefresh()
  if (statuses.some((status) => status?.state === 'pending' || status?.state === 'running')) {
    refreshTimer = window.setTimeout(() => { void refreshReportStatuses() }, 1500)
  }
}

onMounted(() => { void refreshReportStatuses() })
onBeforeUnmount(stopStatusRefresh)
</script>

<template>
  <section class="page-card history-page">
    <div class="report-toolbar"><div><h1>最近任务</h1><p>仅保存任务摘要；访问令牌只保留在当前浏览器会话中。</p></div><el-button type="primary" @click="router.push('/')">新建检查</el-button></div>
    <el-empty v-if="history.tasks.length === 0" description="暂无任务记录" />
    <el-table v-else :data="history.tasks">
      <el-table-column label="DB类型" width="130"><template #default="{ row }">{{ dialectLabel(row.dialect) }}</template></el-table-column>
      <el-table-column label="数据库" min-width="240">
        <template #default="{ row }"><span class="history-db">{{ row.database }}</span><el-tooltip v-if="row.maskedUrl" :content="row.maskedUrl" placement="top"><el-icon class="history-url-help"><QuestionFilled /></el-icon></el-tooltip></template>
      </el-table-column>
      <el-table-column prop="status" label="最近状态" width="140" />
      <el-table-column label="报告生成" width="120"><template #default="{ row }"><el-tag :type="rebuildStateType[reportStatus(row.id).state]">{{ rebuildStateLabel[reportStatus(row.id).state] }}</el-tag></template></el-table-column>
      <el-table-column label="报告生成时间" min-width="270"><template #default="{ row }"><div v-if="reportStatus(row.id).state !== 'idle'" class="report-time-cell"><span v-if="reportStatus(row.id).requestedAt">提交：{{ formatTime(reportStatus(row.id).requestedAt) }}</span><span v-if="reportStatus(row.id).startedAt">开始：{{ formatTime(reportStatus(row.id).startedAt) }}</span><span v-if="reportStatus(row.id).finishedAt">完成：{{ formatTime(reportStatus(row.id).finishedAt) }}</span></div><span v-else>—</span></template></el-table-column>
      <el-table-column label="创建时间" min-width="180"><template #default="{ row }">{{ new Date(row.createdAt).toLocaleString() }}</template></el-table-column>
      <el-table-column label="操作" width="180" fixed="right"><template #default="{ row }"><el-button link type="primary" @click="router.push(`/tasks/${row.id}`)">查看任务</el-button><el-button link type="danger" @click="history.remove(row.id)">移除</el-button></template></el-table-column>
    </el-table>
  </section>
</template>
