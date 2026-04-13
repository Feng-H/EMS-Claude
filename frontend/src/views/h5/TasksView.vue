<template>
  <div class="h5-tasks-view">
    <mobile-header title="我的任务" :show-back="true" />

    <van-tabs v-model:active="activeTab" class="tasks-tabs">
      <van-tab title="点检任务" name="inspection">
        <van-empty v-if="inspectionTasks.length === 0" description="暂无点检任务" />
        <div v-else class="task-list">
          <div
            v-for="task in inspectionTasks"
            :key="task.id"
            class="task-card"
            @click="navigateToTask(task)"
          >
            <div class="task-header">
              <span class="task-equipment">{{ task.equipment_name }}</span>
              <van-tag :type="getTaskStatusType(task.status)">
                {{ getTaskStatusText(task.status) }}
              </van-tag>
            </div>
            <div class="task-info">
              <div class="task-detail">模板：{{ task.template_name }}</div>
              <div class="task-detail">计划时间：{{ task.scheduled_date }}</div>
            </div>
          </div>
        </div>
      </van-tab>

      <van-tab title="保养任务" name="maintenance">
        <van-empty v-if="maintenanceTasks.length === 0" description="暂无保养任务" />
        <div v-else class="task-list">
          <div
            v-for="task in maintenanceTasks"
            :key="task.id"
            class="task-card"
            @click="navigateToTask(task)"
          >
            <div class="task-header">
              <span class="task-equipment">{{ task.equipment_name }}</span>
              <van-tag :type="getTaskStatusType(task.status)">
                {{ getTaskStatusText(task.status) }}
              </van-tag>
            </div>
            <div class="task-info">
              <div class="task-detail">计划：{{ task.plan_name }}</div>
              <div class="task-detail">到期时间：{{ task.due_date }}</div>
              <div class="task-detail">工时：{{ task.work_hours }}小时</div>
            </div>
          </div>
        </div>
      </van-tab>

      <van-tab title="维修工单" name="repair">
        <van-empty v-if="repairOrders.length === 0" description="暂无维修工单" />
        <div v-else class="task-list">
          <div
            v-for="order in repairOrders"
            :key="order.id"
            class="task-card"
            @click="navigateToTask(order)"
          >
            <div class="task-header">
              <span class="task-equipment">{{ order.equipment_name }}</span>
              <van-tag :type="getOrderStatusType(order.status)">
                {{ getOrderStatusText(order.status) }}
              </van-tag>
            </div>
            <div class="task-info">
              <div class="task-detail">故障：{{ order.fault_description }}</div>
              <div class="task-detail">优先级：{{ getPriorityText(order.priority) }}</div>
              <div class="task-detail">创建时间：{{ formatTime(order.created_at) }}</div>
            </div>
          </div>
        </div>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { inspectionTaskApi, type InspectionTask } from '@/api/inspection'
import { maintenanceApi, type MaintenanceTask } from '@/api/maintenance'
import { repairOrderApi, type RepairOrder } from '@/api/repair'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const router = useRouter()
const activeTab = ref('inspection')

const inspectionTasks = ref<InspectionTask[]>([])
const maintenanceTasks = ref<MaintenanceTask[]>([])
const repairOrders = ref<RepairOrder[]>([])

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

const formatTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const navigateToTask = (task: any) => {
  if (task.template_name) {
    // 点检任务
    router.push(`/h5/inspection?taskId=${task.id}`)
  } else if (task.plan_name) {
    // 保养任务
    router.push(`/h5/maintenance?taskId=${task.id}`)
  } else if (task.fault_description) {
    // 维修工单
    router.push(`/h5/repair/execute?orderId=${task.id}`)
  }
}

const loadTasks = async () => {
  try {
    // 加载点检任务
    const inspections = await inspectionTaskApi.getMyTasks()
    inspectionTasks.value = inspections

    // 加载保养任务
    const maintenance = await maintenanceApi.getMyTasks()
    maintenanceTasks.value = maintenance

    // 加载维修工单
    const repairs = await repairOrderApi.getMyOrders()
    repairOrders.value = repairs
  } catch (error: any) {
    showToast('加载任务失败')
  }
}

onMounted(() => {
  loadTasks()
})
</script>

<style scoped>
.h5-tasks-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
}

/* 优化 Tabs 样式 */
.tasks-tabs :deep(.van-tabs__wrap) {
  background: rgba(250, 249, 245, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.tasks-tabs :deep(.van-tab) {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.tasks-tabs :deep(.van-tab--active) {
  color: var(--color-terracotta);
  font-weight: 500;
}

.tasks-tabs :deep(.van-tabs__line) {
  background: var(--color-terracotta);
  height: 3px;
  border-radius: 2px;
}

.task-list {
  padding: var(--space-lg);
}

.task-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  padding: var(--space-lg);
  margin-bottom: var(--space-md);
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-fast);
  border: 1px solid var(--color-border);
}

.task-card:active {
  transform: scale(0.98);
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.06);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-md);
}

.task-equipment {
  font-weight: 500;
  font-size: 16px;
  color: var(--color-text-primary);
  font-family: var(--font-serif);
}

.task-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-left: 14px;
  position: relative;
}

.task-info::before {
  content: '';
  position: absolute;
  left: 0;
  top: 6px;
  bottom: 6px;
  width: 3px;
  background: var(--color-terracotta);
  border-radius: 2px;
}

.task-detail {
  font-size: 13px;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.task-detail::before {
  content: '•';
  color: var(--color-text-muted);
}

/* 优化标签样式 */
.task-card :deep(.van-tag) {
  padding: 4px 10px;
  border-radius: var(--radius-subtle);
  font-size: 12px;
  font-weight: 500;
}

.task-card :deep(.van-tag--primary) {
  background: var(--color-primary-dim);
  color: var(--color-terracotta);
}

.task-card :deep(.van-tag--success) {
  background: var(--color-success-dim);
  color: var(--color-success);
}

.task-card :deep(.van-tag--warning) {
  background: var(--color-warning-dim);
  color: var(--color-warning);
}

.task-card :deep(.van-tag--danger) {
  background: var(--color-danger-dim);
  color: var(--color-danger);
}

.task-card :deep(.van-tag--default) {
  background: var(--color-bg-tertiary);
  color: var(--color-text-tertiary);
}

/* 空状态优化 */
.task-list :deep(.van-empty) {
  padding: 60px 20px;
}

.task-list :deep(.van-empty__description) {
  color: var(--color-text-tertiary);
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  .h5-tasks-view {
    background: var(--color-bg-primary);
  }

  .tasks-tabs :deep(.van-tabs__wrap) {
    background: rgba(26, 26, 19, 0.95);
  }

  .tasks-tabs :deep(.van-tab) {
    color: var(--color-text-tertiary);
  }

  .task-card {
    background: var(--color-bg-card);
    border-color: var(--color-border);
  }

  .task-equipment {
    color: var(--color-text-primary);
  }

  .task-detail {
    color: var(--color-text-tertiary);
  }

  .task-detail::before {
    color: var(--color-text-muted);
  }
}
</style>
