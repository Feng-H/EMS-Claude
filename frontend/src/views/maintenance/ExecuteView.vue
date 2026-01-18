<template>
  <div class="maintenance-execute-view">
    <!-- Header -->
    <van-nav-bar
      title="保养执行"
      left-text="返回"
      left-arrow
      @click-left="goBack"
    />

    <div v-if="!task" class="loading-container">
      <van-loading size="24">加载中...</van-loading>
    </div>

    <div v-else class="content">
      <!-- Task Info -->
      <van-cell-group inset class="info-card">
        <van-cell title="设备编号" :value="task.equipment_code" />
        <van-cell title="设备名称" :value="task.equipment_name" />
        <van-cell title="保养计划" :value="task.plan_name" />
        <van-cell title="到期日期" :value="task.due_date" />
        <van-cell title="状态">
          <template #value>
            <van-tag :type="getStatusTagType(task.status)">
              {{ getStatusName(task.status) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Maintenance Items -->
      <van-cell-group inset class="items-card">
        <van-cell title="保养项目" :value="`${completedCount}/${itemCount}`" />
      </van-cell-group>

      <div class="items-list">
        <div
          v-for="(item, index) in maintenanceItems"
          :key="item.id"
          class="item-card"
          :class="{ completed: records[item.id]?.result }"
        >
          <div class="item-header">
            <span class="item-index">{{ index + 1 }}</span>
            <span class="item-name">{{ item.name }}</span>
            <van-tag
              v-if="records[item.id]?.result"
              :type="records[item.id].result === 'OK' ? 'success' : 'danger'"
              size="small"
            >
              {{ records[item.id].result }}
            </van-tag>
          </div>

          <div v-if="item.method" class="item-method">
            <span class="label">方法：</span>
            <span>{{ item.method }}</span>
          </div>

          <div v-if="item.criteria" class="item-criteria">
            <span class="label">标准：</span>
            <span>{{ item.criteria }}</span>
          </div>

          <!-- Result Input -->
          <div class="item-actions">
            <van-button
              v-if="!records[item.id]?.result"
              type="primary"
              size="small"
              @click="showItemDialog(item)"
            >
              填写结果
            </van-button>
            <van-button
              v-else
              size="small"
              @click="showItemDialog(item)"
            >
              修改
            </van-button>
          </div>
        </div>
      </div>

      <!-- Complete Button -->
      <div class="action-bar">
        <van-button
          type="primary"
          block
          :disabled="!canComplete"
          :loading="completing"
          @click="handleComplete"
        >
          完成保养
        </van-button>
      </div>

      <!-- Remark -->
      <van-cell-group inset class="remark-card">
        <van-field
          v-model="remark"
          label="备注"
          type="textarea"
          placeholder="请输入备注信息"
          rows="2"
        />
        <van-field
          v-model.number="actualHours"
          label="实际工时"
          type="number"
          placeholder="请输入实际工时"
          input-align="right"
        />
      </van-cell-group>
    </div>

    <!-- Item Result Dialog -->
    <van-dialog
      v-model:show="showItemResultDialog"
      :title="`填写结果 - ${currentItem?.name}`"
      show-cancel-button
      confirm-button-text="确定"
      @confirm="confirmItemResult"
    >
      <van-radio-group v-model="itemResultForm.result" class="result-options">
        <van-radio name="OK">
          <span class="result-option ok">合格 (OK)</span>
        </van-radio>
        <van-radio name="NG">
          <span class="result-option ng">不合格 (NG)</span>
        </van-radio>
      </van-radio-group>

      <van-field
        v-model="itemResultForm.remark"
        label="备注"
        type="textarea"
        placeholder="请输入备注"
        rows="2"
        border
      />

      <van-field
        v-model="itemResultForm.photo_url"
        label="照片"
        placeholder="照片URL"
        border
      />
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import {
  getMaintenanceTask,
  startMaintenance,
  completeMaintenance,
  getStatusName,
  type MaintenanceTask
} from '@/api/maintenance'
import { equipmentApi } from '@/api/equipment'
import type { MaintenanceItemRecord } from '@/api/maintenance'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const task = ref<MaintenanceTask | null>(null)
const maintenanceItems = ref<any[]>([])
const records = ref<Record<number, MaintenanceItemRecord>>({})
const remark = ref('')
const actualHours = ref(0)
const completing = ref(false)

const showItemResultDialog = ref(false)
const currentItem = ref<any>(null)

const itemResultForm = reactive({
  result: 'OK',
  remark: '',
  photo_url: ''
})

const completedCount = computed(() => {
  return Object.keys(records.value).length
})

const itemCount = computed(() => {
  return maintenanceItems.value.length || task.value?.item_count || 0
})

const canComplete = computed(() => {
  return completedCount.value >= itemCount.value
})

const getStatusTagType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'warning',
    in_progress: 'primary',
    completed: 'success',
    overdue: 'danger'
  }
  return types[status] || 'default'
}

const goBack = () => {
  router.back()
}

const loadTask = async () => {
  const taskId = route.query.id as string
  if (!taskId) {
    // If no ID is provided (e.g., direct access from sidebar), don't crash or go back.
    // Just show empty state or redirect to task list if desired.
    // Here we just return and show the "No Data" state if handled by template,
    // but the template expects 'task' to be populated for 'v-else'
    // So let's redirect to Task List with a message
    showToast('请从任务列表选择要执行的任务')
    router.replace('/maintenance/tasks')
    return
  }

  try {
    const res = await getMaintenanceTask(Number(taskId))
    task.value = res.data

    // Auto-start if pending
    if (task.value.status === 'pending') {
      await startMaintenance({
        task_id: Number(taskId)
      })
      task.value.status = 'in_progress'
    }

    // Load maintenance items (mock data for now)
    maintenanceItems.value = []
    for (let i = 1; i <= task.value.item_count; i++) {
      maintenanceItems.value.push({
        id: i,
        plan_id: task.value.plan_id,
        name: `保养项目 ${i}`,
        method: '按标准操作',
        criteria: '符合技术要求',
        sequence_order: i
      })
    }
  } catch (err: any) {
    showToast(err.response?.data?.error || '加载任务失败')
    router.replace('/maintenance/tasks')
  }
}

const showItemDialog = (item: any) => {
  currentItem.value = item
  const existing = records.value[item.id]
  if (existing) {
    itemResultForm.result = existing.result
    itemResultForm.remark = existing.remark || ''
    itemResultForm.photo_url = existing.photo_url || ''
  } else {
    itemResultForm.result = 'OK'
    itemResultForm.remark = ''
    itemResultForm.photo_url = ''
  }
  showItemResultDialog.value = true
}

const confirmItemResult = () => {
  if (!currentItem.value) return

  records.value[currentItem.value.id] = {
    item_id: currentItem.value.id,
    result: itemResultForm.result,
    remark: itemResultForm.remark,
    photo_url: itemResultForm.photo_url
  }
}

const handleComplete = async () => {
  if (!canComplete.value) {
    showToast('请先完成所有保养项目')
    return
  }

  try {
    await showConfirmDialog({
      title: '确认完成',
      message: '确认完成所有保养项目？'
    })
  } catch {
    return
  }

  completing.value = true
  try {
    const recordsArray = Object.values(records.value)
    await completeMaintenance({
      task_id: task.value!.id,
      records: recordsArray,
      actual_hours: actualHours.value || undefined,
      remark: remark.value || undefined
    })

    showToast('保养完成')
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (err: any) {
    showToast(err.response?.data?.error || '完成失败')
  } finally {
    completing.value = false
  }
}

onMounted(() => {
  loadTask()
})
</script>

<style scoped>
.maintenance-execute-view {
  min-height: 100vh;
  background: #f5f5f5;
  padding-top: 46px; /* NavBar 高度 */
  padding-bottom: 80px;
}

/* Fix Vant nav bar arrow size */
:deep(.van-nav-bar__arrow) {
  font-size: 18px !important;
}

:deep(.van-icon__image) {
  width: 18px;
  height: 18px;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.content {
  padding: 16px;
}

.info-card,
.items-card,
.remark-card {
  margin-bottom: 16px;
}

.items-list {
  margin-bottom: 16px;
}

.item-card {
  background: white;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
}

.item-card.completed {
  background: #f0f9ff;
}

.item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.item-index {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.item-name {
  flex: 1;
  font-weight: 500;
}

.item-method,
.item-criteria {
  font-size: 13px;
  color: #606266;
  margin-bottom: 4px;
  padding-left: 32px;
}

.label {
  color: #909399;
}

.item-actions {
  margin-top: 8px;
  padding-left: 32px;
}

.result-options {
  padding: 16px;
}

.result-option {
  margin-left: 8px;
}

.result-option.ok {
  color: #67c23a;
}

.result-option.ng {
  color: #f56c6c;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}
</style>
