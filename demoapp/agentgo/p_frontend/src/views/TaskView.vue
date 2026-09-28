<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { cancelTask, getAudit, pauseTask, resumeTask } from '@/api/tasks'
import { useTaskPolling } from '@/composables/useTaskPolling'
import { useHistoryStore } from '@/stores/history'
import type { AuditEvent, ConnectionConfig, WorkerFlowNode, WorkerProgress } from '@/types/domain'

const props = defineProps<{ taskId: string }>()
const router = useRouter()
const history = useHistoryStore()
const resumeURL = ref('')
const resumeMode = ref<'url' | 'fields'>('url')
const resumeConnection = ref<ConnectionConfig>({ dialect: 'mysql', host: '', port: 3306, database: '', username: '', password: '' })
const showResumeDialog = ref(false)
const actionLoading = ref(false)
const auditEvents = ref<AuditEvent[]>([])
const { task, loading, error, refresh } = useTaskPolling(props.taskId, () => router.replace('/history'))

const statusLabel = computed(() => ({
  pending: '等待执行', running: '进行中', pausing: '正在暂停', paused: '已暂停', partial_success: '部分完成', failed: '执行失败', completed: '已完成', cancelled: '已取消',
}[task.value?.status ?? 'pending']))

const statusType = computed(() => ({
  pending: 'info', running: 'primary', pausing: 'warning', paused: 'warning', partial_success: 'warning', failed: 'danger', completed: 'success', cancelled: 'info',
}[task.value?.status ?? 'pending']) as 'primary' | 'success' | 'warning' | 'danger' | 'info')

type FlowStage = { order: number; nodes: WorkerFlowNode[] }
const workerByID = computed(() => new Map((task.value?.progress.workers ?? []).map((worker) => [worker.workerId, worker])))
const flowStages = computed<FlowStage[]>(() => {
  const stages = new Map<number, WorkerFlowNode[]>()
  for (const node of task.value?.progress.workerFlow ?? []) {
    stages.set(node.order, [...(stages.get(node.order) ?? []), node])
  }
  return [...stages.entries()].sort(([left], [right]) => left - right).map(([order, nodes]) => ({ order, nodes }))
})

function workerStatus(status: string) {
  return ({ pending: '等待中', running: '进行中', pausing: '正在暂停', paused: '已暂停', completed: '已完成', failed: '失败', cancelled: '已取消' }[status] ?? status)
}

function stageLabel(stage: string) {
  return ({ collecting_schema: '采集 Schema', analyzing: '并行规则分析', normalizing_findings: '聚合问题', generating_ddl: '生成 DDL 建议', generating_report: '生成报告' }[stage] ?? stage)
}

function flowWorker(node: WorkerFlowNode): WorkerProgress | undefined { return workerByID.value.get(node.workerId) }
function flowStatus(node: WorkerFlowNode) { return flowWorker(node)?.status ?? 'pending' }
function flowPercentage(node: WorkerFlowNode) {
  const worker = flowWorker(node)
  if (!worker) return 0
  if (worker.totalUnits) return Math.floor(worker.completedUnits * 100 / worker.totalUnits)
  return worker.status === 'completed' ? 100 : 0
}
function flowDependency(node: WorkerFlowNode) {
  const names = node.dependsOn.map((id) => task.value?.progress.workerFlow.find((candidate) => candidate.workerId === id)?.name).filter(Boolean)
  return names.length ? `依赖：${names.join('、')}` : '起始节点'
}
function flowTagType(node: WorkerFlowNode): 'primary' | 'success' | 'warning' | 'danger' | 'info' {
  const status = flowStatus(node)
  return status === 'failed' ? 'danger' : status === 'completed' ? 'success' : status === 'running' ? 'primary' : status === 'paused' || status === 'pausing' ? 'warning' : 'info'
}

async function pause() {
  actionLoading.value = true
  try { await pauseTask(props.taskId); await refresh() } catch { ElMessage.error('暂停请求失败，请稍后重试。') } finally { actionLoading.value = false }
}

async function resume() {
  if (task.value?.connectionRequired || task.value?.status === 'failed') { showResumeDialog.value = true; return }
  await performResume()
}

async function performResume() {
  actionLoading.value = true
  try {
    await resumeTask(props.taskId, resumeMode.value === 'url' ? { connectionUrl: resumeURL.value } : { connection: resumeConnection.value })
    resumeURL.value = ''
    showResumeDialog.value = false
    await refresh()
  } catch {
    ElMessage.error('无法继续任务，请确认连接信息。')
  } finally { actionLoading.value = false }
}

async function cancel() {
  try {
    await ElMessageBox.confirm('取消后无法继续本次任务，确认取消吗？', '取消检查', { type: 'warning' })
    actionLoading.value = true
    await cancelTask(props.taskId)
    await refresh()
  } catch {
    // User cancellation and request errors are intentionally not surfaced twice.
  } finally { actionLoading.value = false }
}

function openReport() { router.push(`/reports/${props.taskId}`) }

async function loadAudit() { try { auditEvents.value = await getAudit(props.taskId) } catch { auditEvents.value = [] } }

watch(task, (value) => { if (value) history.update(value.id, value.status, value.target) })
onMounted(async () => { await refresh(); await loadAudit() })
</script>

<template>
  <section v-loading="loading" class="page-card task-page">
    <template v-if="task">
      <div class="task-heading">
        <div>
          <el-button text @click="router.push('/')">← 返回首页</el-button>
          <h1>{{ task.target.database }}</h1>
          <p>{{ task.target.dialect }} · {{ task.target.host }}:{{ task.target.port }}</p>
        </div>
        <div class="task-actions">
          <el-tag :type="statusType" effect="dark">{{ statusLabel }}</el-tag>
          <el-button v-if="task.status === 'running'" type="warning" :loading="actionLoading" @click="pause">暂停</el-button>
          <el-button v-if="task.status === 'paused' || task.status === 'failed'" type="primary" :loading="actionLoading" @click="resume">{{ task.status === 'failed' ? '恢复任务' : '继续' }}</el-button>
          <el-button v-if="['running', 'pausing', 'paused'].includes(task.status)" :loading="actionLoading" @click="cancel">取消检查</el-button>
          <el-button v-if="['completed', 'partial_success'].includes(task.status)" type="primary" @click="openReport">查看报告</el-button>
        </div>
      </div>

      <el-alert v-if="task.status === 'pausing'" title="暂停请求已保存，正在等待各 Worker 完成当前表或阶段。" type="warning" :closable="false" />
      <el-alert v-if="task.status === 'paused'" :title="task.connectionRequired ? '任务已暂停，服务重启后需重新提交连接信息才能继续。' : '任务已暂停，数据库连接已释放。'" type="info" :closable="false" />
      <el-alert v-if="task.failureSummary" :title="task.failureSummary" type="error" :closable="false" />
      <el-alert v-if="error" :title="error" type="warning" :closable="false" />

      <div class="overall-progress"><span>已处理 {{ task.progress.completedUnits }} / {{ task.progress.totalUnits }} 个单元</span><el-progress :percentage="task.progress.percent" /></div>

      <section class="worker-flow-panel">
        <div class="flow-heading"><div><h2>Worker 执行流程</h2><p>实线箭头表示依赖顺序；并列节点会并行检查各自负责的规则集。</p></div><el-tag type="info">当前阶段：{{ task.progress.stage }}</el-tag></div>
        <div class="worker-flow" aria-label="Worker 执行流程图">
          <template v-for="(stage, stageIndex) in flowStages" :key="stage.order">
            <div v-if="stageIndex" class="flow-connector" aria-hidden="true" />
            <div class="flow-stage" :class="{ 'flow-stage--parallel': stage.nodes.length > 1 }">
              <article v-for="node in stage.nodes" :key="node.workerId" class="flow-node" :class="`flow-node--${flowStatus(node)}`">
                <div class="flow-node-title"><strong>{{ node.name }}</strong><el-tag v-if="node.aiInvolved" size="small" type="primary" effect="plain">AI</el-tag></div>
                <el-tag size="small" :type="flowTagType(node)">{{ workerStatus(flowStatus(node)) }}</el-tag>
                <el-progress :percentage="flowPercentage(node)" :stroke-width="5" :show-text="false" />
                <p>{{ stageLabel(node.stage) }} · {{ flowWorker(node)?.completedUnits ?? 0 }}/{{ flowWorker(node)?.totalUnits ?? 0 }} · {{ flowWorker(node)?.findingCount ?? 0 }} 个问题</p>
                <small>Token：{{ flowWorker(node)?.tokenUsage?.totalTokens ?? 0 }}（输入 {{ flowWorker(node)?.tokenUsage?.inputTokens ?? 0 }} / 输出 {{ flowWorker(node)?.tokenUsage?.outputTokens ?? 0 }}）</small>
                <small>{{ flowWorker(node)?.currentObject || flowDependency(node) }}</small>
                <small v-if="flowWorker(node)?.errorSummary" class="flow-node-error">{{ flowWorker(node)?.errorSummary }}</small>
              </article>
            </div>
          </template>
        </div>
      </section>

      <div class="overall-progress"><span>模型总消耗 {{ task.progress.tokenUsage.totalTokens }} Token（输入 {{ task.progress.tokenUsage.inputTokens }} / 输出 {{ task.progress.tokenUsage.outputTokens }}）</span></div>

      <h2>Worker 进度明细</h2>
      <el-table :data="task.progress.workers" size="small" class="worker-table">
        <el-table-column prop="name" label="Worker" min-width="150" />
        <el-table-column label="状态" width="110"><template #default="{ row }"><el-tag size="small" :type="row.status === 'failed' ? 'danger' : row.status === 'completed' ? 'success' : row.status === 'running' ? 'primary' : 'info'">{{ workerStatus(row.status) }}</el-tag></template></el-table-column>
        <el-table-column label="进度" min-width="180"><template #default="{ row }"><el-progress :percentage="row.totalUnits ? Math.floor(row.completedUnits * 100 / row.totalUnits) : 0" :format="() => `${row.completedUnits}/${row.totalUnits}`" /></template></el-table-column>
        <el-table-column prop="currentObject" label="当前对象" min-width="150" />
        <el-table-column prop="findingCount" label="问题数" width="90" />
        <el-table-column label="Token" width="140"><template #default="{ row }">{{ row.tokenUsage?.totalTokens ?? 0 }}<small class="token-detail">{{ row.tokenUsage ? `（${row.tokenUsage.inputTokens}/${row.tokenUsage.outputTokens}）` : '' }}</small></template></el-table-column>
        <el-table-column prop="errorSummary" label="失败摘要" min-width="200" />
      </el-table>
      <el-collapse class="audit-panel"><el-collapse-item title="任务审计事件" name="audit"><el-timeline><el-timeline-item v-for="event in auditEvents" :key="event.id" :timestamp="new Date(event.createdAt).toLocaleString()" :type="event.outcome === 'success' ? 'success' : 'warning'">{{ event.action }}{{ event.detail ? `：${event.detail}` : '' }}</el-timeline-item></el-timeline></el-collapse-item></el-collapse>
    </template>
  </section>

  <el-dialog v-model="showResumeDialog" title="重新提交连接信息" width="520px">
    <p class="dialog-help">服务重启后不会保存密码。请重新提交同一目标数据库的连接信息以继续未完成的分片。</p>
    <el-tabs v-model="resumeMode">
      <el-tab-pane label="数据库 URL" name="url"><el-input v-model="resumeURL" autocomplete="off" placeholder="mysql://用户名:密码@主机:3306/数据库" /></el-tab-pane>
      <el-tab-pane label="逐项填写" name="fields"><div class="form-grid compact-form"><el-form-item label="数据库类型"><el-select v-model="resumeConnection.dialect" @change="resumeConnection.port = resumeConnection.dialect === 'postgres' ? 5432 : 3306"><el-option label="MySQL" value="mysql" /><el-option label="PostgreSQL" value="postgres" /></el-select></el-form-item><el-form-item label="主机"><el-input v-model="resumeConnection.host" /></el-form-item><el-form-item label="端口"><el-input-number v-model="resumeConnection.port" :min="1" :max="65535" /></el-form-item><el-form-item label="数据库"><el-input v-model="resumeConnection.database" /></el-form-item><el-form-item label="用户名"><el-input v-model="resumeConnection.username" /></el-form-item><el-form-item label="密码"><el-input v-model="resumeConnection.password" type="password" show-password /></el-form-item></div></el-tab-pane>
    </el-tabs>
    <template #footer><el-button @click="showResumeDialog = false">取消</el-button><el-button type="primary" :loading="actionLoading" @click="performResume">继续任务</el-button></template>
  </el-dialog>
</template>
