import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/views/layout/MainLayout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/DashboardView.vue'),
        meta: { title: '首页' },
      },
      {
        path: 'equipment',
        name: 'Equipment',
        component: () => import('@/views/equipment/EquipmentListView.vue'),
        meta: { title: '设备台账' },
      },
      {
        path: 'equipment/detail/:id',
        name: 'EquipmentDetail',
        component: () => import('@/views/equipment/EquipmentDetailView.vue'),
        meta: { title: '设备详情' },
      },
      {
        path: 'organization',
        name: 'Organization',
        component: () => import('@/views/equipment/OrganizationView.vue'),
        meta: { title: '组织架构', roles: ['admin', 'engineer'] },
      },
      // 点检管理路由
      {
        path: 'inspection/templates',
        name: 'InspectionTemplates',
        component: () => import('@/views/inspection/TemplateView.vue'),
        meta: { title: '点检模板', roles: ['admin', 'engineer'] },
      },
      {
        path: 'inspection/tasks',
        name: 'InspectionTasks',
        component: () => import('@/views/inspection/TaskListView.vue'),
        meta: { title: '点检任务' },
      },
      // 移动端点检执行（独立页面，不需要主布局）
      {
        path: 'inspection/execute/:taskId?',
        name: 'inspection-execute',
        component: () => import('@/views/inspection/ExecuteView.vue'),
        meta: { requiresAuth: true, layout: 'full-screen' },
      },
      // 维修管理路由
      {
        path: 'repair/orders',
        name: 'RepairOrders',
        component: () => import('@/views/repair/OrderListView.vue'),
        meta: { title: '维修工单' },
      },
      // 移动端报修页面（独立页面）
      {
        path: 'repair/create',
        name: 'repair-create',
        component: () => import('@/views/repair/ReportView.vue'),
        meta: { requiresAuth: true, layout: 'full-screen' },
      },
      // 移动端维修执行页面（独立页面）
      {
        path: 'repair/execute/:orderId?',
        name: 'repair-execute',
        component: () => import('@/views/repair/ExecuteView.vue'),
        meta: { requiresAuth: true, layout: 'full-screen' },
      },
      // 保养管理路由
      {
        path: 'maintenance',
        name: 'Maintenance',
        redirect: '/maintenance/plans',
        meta: { title: '保养管理' },
      },
      {
        path: 'maintenance/plans',
        name: 'MaintenancePlans',
        component: () => import('@/views/maintenance/PlanListView.vue'),
        meta: { title: '保养计划', roles: ['admin', 'engineer'] },
      },
      {
        path: 'maintenance/tasks',
        name: 'MaintenanceTasks',
        component: () => import('@/views/maintenance/TaskListView.vue'),
        meta: { title: '保养任务' },
      },
      // 移动端保养执行页面（独立页面）
      {
        path: 'maintenance/execute/:taskId?',
        name: 'maintenance-execute',
        component: () => import('@/views/maintenance/ExecuteView.vue'),
        meta: { requiresAuth: true, layout: 'full-screen' },
      },
      // 备件管理路由
      {
        path: 'spareparts',
        name: 'SpareParts',
        component: () => import('@/views/sparepart/SparePartListView.vue'),
        meta: { title: '备件管理' },
      },
      // 统计分析路由
      {
        path: 'analytics',
        name: 'Analytics',
        component: () => import('@/views/analytics/AnalyticsView.vue'),
        meta: { title: '统计分析' },
      },
      // 知识库路由
      {
        path: 'knowledge',
        name: 'Knowledge',
        component: () => import('@/views/knowledge/KnowledgeView.vue'),
        meta: { title: '知识库' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta?.requiresAuth !== false
  const requiredRoles = to.meta?.roles as string[] | undefined

  if (requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }

  if (to.name === 'Login' && authStore.isLoggedIn) {
    next({ name: 'Dashboard' })
    return
  }

  if (requiredRoles && requiredRoles.length > 0) {
    if (!authStore.hasRole(...requiredRoles)) {
      next({ name: 'Dashboard' })
      return
    }
  }

  next()
})

export default router
