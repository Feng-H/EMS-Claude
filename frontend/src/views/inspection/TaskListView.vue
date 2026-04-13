<template>
  <div class="task-list-view">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>点检任务管理</span>
          <div class="header-actions">
            <el-button @click="loadTasks" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable @change="loadTasks" style="width: 160px">
            <el-option label="待执行" value="pending" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已逾期" value="overdue" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="onDateChange"
          />
        </el-form-item>
        <el-form-item label="指派给">
          <el-select v-model="filterForm.assigned_to" placeholder="全部" clearable filterable style="width: 160px">
            <el-option
              v-for="user in users"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="6">
          <el-statistic title="总任务" :value="statistics.total_tasks" />
        </el-col>
        <el-col :span="6">
          <el-statistic title="待执行" :value="statistics.pending_tasks">
            <template #suffix>
              <span class="stat-pending">件</span>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-statistic title="进行中" :value="statistics.in_progress_tasks">
            <template #suffix>
              <span class="stat-progress">件</span>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="6">
          <el-statistic title="今日完成" :value="statistics.today_completed">
            <template #suffix>
              <span class="stat-completed">件</span>
            </template>
          </el-statistic>
        </el-col>
      </el-row>

      <!-- 任务表格 -->
      <el-table :data="tasks" v-loading="loading" stripe class="task-table">
        <el-table-column prop="equipment_code" label="设备编号" width="120" />
        <el-table-column prop="equipment_name" label="设备名称" width="150" />
        <el-table-column prop="template_name" label="点检模板" width="120" />
        <el-table-column prop="assignee_name" label="执行人" width="100" />
        <el-table-column prop="scheduled_date" label="计划日期" width="110" />
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="120" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="getProgress(row)"
              :color="getProgressColor(row)"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column prop="started_at" label="开始时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.started_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="completed_at" label="完成时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.completed_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              link
              @click="viewDetail(row)"
            >
              详情
            </el-button>
            <el-button
              v-if="row.status === 'pending' || row.status === 'in_progress'"
              type="success"
              size="small"
              link
              @click="executeTask(row)"
            >
              执行
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadTasks"
          @current-change="loadTasks"
        />
      </div>
    </el-card>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="点检任务详情"
      width="800px"
    >
      <div v-if="currentTask" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="设备编号">
            {{ currentTask.equipment_code }}
          </el-descriptions-item>
          <el-descriptions-item label="设备名称">
            {{ currentTask.equipment_name }}
          </el-descriptions-item>
          <el-descriptions-item label="点检模板">
            {{ currentTask.template_name }}
          </el-descriptions-item>
          <el-descriptions-item label="执行人">
            {{ currentTask.assignee_name }}
          </el-descriptions-item>
          <el-descriptions-item label="计划日期">
            {{ currentTask.scheduled_date }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentTask.status)">
              {{ getStatusText(currentTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="开始时间">
            {{ formatTime(currentTask.started_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="完成时间">
            {{ formatTime(currentTask.completed_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="检查进度" :span="2">
            <el-progress
              :percentage="getProgress(currentTask)"
              :color="getProgressColor(currentTask)"
            />
          </el-descriptions-item>
        </el-descriptions>

        <!-- 点检记录 -->
        <el-divider>点检记录</el-divider>
        <el-table :data="currentTaskRecords" size="small">
          <el-table-column prop="item_name" label="检查项目" />
          <el-table-column label="结果" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.result === 'OK' ? 'success' : 'danger'" size="small">
                {{ row.result }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" />
          <el-table-column label="照片" width="80" align="center">
            <template #default="{ row }">
              <el-button
                v-if="row.photo_url"
                type="primary"
                size="small"
                link
                @click="viewPhoto(row.photo_url)"
              >
                查看
              </el-button>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="记录时间" width="160">
            <template #default="{ row }">
              {{ formatDateTime(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
      <el-empty v-else description="暂无数据" />
    </el-dialog>

    <!-- 照片预览 -->
    <el-dialog v-model="showPhotoDialog" title="照片预览" width="600px">
      <img :src="currentPhoto" style="width: 100%" alt="点检照片" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import {
  inspectionTaskApi,
  type InspectionTask,
  type InspectionRecord,
  type InspectionStatistics
} from '@/api/inspection'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const tasks = ref<InspectionTask[]>([])
const currentTask = ref<InspectionTask | null>(null)
const currentTaskRecords = ref<InspectionRecord[]>([])
const statistics = ref<InspectionStatistics>({
  total_tasks: 0,
  pending_tasks: 0,
  in_progress_tasks: 0,
  completed_tasks: 0,
  overdue_tasks: 0,
  today_completed: 0,
  completion_rate: 0
})

const dateRange = ref<[string, string] | null>(null)
const showDetailDialog = ref(false)
const showPhotoDialog = ref(false)
const currentPhoto = ref('')

const filterForm = reactive({
  status: '',
  assigned_to: undefined as number | undefined,
  date_from: '',
  date_to: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 模拟用户列表（实际应从API获取）
const users = ref([
  { id: 1, name: '管理员' },
  { id: 2, name: '工程师' },
  { id: 3, name: '操作员' }
])

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    pending: 'info',
    in_progress: 'warning',
    completed: 'success',
    overdue: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待执行',
    in_progress: '进行中',
    completed: '已完成',
    overdue: '已逾期'
  }
  return map[status] || status
}

const getProgress = (task: InspectionTask) => {
  if (task.item_count === 0) return 0
  return Math.round((task.completed_count / task.item_count) * 100)
}

const getProgressColor = (task: InspectionTask) => {
  const progress = getProgress(task)
  if (progress === 100) return '#67c23a'
  if (progress > 50) return '#e6a23c'
  return '#f56c6c'
}

const formatTime = (time?: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDateTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

const onDateChange = (dates: [string, string] | null) => {
  if (dates) {
    filterForm.date_from = dates[0]
    filterForm.date_to = dates[1]
  } else {
    filterForm.date_from = ''
    filterForm.date_to = ''
  }
  loadTasks()
}

const loadTasks = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (filterForm.status) params.status = filterForm.status
    if (filterForm.assigned_to) params.assigned_to = filterForm.assigned_to
    if (filterForm.date_from) params.date_from = filterForm.date_from
    if (filterForm.date_to) params.date_to = filterForm.date_to

    const data = await inspectionTaskApi.getTasks(params)
    tasks.value = data.items
    pagination.total = data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载任务列表失败')
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    const data = await inspectionTaskApi.getStatistics()
    statistics.value = data
  } catch (error: any) {
    console.error('加载统计数据失败', error)
  }
}

const viewDetail = async (task: InspectionTask) => {
  currentTask.value = task
  currentTaskRecords.value = []
  showDetailDialog.value = true

  // 如果已完成，加载点检记录
  if (task.status === 'completed') {
    try {
      // TODO: 实现获取点检记录的API
      // const records = await inspectionTaskApi.getRecords(task.id)
      // currentTaskRecords.value = records
    } catch (error: any) {
      console.error('加载点检记录失败', error)
    }
  }
}

const viewPhoto = (url: string) => {
  currentPhoto.value = url
  showPhotoDialog.value = true
}

const executeTask = (task: InspectionTask) => {
  router.push({
    name: 'H5Inspection',
    query: { taskId: task.id }
  })
}

onMounted(() => {
  loadTasks()
  loadStatistics()
})
</script>

<style scoped>
.task-list-view {
  padding: 24px;
  max-width: 1600px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-header span::before {
  content: '';
  width: 4px;
  height: 18px;
  background: linear-gradient(180deg, var(--color-primary), #764ba2);
  border-radius: 2px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* 筛选表单优化 */
.filter-form {
  background: var(--color-bg-secondary);
  padding: 20px;
  border-radius: var(--radius-lg);
  margin-bottom: 20px;
  border: 1px solid var(--color-border);
}

.filter-form :deep(.el-form-item) {
  margin-bottom: 0;
  margin-right: 20px;
}

.filter-form :deep(.el-form-item__label) {
  color: var(--color-text-secondary);
  font-weight: 500;
  font-size: 14px;
}

/* 统计卡片行优化 */
.stats-row {
  margin-bottom: 20px;
}

.stats-row :deep(.el-col) {
  margin-bottom: 12px;
}

.stats-row :deep(.el-statistic) {
  background: var(--color-bg-card);
  padding: 20px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  transition: all var(--transition-base);
  position: relative;
  overflow: hidden;
}

.stats-row :deep(.el-statistic::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--color-primary), #764ba2);
  opacity: 0;
  transition: opacity var(--transition-base);
}

.stats-row :deep(.el-statistic:hover) {
  border-color: var(--color-primary-dim);
  transform: translateY(-4px);
  box-shadow: var(--shadow-card);
}

.stats-row :deep(.el-statistic:hover::before) {
  opacity: 1;
}

.stats-row :deep(.el-statistic__head) {
  font-size: 13px;
  color: var(--color-text-tertiary);
  margin-bottom: 8px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stats-row :deep(.el-statistic__number) {
  font-size: 32px;
  font-weight: 700;
  color: var(--color-text-primary);
  font-family: var(--font-numbers);
}

.stats-row :deep(.el-statistic__content) {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.stat-pending {
  color: var(--color-text-tertiary);
}

.stat-progress {
  color: var(--color-warning);
}

.stat-completed {
  color: var(--color-success);
}

/* 表格优化 */
.task-table {
  margin-bottom: 20px;
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.task-table :deep(.el-table__header-wrapper) {
  background: var(--color-bg-secondary);
}

.task-table :deep(.el-table__header th) {
  background: var(--color-bg-elevated);
  color: var(--color-text-secondary);
  font-weight: 600;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 2px solid var(--color-border);
  padding: 14px 0;
}

.task-table :deep(.el-table__body tr) {
  transition: background var(--transition-fast);
}

.task-table :deep(.el-table__body tr:hover > td) {
  background: var(--color-bg-tertiary);
}

.task-table :deep(.el-table__body td) {
  padding: 14px 0;
  border-bottom: 1px solid var(--color-divider);
}

/* 标签优化 */
.task-table :deep(.el-tag) {
  padding: 4px 10px;
  border-radius: var(--radius-sm);
  font-weight: 500;
  font-size: 12px;
}

/* 进度条优化 */
.task-table :deep(.el-progress) {
  width: 120px;
}

.task-table :deep(.el-progress__bar) {
  border-radius: var(--radius-sm);
}

/* 按钮优化 */
.task-table :deep(.el-button--small.is-link) {
  font-weight: 500;
  padding: 4px 8px;
}

/* 分页优化 */
.pagination {
  display: flex;
  justify-content: flex-end;
  padding: 16px 0;
  gap: 8px;
}

.pagination :deep(.el-pagination) {
  gap: 8px;
}

.pagination :deep(.el-pagination.is-background .el-pager li) {
  border-radius: var(--radius-sm);
  font-weight: 500;
  min-width: 32px;
  height: 32px;
  line-height: 32px;
}

.pagination :deep(.el-pagination.is-background .el-pager li:not(.disabled).is-active) {
  background: linear-gradient(135deg, var(--color-primary), #764ba2);
}

.pagination :deep(.el-pagination.is-background .btn-prev),
.pagination :deep(.el-pagination.is-background .btn-next) {
  border-radius: var(--radius-sm);
  min-width: 32px;
  height: 32px;
}

/* 详情对话框优化 */
.task-detail {
  padding: 8px 0;
}

.task-detail :deep(.el-descriptions) {
  margin-bottom: 20px;
}

.task-detail :deep(.el-descriptions__label) {
  font-weight: 500;
  width: 120px;
}

.task-detail :deep(.el-divider__text) {
  font-weight: 600;
  color: var(--color-text-primary);
}

/* 空状态优化 */
:deep(.el-empty) {
  padding: 60px 0;
}

:deep(.el-empty__description) {
  color: var(--color-text-tertiary);
}

/* 对话框优化 */
:deep(.el-dialog) {
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
}

:deep(.el-dialog__header) {
  border-bottom: 1px solid var(--color-border);
  padding: 20px;
}

:deep(.el-dialog__title) {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text-primary);
}

:deep(.el-dialog__body) {
  padding: 24px;
}

/* 卡片优化 */
:deep(.el-card) {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
}

:deep(.el-card:hover) {
  box-shadow: var(--shadow-md);
}

:deep(.el-card__header) {
  border-bottom: 1px solid var(--color-border);
  padding: 16px 20px;
  background: var(--color-bg-elevated);
}

:deep(.el-card__body) {
  padding: 20px;
}

/* 响应式 */
@media (max-width: 1200px) {
  .task-list-view {
    padding: 16px;
  }

  .stats-row :deep(.el-statistic__number) {
    font-size: 24px;
  }
}

@media (max-width: 768px) {
  .filter-form :deep(.el-form-item) {
    width: 100%;
    margin-right: 0;
    margin-bottom: 12px;
  }

  .stats-row :deep(.el-statistic__number) {
    font-size: 20px;
  }
}
</style>
