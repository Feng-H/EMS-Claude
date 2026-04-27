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
      <mobile-header title="EMS 工作台" :show-back="false" />

      <div class="content">
        <!-- 待办数据统计 -->
        <div class="stats-overview">
          <div class="stats-card" @click="navigateTo('/h5/inspection')">
            <div class="stats-num">{{ pendingInspections }}</div>
            <div class="stats-label">待点检</div>
          </div>
          <div class="stats-card" @click="navigateTo('/h5/maintenance')">
            <div class="stats-num">{{ pendingMaintenance }}</div>
            <div class="stats-label">待保养</div>
          </div>
          <div class="stats-card" @click="navigateTo('/h5/repair/execute')">
            <div class="stats-num">{{ pendingRepairs }}</div>
            <div class="stats-label">待维修</div>
          </div>
        </div>

        <!-- 快捷功能 -->
        <van-cell-group inset class="section-title">
          <van-cell title="核心作业" />
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

      <!-- 辅助功能 -->
      <van-cell-group inset class="section-title">
        <van-cell title="辅助查询" />
      </van-cell-group>

      <van-cell-group inset>
        <van-cell
          title="知识库"
          label="故障处理经验沉淀"
          is-link
          icon="description-o"
          @click="navigateTo('/h5/knowledge')"
        />
        <van-cell
          title="我的设备"
          label="快速搜索与查看设备"
          is-link
          icon="cluster-o"
          @click="navigateTo('/h5/equipment')"
        />
        <van-cell
          title="备件库存"
          label="备件余量即时查询"
          is-link
          icon="shop-o"
          @click="navigateTo('/h5/spareparts')"
        />
      </van-cell-group>
    </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
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
  const tokenInStorage = localStorage.getItem('ems_token')
  if (tokenInStorage && !authStore.user) {
    try {
      await authStore.getUserInfo()
    } catch (error) {
      authStore.logout()
    }
  }
}

const navigateTo = (path: string) => {
  router.push(path)
}

const goToLogin = () => {
  authStore.logout()
  window.location.href = '/login'
}

const loadPendingCounts = async () => {
  try {
    const inspectionsResponse = await inspectionTaskApi.getMyTasks()
    pendingInspections.value = inspectionsResponse.data.filter(t => t.status === 'pending' || t.status === 'in_progress').length

    const maintenanceResponse = await maintenanceApi.getMyTasks()
    pendingMaintenance.value = maintenanceResponse.data.filter(t => t.status === 'pending' || t.status === 'in_progress').length

    const repairsResponse = await repairOrderApi.getMyTasks()
    pendingRepairs.value = repairsResponse.data.filter(r => r.status === 'pending' || r.status === 'assigned' || r.status === 'in_progress').length
  } catch (error) {
    console.error('加载待办数量失败:', error)
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
  background: var(--color-bg-primary);
}

/* 未登录状态 */
.login-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 40px;
}

.login-btn {
  width: 200px;
  height: 50px;
  border-radius: var(--radius-high);
  background: var(--color-terracotta);
  border: none;
  font-size: 16px;
  font-weight: 500;
  box-shadow: 0 8px 24px rgba(201, 100, 66, 0.3);
  color: #faf9f5;
}

/* 已登录状态 */
.logged-in-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 80px;
}

.content {
  padding: var(--space-lg);
}

/* 统计卡片优化 */
.stats-overview {
  display: flex;
  gap: var(--space-md);
  margin-bottom: var(--space-xl);
}

.stats-card {
  flex: 1;
  background: var(--color-bg-card);
  padding: var(--space-lg) var(--space-sm);
  border-radius: var(--radius-very);
  text-align: center;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  transition: all var(--transition-fast);
}

.stats-card:active {
  transform: scale(0.95);
  background: var(--color-bg-tertiary);
}

.stats-num {
  font-size: 24px;
  font-weight: 600;
  color: var(--color-terracotta);
  margin-bottom: 4px;
  font-family: var(--font-serif);
}

.stats-label {
  font-size: 12px;
  color: var(--color-text-secondary);
}

/* Section 标题 */
.section-title :deep(.van-cell) {
  background: transparent;
  padding: var(--space-md) 0 var(--space-sm);
}

.section-title :deep(.van-cell__title) {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  font-family: var(--font-serif);
}

/* 快捷操作网格 */
.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-md);
  padding: var(--space-xl) var(--space-md);
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  margin-bottom: var(--space-xl);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.action-item :deep(.van-icon) {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-tertiary);
  border-radius: 14px;
  color: var(--color-terracotta);
  font-size: 24px;
}

.action-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  font-weight: 500;
}

/* 列表项优化 */
.content :deep(.van-cell-group) {
  background: transparent;
}

.content :deep(.van-cell) {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  margin-bottom: var(--space-sm);
  padding: 16px;
  border: 1px solid var(--color-border);
  box-shadow: var(--shadow-sm);
}

.content :deep(.van-cell::after) {
  display: none;
}

.content :deep(.van-cell__title) {
  font-weight: 500;
}
</style>
