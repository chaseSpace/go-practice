import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import TaskView from '@/views/TaskView.vue'
import ReportView from '@/views/ReportView.vue'
import PreviewView from '@/views/PreviewView.vue'
import HistoryView from '@/views/HistoryView.vue'
import NotFoundView from '@/views/NotFoundView.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomeView, meta: { title: '新建检查' } },
    { path: '/tasks/:taskId', component: TaskView, props: true, meta: { title: '检查进度' } },
    { path: '/reports/:taskId', component: ReportView, props: true, meta: { title: '体检报告' } },
    { path: '/reports/:taskId/preview', component: PreviewView, props: true, meta: { title: '报告预览' } },
    { path: '/history', component: HistoryView, meta: { title: '最近任务' } },
    { path: '/:pathMatch(.*)*', component: NotFoundView, meta: { title: '页面不存在' } },
  ],
})
