<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { createTask } from '@/api/tasks'
import { saveTaskToken } from '@/api/client'
import { useHistoryStore } from '@/stores/history'
import type { ConnectionConfig, CreateTaskPayload } from '@/types/domain'
import { buildTaskPayload, connectionDatabaseName, connectionDialect, isValidConnection, isValidDatabaseURL, maskedConnectionURL } from '@/utils/connection'

type ConnectionMode = 'url' | 'fields'

const router = useRouter()
const history = useHistoryStore()
const submitting = ref(false)
const mode = ref<ConnectionMode>('url')
const form = reactive({
  connectionUrl: '',
  connection: { dialect: 'mysql', host: '127.0.0.1', port: 3306, database: '', username: 'readonly', password: '' } as ConnectionConfig,
  excludeTables: '',
  redactObjectNames: false,
})

function clearInactiveCredentials() {
  if (mode.value === 'url') {
    form.connection.password = ''
  } else {
    form.connectionUrl = ''
  }
}

function updateDefaultPort() {
  form.connection.port = form.connection.dialect === 'postgres' ? 5432 : 3306
}

function validateCurrentConnection() {
  if (mode.value === 'url') {
    if (!isValidDatabaseURL(form.connectionUrl)) {
      ElMessage.warning('请输入有效的数据库 URL，例如 mysql://用户名:密码@主机:3306/数据库。')
      return false
    }
    return true
  }
  if (!isValidConnection(form.connection)) {
    ElMessage.warning('请完整填写主机、端口、数据库和用户名。')
    return false
  }
  return true
}

async function submit() {
  if (!validateCurrentConnection()) return
  submitting.value = true
  try {
    const payload: CreateTaskPayload = buildTaskPayload(mode.value, form.connectionUrl, form.connection, form.excludeTables, form.redactObjectNames)
    const maskedUrl = maskedConnectionURL(mode.value, form.connectionUrl, form.connection)
    const database = connectionDatabaseName(mode.value, form.connectionUrl, form.connection)
    const dialect = connectionDialect(mode.value, form.connectionUrl, form.connection)
    const created = await createTask(payload)
    saveTaskToken(created.taskId, created.taskAccessToken)
    history.remember({ id: created.taskId, database, dialect, maskedUrl, createdAt: new Date().toISOString(), status: 'pending' })
    form.connectionUrl = ''
    form.connection.password = ''
    await router.push(`/tasks/${created.taskId}`)
  } catch {
    ElMessage.error('创建任务失败，请检查连接信息和后端服务。')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <section class="page-card connection-page">
    <div class="page-intro">
      <h1>数据库体检</h1>
      <p>使用只读账号检查 Schema。不会读取业务数据，也不会直接执行数据库修改。</p>
    </div>
    <el-alert title="密码和完整数据库 URL 仅用于本次请求；最近任务仅保存密码遮掩后的 URL 标签。" type="info" :closable="false" show-icon />

    <el-form :model="form" label-position="top" class="connection-form" @submit.prevent>
      <el-tabs v-model="mode" @tab-change="clearInactiveCredentials">
        <el-tab-pane label="数据库 URL" name="url">
          <el-form-item label="数据库 URL">
            <el-input v-model="form.connectionUrl" autocomplete="off" placeholder="mysql://用户名:密码@主机:3306/数据库" />
          </el-form-item>
          <p class="field-help">格式：mysql://用户名:密码@主机:端口/数据库</p>
        </el-tab-pane>
        <el-tab-pane label="逐项填写" name="fields">
          <div class="form-grid">
            <el-form-item label="数据库类型"><el-select v-model="form.connection.dialect" @change="updateDefaultPort"><el-option label="MySQL" value="mysql" /><el-option label="PostgreSQL" value="postgres" /></el-select></el-form-item>
            <el-form-item label="主机"><el-input v-model="form.connection.host" autocomplete="off" /></el-form-item>
            <el-form-item label="端口"><el-input-number v-model="form.connection.port" :min="1" :max="65535" controls-position="right" /></el-form-item>
            <el-form-item label="数据库"><el-input v-model="form.connection.database" autocomplete="off" /></el-form-item>
            <el-form-item label="用户名"><el-input v-model="form.connection.username" autocomplete="off" /></el-form-item>
            <el-form-item label="密码"><el-input v-model="form.connection.password" type="password" show-password autocomplete="new-password" /></el-form-item>
          </div>
        </el-tab-pane>
      </el-tabs>

      <div class="form-grid options-grid">
        <el-form-item label="排除表（可选，英文逗号分隔）"><el-input v-model="form.excludeTables" placeholder="audit_log, temporary_data" /></el-form-item>
        <el-form-item label="报告脱敏"><el-switch v-model="form.redactObjectNames" active-text="隐藏对象名称" /></el-form-item>
      </div>
      <div class="form-actions"><el-button type="primary" :loading="submitting" @click="submit">开始检查</el-button></div>
    </el-form>
  </section>
</template>
