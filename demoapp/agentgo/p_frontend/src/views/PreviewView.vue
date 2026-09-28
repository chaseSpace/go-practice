<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getPreview } from '@/api/tasks'
import MarkdownArticle from '@/components/MarkdownArticle.vue'

const props = defineProps<{ taskId: string }>()
const router = useRouter()
const viewMode = ref<'article' | 'source'>('article')
const content = ref('')
const loading = ref(false)
const error = ref('')

async function load() {
  loading.value = true
  try { content.value = await getPreview(props.taskId, 'markdown'); error.value = '' } catch { error.value = '报告尚未生成或任务已过期。' } finally { loading.value = false }
}

watch(() => props.taskId, load)
onMounted(load)
</script>

<template>
  <section class="preview-page">
    <div class="report-toolbar"><div><el-button text @click="router.push(`/reports/${taskId}`)">← 返回问题列表</el-button><h1>在线报告预览</h1></div><el-radio-group v-model="viewMode"><el-radio-button value="article">阅读预览</el-radio-button><el-radio-button value="source">Markdown 源码</el-radio-button></el-radio-group></div>
    <el-alert v-if="error" :title="error" type="warning" :closable="false" />
    <div v-loading="loading" class="page-card preview-content">
      <MarkdownArticle v-if="viewMode === 'article' && content" :content="content" />
      <pre v-else-if="content" class="markdown-content">{{ content }}</pre>
    </div>
  </section>
</template>
