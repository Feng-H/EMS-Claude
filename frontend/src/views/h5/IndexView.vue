<template>
  <div class="h5-home-view">
    <!-- 未登录状态 -->
    <div v-if="!isLoggedIn || !authStore.user" class="login-prompt">
      <div class="login-icon">
        <van-icon name="lock" size="80" color="#999" />
      </div>
      <div class="login-text">请先登录</div>
      <van-button type="primary" size="large" block @click="goToLogin" class="login-btn">
        去登录
      </van-button>
    </div>

    <!-- 已登录状态 -->
    <div v-else class="logged-in-view">
      <mobile-header title="EMS移动端" :show-back="false">
        <template #right>
          <van-icon name="switch" size="20" @click="handleLogout" />
        </template>
      </mobile-header>

      <div class="content">
      <!-- 用户信息卡片 -->
      <van-cell-group inset class="user-card">
        <van-cell center>
          <template #title>
            <div class="user-info">
              <van-icon name="user-circle-o" size="48" color="#1989fa" />
              <div class="user-details">
                <div class="user-name">{{ authStore.user?.name || '未登录' }}</div>
                <div class="user-role">{{ getRoleText(authStore.user?.role) }}</div>
              </div>
            </div>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 快捷功能 -->
      <van-cell-group inset class="section">
        <van-cell title="快捷操作" />
      </van-cell-group>

      <div class="action-grid">
        <div class="action-item" @click="navigateTo('/h5/inspection')">
          <van-icon name="scan" size="32" color="#1989fa" />
          <span class="action-label">设备点检</span>
        </div>
        <div class="action-item" @click="navigateTo('/h5/repair/report')">
          <van-icon name="setting-o" size="32" color="#ff976a" />
          <span class="action-label">故障报修</span>
        </div>
        <div class="action-item" @click="navigateTo('/h5/maintenance')">
          <van-icon name="todo-list-o" size="32" color="#07c160" />
          <span class="action-label">保养任务</span>
        </div>
        <div class="action-item" @click="navigateTo('/h5/repair/execute')">
          <van-icon name="tool-job-o" size="32" color="#ee0a24" />
          <span class="action-label">维修执行</span>
        </div>
      </div>

      <!-- 待办事项 -->
      <van-cell-group inset class="section">
        <van-cell title="我的任务" is-link @click="navigateTo('/h5/tasks')" />
      </van-cell-group>

      <van-cell-group inset>
        <van-cell
          title="待处理点检"
          :value="`${pendingInspections}项`"
          is-link
          @click="navigateTo('/h5/inspection')"
        />
        <van-cell
          title="待处理保养"
          :value="`${pendingMaintenance}项`"
          is-link
          @click="navigateTo('/h5/maintenance')"
        />
        <van-cell
          title="待处理维修"
          :value="`${pendingRepairs}项`"
          is-link
          @click="navigateTo('/h5/repair/execute')"
        />
      </van-cell-group>

      <!-- 其他功能 -->
      <van-cell-group inset class="section">
        <van-cell title="其他" />
      </van-cell-group>

      <van-cell-group inset>
        <van-cell
          title="知识库"
          label="故障处理方法"
          is-link
          @click="navigateTo('/knowledge')"
        />
        <van-cell
          title="我的设备"
          label="查看设备列表"
          is-link
          @click="navigateTo('/equipment')"
        />
        <van-cell
          title="备件库存"
          label="查看备件信息"
          is-link
          @click="navigateTo('/spareparts')"
        />
        <van-cell
          title="退出登录"
          @click="handleLogout"
        >
          <template #right-icon>
            <van-icon name="switch" color="#ee0a24" />
          </template>
        </van-cell>
      </van-cell-group>
    </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useAuthStore } from '@/stores/auth'
import { inspectionTaskApi } from '@/api/inspection'
import { maintenanceApi } from '@/api/maintenance'
import { repairOrderApi } from '@/api/repair'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const router = useRouter()
const authStore = useAuthStore()

const isLoggedIn = computed(() => authStore.isLoggedIn)

const pendingInspections = ref(0)
const pendingMaintenance = ref(0)
const pendingRepairs = ref(0)

// 初始化：检查登录状态并获取用户信息
const initializeAuth = async () => {
  // 检查localStorage中的token是否有效
  const tokenInStorage = localStorage.getItem('ems_token')

  if (tokenInStorage && !authStore.user) {
    try {
      await authStore.getUserInfo()
    } catch (error) {
      // 获取用户信息失败，清除登录状态
      console.error('获取用户信息失败:', error)
      authStore.logout()
    }
  } else if (!tokenInStorage && authStore.isLoggedIn) {
    // localStorage中没有token但store认为已登录，同步状态
    authStore.logout()
  }
}

const getRoleText = (role?: string) => {
  const roles: Record<string, string> = {
    admin: '系统管理员',
    supervisor: '设备主管',
    engineer: '设备工程师',
    maintenance: '维修工',
    operator: '操作工'
  }
  return roles[role || ''] || '未知角色'
}

const navigateTo = (path: string) => {
  router.push(path)
}

const goToLogin = () => {
  // 清除登录状态
  authStore.logout()
  // 使用window.location确保完整页面跳转
  window.location.href = '/login'
}

const loadPendingCounts = async () => {
  try {
    // 加载待处理点检
    const inspections = await inspectionTaskApi.getMyTasks()
    pendingInspections.value = inspections.filter(t => t.status === 'pending' || t.status === 'in_progress').length

    // 加载待处理保养
    const maintenance = await maintenanceApi.getMyTasks()
    pendingMaintenance.value = maintenance.filter(t => t.status === 'pending' || t.status === 'in_progress').length

    // 加载待处理维修
    const repairs = await repairOrderApi.getMyOrders()
    pendingRepairs.value = repairs.filter(r => r.status === 'pending' || r.status === 'in_progress').length
  } catch (error) {
    console.error('加载待办数量失败:', error)
  }
}

const handleLogout = async () => {
  try {
    await showConfirmDialog({
      title: '退出登录',
      message: '确定要退出登录吗？'
    })
    authStore.logout()
    showToast('已退出登录')
    router.push('/login')
  } catch {
    // 用户取消
  }
}

onMounted(async () => {
  await initializeAuth()
  if (authStore.user) {
    loadPendingCounts()
  }
})
</script>

<style scoped>
.h5-home-view {
  min-height: 100vh;
  background: linear-gradient(180deg, #F8FAFC 0%, #F1F5F9 100%);
}

/* 未登录状态 */
.login-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 40px;
  background: linear-gradient(180deg, #F8FAFC 0%, #F1F5F9 100%);
}

.login-icon {
  margin-bottom: 24px;
  opacity: 0.6;
}

.login-text {
  font-size: 17px;
  color: #64748B;
  margin-bottom: 32px;
  font-weight: 500;
}

.login-btn {
  width: 200px;
  height: 50px;
  border-radius: 14px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.3);
}

.login-btn:active {
  transform: scale(0.96);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* 已登录状态 */
.logged-in-view {
  min-height: 100vh;
  background: linear-gradient(180deg, #F8FAFC 0%, #F1F5F9 100%);
  padding-top: 46px;
  padding-bottom: env(safe-area-inset-bottom, 20px);
}

.content {
  padding: 16px;
}

/* 用户卡片优化 */
.user-card {
  margin-bottom: 20px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
}

.user-card :deep(.van-cell-group) {
  background: transparent;
}

.user-card :deep(.van-cell) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  color: #FFFFFF;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
}

.user-info :deep(.van-icon) {
  color: rgba(255, 255, 255, 0.9);
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 18px;
  font-weight: 700;
  margin-bottom: 4px;
  color: #FFFFFF;
}

.user-role {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
}

/* Section 标题优化 */
.section :deep(.van-cell) {
  background: transparent;
  padding: 12px 16px 8px;
}

.section :deep(.van-cell__title) {
  font-size: 15px;
  font-weight: 600;
  color: #1A202C;
}

/* 快捷操作网格优化 */
.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 20px 16px;
  background: #FFFFFF;
  border-radius: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  cursor: pointer;
  transition: all 0.2s;
}

.action-item :deep(.van-icon) {
  width: 52px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #F8FAFC, #F1F5F9);
  border-radius: 14px;
  transition: all 0.2s;
}

.action-item:active :deep(.van-icon) {
  transform: scale(0.92);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.15), rgba(118, 75, 162, 0.15));
}

.action-item:active .action-label {
  color: #667eea;
}

.action-label {
  font-size: 12px;
  color: #64748B;
  text-align: center;
  font-weight: 500;
  transition: color 0.2s;
}

/* Cell Group 优化 */
.content :deep(.van-cell-group) {
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.content :deep(.van-cell) {
  background: #FFFFFF;
  padding: 14px 16px;
}

.content :deep(.van-cell:not(:last-child)::after) {
  left: 16px;
  border-color: #F1F5F9;
}

.content :deep(.van-cell__title) {
  font-size: 15px;
  font-weight: 500;
  color: #1A202C;
}

.content :deep(.van-cell__label) {
  font-size: 13px;
  color: #94A3B8;
}

.content :deep(.van-cell__value) {
  font-size: 14px;
  color: #667eea;
  font-weight: 500;
}

.content :deep(.van-icon__image) {
  color: #667eea;
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  .h5-home-view,
  .login-prompt {
    background: linear-gradient(180deg, #0F1419 0%, #1A2332 100%);
  }

  .login-text {
    color: #6B7280;
  }

  .action-grid,
  .content :deep(.van-cell) {
    background: #1A2332;
  }

  .section :deep(.van-cell__title),
  .content :deep(.van-cell__title) {
    color: #F0F3FF;
  }

  .content :deep(.van-cell__label) {
    color: #6B7280;
  }

  .content :deep(.van-cell:not(:last-child)::after) {
    border-color: #2D3748;
  }

  .action-label {
    color: #6B7280;
  }

  .action-item :deep(.van-icon) {
    background: linear-gradient(135deg, #2D3748, #1F2937);
  }
}
</style>
