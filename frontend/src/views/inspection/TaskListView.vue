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
          <el-select v-model="filterForm.status" placeholder="全部" clearable @change="loadTasks">
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
          <el-select v-model="filterForm.assigned_to" placeholder="全部" clearable filterable>
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
    name: 'inspection-execute',
    params: { taskId: task.id }
  })
}

onMounted(() => {
  loadTasks()
  loadStatistics()
})
</script>

<style scoped>
.task-list-view {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-form {
  margin-bottom: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-pending {
  color: #909399;
}

.stat-progress {
  color: #e6a23c;
}

.stat-completed {
  color: #67c23a;
}

.task-table {
  margin-bottom: 16px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
}

.task-detail :deep(.el-descriptions) {
  margin-bottom: 16px;
}
</style>
