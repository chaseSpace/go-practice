<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { download, explainFinding, getPreview, getReport, getReportRebuildStatus, reanalyzeReport, rebuildReport } from '@/api/tasks'
import MarkdownArticle from '@/components/MarkdownArticle.vue'
import type { Finding, LLMExplanation, ReportPage, ReportRebuildStatus, Severity } from '@/types/domain'

const props = defineProps<{ taskId: string }>()
const router = useRouter()
const report = ref<ReportPage>()
const loading = ref(false)
const error = ref('')
const severity = ref<Severity | ''>('')
const keyword = ref('')
const page = ref(1)
const pageSize = ref(20)
const selected = ref<Finding>()
const drawer = ref(false)
const explanation = ref<LLMExplanation>()
const explaining = ref(false)
const markdown = ref('')
const markdownLoading = ref(false)
const rebuildStatus = ref<ReportRebuildStatus>()
let rebuildPollTimer: number | undefined
let rebuildRequestedByUser = false
let rebuildMode: 'regenerate' | 'reanalyze' = 'regenerate'

const severityLabel: Record<Severity, string> = { urgent: '紧急', normal: '一般', suggestion: '建议' }
const severityType: Record<Severity, 'danger' | 'warning' | 'info'> = { urgent: 'danger', normal: 'warning', suggestion: 'info' }
const summary = computed(() => report.value?.summary ?? { urgent: 0, normal: 0, suggestion: 0 })
const reportFindingCount = computed(() => (summary.value.urgent ?? 0) + (summary.value.normal ?? 0) + (summary.value.suggestion ?? 0))
const workerFindingCount = computed(() => (report.value?.workers ?? []).reduce((total, worker) => total + worker.findingCount, 0))
const needsReportRepair = computed(() => workerFindingCount.value > reportFindingCount.value)
const rebuilding = computed(() => rebuildStatus.value?.state === 'pending' || rebuildStatus.value?.state === 'running')
const rebuildMessage = computed(() => {
  if (rebuildStatus.value?.state === 'pending') return rebuildMode === 'reanalyze' ? '已提交 Schema 重新分析，正在等待后台执行。' : '报告重建已提交，正在等待后台执行。'
  return rebuildMode === 'reanalyze' ? '正在重新分析 Schema 并生成报告；完成后会自动刷新。' : '正在根据已保存 Finding 生成报告；完成后会自动刷新。'
})
const rebuildStateLabel: Record<ReportRebuildStatus['state'], string> = { idle: '未生成', pending: '排队中', running: '生成中', succeeded: '已完成', failed: '生成失败' }
const rebuildStateType: Record<ReportRebuildStatus['state'], 'info' | 'warning' | 'success' | 'danger'> = { idle: 'info', pending: 'warning', running: 'warning', succeeded: 'success', failed: 'danger' }

function formatTime(value?: string) {
  return value ? new Date(value).toLocaleString() : '—'
}

async function load() {
  loading.value = true
  try {
    report.value = await getReport(props.taskId, { severity: severity.value || undefined, keyword: keyword.value || undefined, page: page.value, page_size: pageSize.value })
    error.value = ''
  } catch {
    error.value = '报告尚未生成或任务已过期。'
  } finally { loading.value = false }
}

async function loadMarkdown() {
  markdownLoading.value = true
  try { markdown.value = await getPreview(props.taskId, 'markdown') } catch { markdown.value = '' } finally { markdownLoading.value = false }
}

function stopRebuildPolling() {
  if (rebuildPollTimer !== undefined) window.clearTimeout(rebuildPollTimer)
  rebuildPollTimer = undefined
}

function scheduleRebuildPolling() {
  stopRebuildPolling()
  rebuildPollTimer = window.setTimeout(() => { void loadRebuildStatus() }, 1500)
}

async function loadRebuildStatus() {
  try {
    const status = await getReportRebuildStatus(props.taskId)
    rebuildStatus.value = status
    if (status.state === 'pending' || status.state === 'running') {
      scheduleRebuildPolling()
      return
    }
    stopRebuildPolling()
    if (status.state === 'succeeded' && rebuildRequestedByUser) {
      rebuildRequestedByUser = false
      await Promise.all([load(), loadMarkdown()])
      ElMessage.success('报告已在后台重新生成')
    }
    if (status.state === 'failed' && rebuildRequestedByUser) {
      rebuildRequestedByUser = false
      ElMessage.error(status.failureSummary || '报告重建失败，请查看后端日志后重试。')
    }
  } catch {
    stopRebuildPolling()
  }
}

async function submitReportGeneration(mode: 'regenerate' | 'reanalyze') {
  try {
    rebuildMode = mode
    rebuildStatus.value = mode === 'reanalyze' ? await reanalyzeReport(props.taskId) : await rebuildReport(props.taskId)
    rebuildRequestedByUser = true
    ElMessage.success(mode === 'reanalyze' ? '已提交后台 Schema 重新分析' : '已提交后台报告生成')
    await loadRebuildStatus()
  } catch {
    ElMessage.error('提交报告生成失败，请确认后端服务可用后重试。')
  }
}

function regenerateReport() { return submitReportGeneration('regenerate') }
function reanalyzeAndGenerateReport() { return submitReportGeneration('reanalyze') }

function detail(finding: Finding) { selected.value = finding; explanation.value = undefined; drawer.value = true }
async function copy(value?: string) {
  if (!value) return
  await navigator.clipboard.writeText(value)
  ElMessage.success('已复制')
}
async function explain() {
  if (!selected.value) return
  explaining.value = true
  try { explanation.value = await explainFinding(props.taskId, selected.value.id) }
  catch { ElMessage.warning('LLM 辅助解释未配置或暂时不可用。') }
  finally { explaining.value = false }
}

async function save(resource: 'sql' | 'export', format?: 'html' | 'markdown' | 'json') {
  try {
    const response = await download(props.taskId, resource, format)
    const href = URL.createObjectURL(response.data)
    const link = document.createElement('a')
    link.href = href
    link.download = resource === 'sql' ? 'database-suggestions.sql' : `database-report.${format ?? 'html'}`
    link.click()
    URL.revokeObjectURL(href)
  } catch { ElMessage.error('下载失败') }
}

watch([severity, keyword], () => { page.value = 1; load() })
watch(() => props.taskId, () => { stopRebuildPolling(); rebuildStatus.value = undefined; rebuildRequestedByUser = false; rebuildMode = 'regenerate'; load(); loadMarkdown(); loadRebuildStatus() })
onMounted(() => { load(); loadMarkdown(); loadRebuildStatus() })
onBeforeUnmount(stopRebuildPolling)
</script>

<template>
  <section v-loading="loading" class="report-page">
    <div class="report-toolbar">
      <div><el-button text @click="router.push(`/tasks/${taskId}`)">← 返回任务</el-button><h1>数据库体检报告</h1><p v-if="report">规则版本：{{ report.ruleVersion }} · {{ new Date(report.generatedAt).toLocaleString() }}</p></div>
      <div class="report-actions"><el-button :loading="rebuilding" :disabled="rebuilding" @click="regenerateReport">重新生成报告</el-button><el-popconfirm title="将重新运行 Schema 快照上的所有 AI 分析，耗时和 Token 消耗都会增加。" confirm-button-text="重新分析" cancel-button-text="取消" @confirm="reanalyzeAndGenerateReport"><template #reference><el-button type="warning" plain :disabled="rebuilding">重新分析并生成</el-button></template></el-popconfirm><el-button @click="router.push(`/reports/${taskId}/preview`)">在线预览</el-button><el-dropdown @command="(format) => save('export', format)"><el-button>下载报告</el-button><template #dropdown><el-dropdown-menu><el-dropdown-item command="html">HTML</el-dropdown-item><el-dropdown-item command="markdown">Markdown</el-dropdown-item><el-dropdown-item command="json">JSON</el-dropdown-item></el-dropdown-menu></template></el-dropdown><el-button type="primary" @click="save('sql')">下载 SQL 建议</el-button></div>
    </div>
    <el-alert v-if="error" :title="error" type="warning" :closable="false" />
    <template v-if="report">
      <el-alert v-if="rebuilding" :title="rebuildMessage" type="info" :closable="false" show-icon />
      <el-alert v-else-if="rebuildStatus?.state === 'failed'" :title="rebuildStatus.failureSummary || '报告重建失败，请查看后端日志后重试。'" type="error" :closable="false" show-icon />
      <section v-if="rebuildStatus && rebuildStatus.state !== 'idle'" class="page-card report-rebuild-status">
        <div><h2>报告生成状态</h2><p>状态会在后台生成期间自动刷新。</p></div>
        <el-tag :type="rebuildStateType[rebuildStatus.state]">{{ rebuildStateLabel[rebuildStatus.state] }}</el-tag>
        <div class="report-rebuild-times"><span v-if="rebuildStatus.requestedAt">提交：{{ formatTime(rebuildStatus.requestedAt) }}</span><span v-if="rebuildStatus.startedAt">开始：{{ formatTime(rebuildStatus.startedAt) }}</span><span v-if="rebuildStatus.finishedAt">完成：{{ formatTime(rebuildStatus.finishedAt) }}</span></div>
      </section>
      <el-alert v-if="needsReportRepair" title="检测到 Worker 已发现的问题数与报告快照不一致，需要重新分析已保存的 Schema 快照以修复。" type="warning" :closable="false"><template #default><el-button size="small" type="warning" :loading="rebuilding" @click="reanalyzeAndGenerateReport">重新分析并生成</el-button></template></el-alert>
      <section v-loading="markdownLoading" class="page-card report-article-section">
        <div class="report-article-heading"><div><h2>LLM 总结报告</h2><p>根据已完成 Worker 的结果、分级统计和 Finding 自动生成。</p></div><el-button text @click="router.push(`/reports/${taskId}/preview`)">单页阅读</el-button></div>
        <MarkdownArticle v-if="markdown" :content="markdown" />
        <el-empty v-else description="Markdown 报告快照暂不可用" :image-size="72" />
      </section>
      <div class="summary-grid"><el-card><el-statistic title="紧急" :value="summary.urgent"><template #suffix><el-tag type="danger">优先处理</el-tag></template></el-statistic></el-card><el-card><el-statistic title="一般" :value="summary.normal"><template #suffix><el-tag type="warning">近期整改</el-tag></template></el-statistic></el-card><el-card><el-statistic title="建议" :value="summary.suggestion"><template #suffix><el-tag type="info">持续优化</el-tag></template></el-statistic></el-card></div>
      <div class="page-card filters"><el-segmented v-model="severity" :options="[{ label: '全部', value: '' }, { label: `紧急 ${summary.urgent}`, value: 'urgent' }, { label: `一般 ${summary.normal}`, value: 'normal' }, { label: `建议 ${summary.suggestion}`, value: 'suggestion' }]" /><el-input v-model="keyword" clearable placeholder="搜索规则、对象或问题" /></div>
      <el-table :data="report.findings" class="page-card finding-table" empty-text="暂无符合条件的问题">
        <el-table-column label="等级" width="100"><template #default="{ row }"><el-tag :type="severityType[row.severity]">{{ severityLabel[row.severity] }}</el-tag></template></el-table-column>
        <el-table-column prop="objectName" label="对象" min-width="180" show-overflow-tooltip />
        <el-table-column prop="title" label="问题" min-width="220" />
        <el-table-column prop="ruleId" label="规则" width="120" />
        <el-table-column label="操作" width="100"><template #default="{ row }"><el-button link type="primary" @click="detail(row)">查看详情</el-button></template></el-table-column>
      </el-table>
      <el-pagination v-if="report.total > 0" v-model:current-page="page" v-model:page-size="pageSize" :total="report.total" :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next" class="report-pagination" @current-change="load" @size-change="() => { page = 1; load() }" />
    </template>
  </section>

  <el-drawer v-model="drawer" :title="selected?.title" size="560px"><template v-if="selected"><el-space direction="vertical" alignment="stretch" :size="18" class="drawer-content"><el-tag :type="severityType[selected.severity]">{{ severityLabel[selected.severity] }} · {{ selected.ruleId }}</el-tag><el-descriptions :column="1" border><el-descriptions-item label="对象">{{ selected.objectName }}</el-descriptions-item><el-descriptions-item label="检测证据">{{ selected.evidence }}</el-descriptions-item><el-descriptions-item label="修复建议">{{ selected.recommendation }}</el-descriptions-item><el-descriptions-item label="回滚说明">{{ selected.rollbackPlan || '需人工审核后制定回滚方案' }}</el-descriptions-item></el-descriptions><el-button :loading="explaining" @click="explain">AI 辅助解释</el-button><el-alert v-if="explanation" :title="`模型：${explanation.model}`" type="info" :closable="false"><template #default>{{ explanation.content }}</template></el-alert><template v-if="selected.dmlSuggestion"><h3>DML 修复建议（需人工复核）</h3><pre class="sql-block">{{ selected.precheckSql }}

{{ selected.dmlSuggestion }}</pre><el-button @click="copy(`${selected.precheckSql}\n${selected.dmlSuggestion}`)">复制 DML</el-button></template><template v-if="selected.ddlSuggestion"><h3>DDL 变更建议（需人工复核）</h3><pre class="sql-block">{{ selected.precheckSql }}

{{ selected.ddlSuggestion }}</pre><el-button @click="copy(`${selected.precheckSql}\n${selected.ddlSuggestion}`)">复制 DDL</el-button></template></el-space></template></el-drawer>
</template>
