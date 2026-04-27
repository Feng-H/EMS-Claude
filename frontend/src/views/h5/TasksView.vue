<template>
  <div class="h5-tasks-view">
    <mobile-header title="我的待办" :show-back="false" />

    <van-tabs v-model:active="activeTab" class="tasks-tabs" sticky offset-top="46">
      <van-tab title="点检" name="inspection">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <div class="task-list-container">
            <template v-if="loading && !refreshing">
              <van-skeleton v-for="i in 3" :key="i" title avatar :row="3" class="skeleton-item" />
            </template>
            <van-list
              v-else
              v-model:loading="loading"
              :finished="finished"
              finished-text="没有更多任务了"
              @load="onLoad"
            >
              <van-empty v-if="inspectionTasks.length === 0" description="今日暂无点检任务" />
              <div
                v-for="task in inspectionTasks"
                :key="task.id"
                class="task-card"
                @click="navigateToTask(task)"
              >
                <div class="task-header">
                  <span class="task-equipment">{{ task.equipment_name }}</span>
                  <van-tag :type="getTaskStatusType(task.status)" plain>
                    {{ getTaskStatusText(task.status) }}
                  </van-tag>
                </div>
                <div class="task-info">
                  <div class="task-detail">
                    <van-icon name="notes-o" />
                    {{ task.template_name }}
                  </div>
                  <div class="task-detail">
                    <van-icon name="clock-o" />
                    {{ task.scheduled_date }}
                  </div>
                </div>
              </div>
            </van-list>
          </div>
        </van-pull-refresh>
      </van-tab>

      <van-tab title="保养" name="maintenance">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <div class="task-list-container">
            <template v-if="loading && !refreshing">
              <van-skeleton v-for="i in 3" :key="i" title avatar :row="3" class="skeleton-item" />
            </template>
            <van-list
              v-else
              v-model:loading="loading"
              :finished="finished"
              finished-text="没有更多任务了"
              @load="onLoad"
            >
              <van-empty v-if="maintenanceTasks.length === 0" description="暂无待办保养任务" />
              <div
                v-for="task in maintenanceTasks"
                :key="task.id"
                class="task-card"
                @click="navigateToTask(task)"
              >
                <div class="task-header">
                  <span class="task-equipment">{{ task.equipment_name }}</span>
                  <van-tag :type="getTaskStatusType(task.status)" plain>
                    {{ getTaskStatusText(task.status) }}
                  </van-tag>
                </div>
                <div class="task-info">
                  <div class="task-detail">
                    <van-icon name="description" />
                    {{ task.plan_name }}
                  </div>
                  <div class="task-detail">
                    <van-icon name="clock-o" />
                    到期：{{ task.due_date }}
                  </div>
                </div>
              </div>
            </van-list>
          </div>
        </van-pull-refresh>
      </van-tab>

      <van-tab title="维修" name="repair">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <div class="task-list-container">
            <template v-if="loading && !refreshing">
              <van-skeleton v-for="i in 3" :key="i" title avatar :row="3" class="skeleton-item" />
            </template>
            <van-list
              v-else
              v-model:loading="loading"
              :finished="finished"
              finished-text="没有更多任务了"
              @load="onLoad"
            >
              <van-empty v-if="repairOrders.length === 0" description="暂无待办维修单" />
              <div
                v-for="order in repairOrders"
                :key="order.id"
                class="task-card"
                @click="navigateToTask(order)"
              >
                <div class="task-header">
                  <span class="task-equipment">{{ order.equipment_name }}</span>
                  <div class="header-right">
                    <van-tag :type="getPriorityType(order.priority)" style="margin-right: 4px">
                      {{ getPriorityText(order.priority) }}
                    </van-tag>
                    <van-tag :type="getOrderStatusType(order.status)" plain>
                      {{ getOrderStatusText(order.status) }}
                    </van-tag>
                  </div>
                </div>
                <div class="task-info">
                  <div class="task-detail fault-desc">
                    <van-icon name="warn-o" />
                    {{ order.fault_description }}
                  </div>
                  <div class="task-detail">
                    <van-icon name="clock-o" />
                    {{ formatTime(order.created_at) }}
                  </div>
                </div>
              </div>
            </van-list>
          </div>
        </van-pull-refresh>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { inspectionTaskApi, type InspectionTask } from '@/api/inspection'
import { maintenanceApi, type MaintenanceTask } from '@/api/maintenance'
import { repairOrderApi, type RepairOrder } from '@/api/repair'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const router = useRouter()
const activeTab = ref('inspection')
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)

const inspectionTasks = ref<InspectionTask[]>([])
const maintenanceTasks = ref<MaintenanceTask[]>([])
const repairOrders = ref<RepairOrder[]>([])

// 监听 Tab 切换，重置加载状态
watch(activeTab, () => {
  onRefresh()
})

const getTaskStatusType = (status: string) => {
  const types: Record<string, any> = {
    pending: 'primary',
    in_progress: 'warning',
    completed: 'success',
    cancelled: 'default'
  }
  return types[status] || 'default'
}

const getTaskStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待执行',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const getOrderStatusType = (status: string) => {
  const types: Record<string, any> = {
    pending: 'primary',
    assigned: 'primary',
    in_progress: 'warning',
    testing: 'success',
    completed: 'success',
    rejected: 'danger'
  }
  return types[status] || 'default'
}

const getOrderStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待处理',
    assigned: '已指派',
    in_progress: '维修中',
    testing: '待确认',
    completed: '已完成',
    rejected: '已驳回'
  }
  return texts[status] || status
}

const getPriorityText = (priority: number) => {
  const texts: Record<number, string> = {
    1: '高',
    2: '中',
    3: '低'
  }
  return texts[priority] || '未知'
}

const getPriorityType = (priority: number) => {
  const types: Record<number, any> = {
    1: 'danger',
    2: 'warning',
    3: 'primary'
  }
  return types[priority] || 'default'
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
}

const navigateToTask = (task: any) => {
  if (task.template_name) {
    router.push(`/h5/inspection?taskId=${task.id}`)
  } else if (task.plan_name) {
    router.push(`/h5/maintenance?taskId=${task.id}`)
  } else if (task.fault_description) {
    router.push(`/h5/repair/execute?orderId=${task.id}`)
  }
}

const loadTasks = async () => {
  try {
    if (activeTab.value === 'inspection') {
      const res = await inspectionTaskApi.getMyTasks()
      inspectionTasks.value = res.data
    } else if (activeTab.value === 'maintenance') {
      const res = await maintenanceApi.getMyTasks()
      maintenanceTasks.value = res.data
    } else if (activeTab.value === 'repair') {
      const res = await repairOrderApi.getMyTasks()
      repairOrders.value = res.data
    }
    // 目前后端接口一次性返回全部今日/我的任务，所以直接标记完成
    finished.value = true
  } catch (error: any) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const onRefresh = () => {
  finished.value = false
  loading.value = true
  loadTasks()
}

const onLoad = () => {
  if (refreshing.value) return
  loadTasks()
}
</script>

<style scoped>
.h5-tasks-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 60px;
}

.tasks-tabs :deep(.van-tabs__wrap) {
  height: 50px;
  background: rgba(250, 249, 245, 0.95);
  backdrop-filter: blur(10px);
}

.task-list-container {
  padding: var(--space-md);
  min-height: calc(100vh - 160px);
}

.skeleton-item {
  background: var(--color-bg-card);
  padding: 20px;
  border-radius: var(--radius-very);
  margin-bottom: var(--space-md);
}

.task-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  padding: 16px;
  margin-bottom: var(--space-sm);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  transition: all 0.2s;
}

.task-card:active {
  background: var(--color-bg-tertiary);
  transform: scale(0.98);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.task-equipment {
  font-weight: 600;
  font-size: 16px;
  color: var(--color-text-primary);
  font-family: var(--font-serif);
  flex: 1;
  padding-right: 8px;
}

.header-right {
  display: flex;
  align-items: center;
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-detail {
  font-size: 13px;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.task-detail :deep(.van-icon) {
  font-size: 14px;
  color: var(--color-terracotta);
}

.fault-desc {
  color: var(--color-text-primary);
  font-weight: 500;
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  .tasks-tabs :deep(.van-tabs__wrap) {
    background: rgba(26, 26, 19, 0.95);
  }
  .task-card {
    background: var(--color-bg-card);
  }
}
</style>
</style>
