<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section animate-fade-in">
      <div class="welcome-content">
        <h1 class="welcome-title">
          <span class="title-greeting">欢迎回来，</span>
          <span class="title-name">{{ authStore.userName }}</span>
        </h1>
        <p class="welcome-subtitle">EMS 设备管理系统 · 实时监控 · 智能管理</p>
      </div>
      <div class="welcome-time">
        <div class="time-display">{{ currentTime }}</div>
        <div class="date-display">{{ currentDate }}</div>
      </div>
    </div>

    <!-- 设备统计卡片 -->
    <div class="stats-grid animate-fade-in animate-delay-1">
      <div class="stat-card" data-type="total">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
              <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
              <line x1="12" y1="22.08" x2="12" y2="12"/>
            </svg>
          </div>
          <div class="stat-trend">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.total }}</div>
        <div class="stat-label">设备总数</div>
        <div class="stat-bar">
          <div class="bar-fill" style="width: 100%"></div>
        </div>
      </div>

      <div class="stat-card" data-type="running">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polygon points="5 3 19 12 5 21 5 3"/>
            </svg>
          </div>
          <div class="stat-trend up">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.running }}</div>
        <div class="stat-label">运行中</div>
        <div class="stat-bar">
          <div class="bar-fill success" :style="{ width: runningPercent + '%' }"></div>
        </div>
      </div>

      <div class="stat-card" data-type="maintenance">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
            </svg>
          </div>
          <div class="stat-trend">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="6" x2="12" y2="12"/>
              <line x1="12" y1="12" x2="16" y2="14"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.maintenance }}</div>
        <div class="stat-label">维修中</div>
        <div class="stat-bar">
          <div class="bar-fill warning" :style="{ width: maintenancePercent + '%' }"></div>
        </div>
      </div>

      <div class="stat-card" data-type="stopped">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="6" y="4" width="4" height="16"/>
              <rect x="14" y="4" width="4" height="16"/>
            </svg>
          </div>
          <div class="stat-trend down">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/>
              <polyline points="17 18 23 18 23 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.stopped }}</div>
        <div class="stat-label">已停机</div>
        <div class="stat-bar">
          <div class="bar-fill danger" :style="{ width: stoppedPercent + '%' }"></div>
        </div>
      </div>
    </div>

    <!-- 点检统计面板 -->
    <div class="inspection-panel animate-fade-in animate-delay-2">
      <div class="panel-header">
        <h3 class="panel-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          点检统计
        </h3>
        <router-link to="/inspection/tasks" class="panel-link">
          查看全部
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </router-link>
      </div>

      <div class="inspection-stats">
        <div class="istat-item">
          <div class="istat-value">{{ inspectionStats.total_tasks }}</div>
          <div class="istat-label">总任务数</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value pending">{{ inspectionStats.pending_tasks }}</div>
          <div class="istat-label">待执行</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value progress">{{ inspectionStats.in_progress_tasks }}</div>
          <div class="istat-label">进行中</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value completed">{{ inspectionStats.completed_tasks }}</div>
          <div class="istat-label">已完成</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value today">{{ inspectionStats.today_completed }}</div>
          <div class="istat-label">今日完成</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item highlight">
          <div class="istat-value">{{ inspectionStats.completion_rate.toFixed(1) }}%</div>
          <div class="istat-label">完成率</div>
        </div>
      </div>

      <!-- 进度条 -->
      <div class="progress-track">
        <div
          class="progress-fill"
          :style="{ width: inspectionStats.completion_rate + '%' }"
        ></div>
      </div>
    </div>

    <!-- 下方区域 -->
    <div class="dashboard-grid animate-fade-in animate-delay-3">
      <!-- 快速操作 -->
      <div class="dashboard-card quick-actions">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          快速操作
        </div>
        <div class="action-list">
          <router-link to="/equipment" class="action-item">
            <div class="action-icon primary">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">设备台账</div>
              <div class="action-desc">查看所有设备信息</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>

          <router-link to="/inspection/execute" class="action-item">
            <div class="action-icon success">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">开始点检</div>
              <div class="action-desc">执行设备点检任务</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>

          <router-link to="/inspection/tasks" class="action-item">
            <div class="action-icon warning">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <line x1="8" y1="6" x2="21" y2="6"/>
                <line x1="8" y1="12" x2="21" y2="12"/>
                <line x1="8" y1="18" x2="21" y2="18"/>
                <line x1="3" y1="6" x2="3.01" y2="6"/>
                <line x1="3" y1="12" x2="3.01" y2="12"/>
                <line x1="3" y1="18" x2="3.01" y2="18"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">点检任务</div>
              <div class="action-desc">管理点检任务列表</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>

          <router-link to="/analytics" class="action-item">
            <div class="action-icon info">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <line x1="18" y1="20" x2="18" y2="10"/>
                <line x1="12" y1="20" x2="12" y2="4"/>
                <line x1="6" y1="20" x2="6" y2="14"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">数据统计</div>
              <div class="action-desc">查看数据分析报告</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>
        </div>
      </div>

      <!-- 我的今日任务 -->
      <div class="dashboard-card my-tasks">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
          我的今日任务
        </div>
        <div v-if="myStats.loading" class="tasks-loading">
          <div class="loading-spinner"></div>
        </div>
        <div v-else class="tasks-summary">
          <div class="task-metric">
            <div class="metric-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value">{{ myStats.pending_count }}</div>
              <div class="metric-label">待执行</div>
            </div>
          </div>
          <div class="task-metric">
            <div class="metric-icon active">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value">{{ myStats.in_progress_count }}</div>
              <div class="metric-label">进行中</div>
            </div>
          </div>
          <div class="task-metric">
            <div class="metric-icon success">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value">{{ myStats.today_tasks }}</div>
              <div class="metric-label">今日完成</div>
            </div>
          </div>
        </div>
        <router-link to="/inspection/execute" class="tasks-action">
          <span>开始执行</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </router-link>
      </div>

      <!-- 待办事项 -->
      <div class="dashboard-card todo-list">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
            <polyline points="10 9 9 9 8 9"/>
          </svg>
          待办事项
        </div>
        <div v-if="pendingTasks.length === 0" class="todo-empty">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <p>暂无待办任务</p>
        </div>
        <div v-else class="todo-items">
          <router-link
            v-for="task in pendingTasks.slice(0, 4)"
            :key="task.id"
            :to="`/inspection/execute/${task.id}`"
            class="todo-item"
          >
            <div class="todo-status"></div>
            <div class="todo-content">
              <div class="todo-title">{{ task.equipment_name || task.equipment_code }}</div>
              <div class="todo-time">{{ task.scheduled_date }}</div>
            </div>
            <div class="todo-tag">待点检</div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { equipmentApi, type EquipmentStatistics } from '@/api/equipment'
import {
  inspectionTaskApi,
  type InspectionStatistics,
  type MyTasksStatistics,
  type InspectionTask
} from '@/api/inspection'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const statistics = ref<EquipmentStatistics>({
  total: 0,
  running: 0,
  stopped: 0,
  maintenance: 0,
  scrapped: 0,
})

const inspectionStats = ref<InspectionStatistics>({
  total_tasks: 0,
  pending_tasks: 0,
  in_progress_tasks: 0,
  completed_tasks: 0,
  overdue_tasks: 0,
  today_completed: 0,
  completion_rate: 0,
})

const myStats = ref<{
  pending_count: number
  in_progress_count: number
  today_tasks: number
  loading: boolean
}>({
  pending_count: 0,
  in_progress_count: 0,
  today_tasks: 0,
  loading: true
})

const pendingTasks = ref<InspectionTask[]>([])
const currentTime = ref('')
const currentDate = ref('')

// 计算百分比
const runningPercent = computed(() => {
  if (statistics.value.total === 0) return 0
  return (statistics.value.running / statistics.value.total) * 100
})

const maintenancePercent = computed(() => {
  if (statistics.value.total === 0) return 0
  return (statistics.value.maintenance / statistics.value.total) * 100
})

const stoppedPercent = computed(() => {
  if (statistics.value.total === 0) return 0
  return (statistics.value.stopped / statistics.value.total) * 100
})

// 更新时间
function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  currentDate.value = now.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'long' })
}

async function loadStatistics() {
  try {
    statistics.value = await equipmentApi.getStatistics()
  } catch (error) {
    ElMessage.error('获取设备统计数据失败')
  }
}

async function loadInspectionStats() {
  try {
    inspectionStats.value = await inspectionTaskApi.getStatistics()
  } catch (error) {
    console.error('获取点检统计失败', error)
  }
}

async function loadMyStats() {
  myStats.value.loading = true
  try {
    const data = await inspectionTaskApi.getMyStats()
    myStats.value = {
      ...data,
      loading: false
    }
  } catch (error) {
    console.error('获取我的任务统计失败', error)
    myStats.value.loading = false
  }
}

async function loadPendingTasks() {
  try {
    const data = await inspectionTaskApi.getMyTasks()
    pendingTasks.value = data.filter(t => t.status === 'pending')
  } catch (error) {
    console.error('获取待办任务失败', error)
  }
}

let timeInterval: ReturnType<typeof setInterval>

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  loadStatistics()
  loadInspectionStats()
  loadMyStats()
  loadPendingTasks()
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

/* 欢迎区域 */
.welcome-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
  padding: 24px 32px;
  background: linear-gradient(135deg, var(--color-bg-card), var(--color-bg-elevated));
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  position: relative;
  overflow: hidden;
}

.welcome-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--color-primary), var(--color-success));
}

.welcome-content {
  position: relative;
  z-index: 1;
}

.welcome-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
}

.title-greeting {
  color: var(--color-text-secondary);
}

.title-name {
  background: linear-gradient(135deg, var(--color-primary), var(--color-success));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.welcome-subtitle {
  font-size: 14px;
  color: var(--color-text-tertiary);
  letter-spacing: 2px;
  text-transform: uppercase;
}

.welcome-time {
  text-align: right;
}

.time-display {
  font-family: var(--font-numbers);
  font-size: 36px;
  font-weight: 700;
  color: var(--color-primary);
  line-height: 1;
}

.date-display {
  font-size: 13px;
  color: var(--color-text-tertiary);
  margin-top: 4px;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  position: relative;
  padding: 24px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  transition: all var(--transition-base);
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: var(--color-primary);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-card), var(--shadow-glow);
  border-color: var(--color-primary-dim);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-card[data-type="total"]::before {
  background: linear-gradient(90deg, #667eea, #764ba2);
}

.stat-card[data-type="running"]::before {
  background: var(--color-success);
}

.stat-card[data-type="maintenance"]::before {
  background: var(--color-warning);
}

.stat-card[data-type="stopped"]::before {
  background: var(--color-danger);
}

.stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  color: var(--color-primary);
}

.stat-icon svg {
  width: 24px;
  height: 24px;
}

.stat-trend {
  padding: 4px;
  color: var(--color-text-tertiary);
}

.stat-trend svg {
  width: 18px;
  height: 18px;
}

.stat-trend.up {
  color: var(--color-success);
}

.stat-trend.down {
  color: var(--color-danger);
}

.stat-value {
  font-family: var(--font-numbers);
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--color-text-tertiary);
  margin-top: 4px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.stat-bar {
  height: 4px;
  background: var(--color-bg-secondary);
  border-radius: 2px;
  margin-top: 16px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  background: var(--color-primary);
  border-radius: 2px;
  transition: width 0.6s ease;
}

.bar-fill.success {
  background: var(--color-success);
}

.bar-fill.warning {
  background: var(--color-warning);
}

.bar-fill.danger {
  background: var(--color-danger);
}

/* 点检统计面板 */
.inspection-panel {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  margin-bottom: 24px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.panel-title svg {
  width: 20px;
  height: 20px;
  color: var(--color-primary);
}

.panel-link {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--color-text-tertiary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

.panel-link:hover {
  color: var(--color-primary);
}

.panel-link svg {
  width: 16px;
  height: 16px;
}

.inspection-stats {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.istat-item {
  flex: 1;
  text-align: center;
}

.istat-value {
  font-family: var(--font-numbers);
  font-size: 28px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.istat-value.pending {
  color: var(--color-text-tertiary);
}

.istat-value.progress {
  color: var(--color-warning);
}

.istat-value.completed {
  color: var(--color-success);
}

.istat-value.today {
  color: var(--color-primary);
}

.istat-item.highlight .istat-value {
  color: var(--color-success);
  text-shadow: 0 0 20px var(--color-success-glow);
}

.istat-label {
  font-size: 12px;
  color: var(--color-text-tertiary);
  margin-top: 4px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.istat-divider {
  width: 1px;
  height: 40px;
  background: var(--color-divider);
  margin: 0 8px;
}

.progress-track {
  height: 6px;
  background: var(--color-bg-secondary);
  border-radius: 3px;
  margin-top: 20px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary), var(--color-success));
  border-radius: 3px;
  transition: width 0.6s ease;
}

/* 下方网格 */
.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 20px;
}

.dashboard-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 24px;
  transition: all var(--transition-base);
}

.dashboard-card:hover {
  border-color: var(--color-border-dim);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 20px;
}

.card-title svg {
  width: 20px;
  height: 20px;
  color: var(--color-primary);
}

/* 快速操作 */
.action-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.action-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--color-bg-secondary);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: all var(--transition-fast);
}

.action-item:hover {
  background: var(--color-bg-tertiary);
  border-color: var(--color-primary-dim);
  transform: translateX(4px);
}

.action-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-primary-dim);
  border-radius: var(--radius-sm);
  color: var(--color-primary);
}

.action-icon.primary {
  background: rgba(102, 126, 234, 0.15);
  color: #667eea;
}

.action-icon.success {
  background: var(--color-success-dim);
  color: var(--color-success);
}

.action-icon.warning {
  background: var(--color-warning-dim);
  color: var(--color-warning);
}

.action-icon.info {
  background: var(--color-info-dim);
  color: var(--color-info);
}

.action-icon svg {
  width: 20px;
  height: 20px;
}

.action-content {
  flex: 1;
}

.action-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.action-desc {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.action-arrow {
  width: 18px;
  height: 18px;
  color: var(--color-text-tertiary);
  transition: transform var(--transition-fast);
}

.action-item:hover .action-arrow {
  transform: translateX(4px);
  color: var(--color-primary);
}

/* 我的任务 */
.tasks-loading {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-bg-secondary);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.tasks-summary {
  display: flex;
  justify-content: space-around;
  margin-bottom: 16px;
}

.task-metric {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.metric-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  color: var(--color-text-tertiary);
}

.metric-icon.active {
  background: var(--color-warning-dim);
  color: var(--color-warning);
  animation: pulse 2s ease-in-out infinite;
}

.metric-icon.success {
  background: var(--color-success-dim);
  color: var(--color-success);
}

@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(255, 184, 0, 0.4); }
  50% { box-shadow: 0 0 0 8px rgba(255, 184, 0, 0); }
}

.metric-icon svg {
  width: 22px;
  height: 22px;
}

.metric-value {
  font-family: var(--font-numbers);
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-primary);
}

.metric-label {
  font-size: 12px;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.tasks-action {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  background: var(--color-primary-dim);
  border: 1px solid var(--color-primary-dim);
  border-radius: var(--radius-md);
  color: var(--color-primary);
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: all var(--transition-fast);
}

.tasks-action:hover {
  background: var(--color-primary);
  color: var(--color-bg-primary);
}

.tasks-action svg {
  width: 16px;
  height: 16px;
}

/* 待办事项 */
.todo-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 0;
  color: var(--color-text-tertiary);
}

.todo-empty svg {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
  opacity: 0.5;
}

.todo-empty p {
  font-size: 14px;
}

.todo-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.todo-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: all var(--transition-fast);
}

.todo-item:hover {
  background: var(--color-bg-tertiary);
  transform: translateX(4px);
}

.todo-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-warning);
  box-shadow: 0 0 8px var(--color-warning);
}

.todo-content {
  flex: 1;
}

.todo-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.todo-time {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.todo-tag {
  padding: 4px 10px;
  background: var(--color-warning-dim);
  color: var(--color-warning);
  font-size: 11px;
  font-weight: 500;
  border-radius: var(--radius-sm);
}

/* 响应式 */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .dashboard-grid {
    grid-template-columns: 1fr 1fr;
  }

  .todo-list {
    grid-column: span 2;
  }
}

@media (max-width: 768px) {
  .welcome-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .welcome-time {
    text-align: left;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .todo-list {
    grid-column: span 1;
  }

  .inspection-stats {
    flex-wrap: wrap;
    gap: 16px;
  }

  .istat-divider {
    display: none;
  }

  .istat-item {
    min-width: calc(50% - 8px);
  }
}
</style>
