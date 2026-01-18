import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false, layout: 'full-screen' },
  },
  {
    path: '/change-password',
    name: 'ChangePassword',
    component: () => import('@/views/auth/ChangePasswordView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
  },
  // H5移动端首页（全屏布局）
  {
    path: '/h5',
    name: 'H5Home',
    component: () => import('@/views/h5/IndexView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
  },
  // H5移动端路由
  {
    path: '/h5/inspection',
    name: 'H5Inspection',
    component: () => import('@/views/inspection/ExecuteView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
  },
  {
    path: '/h5/repair/report',
    name: 'H5RepairReport',
    component: () => import('@/views/repair/ReportView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
  },
  {
    path: '/h5/repair/execute',
    name: 'H5RepairExecute',
    component: () => import('@/views/repair/ExecuteView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
  },
  {
    path: '/h5/maintenance',
    name: 'H5Maintenance',
    component: () => import('@/views/maintenance/ExecuteView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
  },
  {
    path: '/h5/tasks',
    name: 'H5Tasks',
    component: () => import('@/views/h5/TasksView.vue'),
    meta: { requiresAuth: true, layout: 'full-screen' },
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
      {
        path: 'organization/users',
        name: 'UserManagement',
        component: () => import('@/views/user/UserManagementView.vue'),
        meta: { title: '人员管理', roles: ['admin'] },
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
      // 维修管理路由
      {
        path: 'repair/orders',
        name: 'RepairOrders',
        component: () => import('@/views/repair/OrderListView.vue'),
        meta: { title: '维修工单' },
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

  // 允许通过 ?force=true 参数强制访问登录页
  if (to.name === 'Login' && authStore.isLoggedIn && !to.query.force) {
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
