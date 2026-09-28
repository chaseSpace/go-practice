<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const title = computed(() => route.meta.title ?? 'DB Inspector')
</script>

<template>
  <el-container class="app-shell">
    <el-header class="app-header">
      <RouterLink to="/" class="brand">DB Inspector</RouterLink>
      <span class="page-title">{{ title }}</span>
      <nav><RouterLink to="/history">最近任务</RouterLink></nav>
    </el-header>
    <el-main class="app-main"><RouterView /></el-main>
  </el-container>
</template>
