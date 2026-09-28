<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <section class="page-card not-found"><el-result icon="warning" title="页面不存在" sub-title="该地址无效，或页面已被移除。"><template #extra><el-button type="primary" @click="router.replace('/')">返回首页</el-button></template></el-result></section>
</template>
