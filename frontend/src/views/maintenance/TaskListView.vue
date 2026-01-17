<template>
  <div class="maintenance-task-view">
    <div class="header">
      <h2>保养任务管理</h2>
    </div>

    <!-- Filters -->
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="待执行" value="pending" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已逾期" value="overdue" />
          </el-select>
        </el-form-item>
        <el-form-item label="指派人">
          <el-select v-model="filterForm.assigned_to" placeholder="全部" clearable filterable style="width: 150px">
            <el-option
              v-for="user in users"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="计划日期">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Statistics Cards -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="4">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-value">{{ stats.total_tasks }}</div>
          <div class="stat-label">总任务数</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card pending" shadow="hover">
          <div class="stat-value">{{ stats.pending_tasks }}</div>
          <div class="stat-label">待执行</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card progress" shadow="hover">
          <div class="stat-value">{{ stats.in_progress_tasks }}</div>
          <div class="stat-label">进行中</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card success" shadow="hover">
          <div class="stat-value">{{ stats.completed_tasks }}</div>
          <div class="stat-label">已完成</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card danger" shadow="hover">
          <div class="stat-value">{{ stats.overdue_tasks }}</div>
          <div class="stat-label">已逾期</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-value">{{ stats.completion_rate.toFixed(1) }}%</div>
          <div class="stat-label">完成率</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Tasks Table -->
    <el-card class="table-card">
      <el-table :data="tasks" v-loading="loading" stripe>
        <el-table-column prop="id" label="任务ID" width="80" />
        <el-table-column label="设备" min-width="180">
          <template #default="{ row }">
            <div>{{ row.equipment_code }}</div>
            <div class="text-secondary">{{ row.equipment_name }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="plan_name" label="保养计划" min-width="150" />
        <el-table-column label="保养人" width="100">
          <template #default="{ row }">{{ row.assignee_name || '-' }}</template>
        </el-table-column>
        <el-table-column label="计划日期" width="110">
          <template #default="{ row }">{{ row.scheduled_date }}</template>
        </el-table-column>
        <el-table-column label="到期日期" width="110">
          <template #default="{ row }">{{ row.due_date }}</template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="100">
          <template #default="{ row }">
            {{ row.completed_count }}/{{ row.item_count }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending' || row.status === 'overdue'"
              link
              type="primary"
              size="small"
              @click="handleStart(row)"
            >
              开始
            </el-button>
            <el-button
              v-if="row.status === 'in_progress'"
              link
              type="success"
              size="small"
              @click="handleComplete(row)"
            >
              完成
            </el-button>
            <el-button
              v-if="row.status === 'completed'"
              link
              type="info"
              size="small"
              @click="viewDetail(row)"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadTasks"
          @current-change="loadTasks"
        />
      </div>
    </el-card>

    <!-- Start Confirm Dialog -->
    <el-dialog v-model="showStartDialog" title="开始保养" width="400px">
      <el-form :model="startForm" label-width="80px">
        <el-form-item label="设备">
          <span>{{ currentTask?.equipment_code }} - {{ currentTask?.equipment_name }}</span>
        </el-form-item>
        <el-form-item label="保养计划">
          <span>{{ currentTask?.plan_name }}</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showStartDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmStart" :loading="starting">确认开始</el-button>
      </template>
    </el-dialog>

    <!-- Detail Dialog -->
    <el-dialog v-model="showDetailDialog" title="保养详情" width="600px">
      <el-descriptions :column="2" border v-if="currentTask">
        <el-descriptions-item label="任务ID">{{ currentTask.id }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentTask.status)">
            {{ getStatusName(currentTask.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="设备">{{ currentTask.equipment_code }}</el-descriptions-item>
        <el-descriptions-item label="设备名称">{{ currentTask.equipment_name }}</el-descriptions-item>
        <el-descriptions-item label="保养计划">{{ currentTask.plan_name }}</el-descriptions-item>
        <el-descriptions-item label="保养人">{{ currentTask.assignee_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="计划日期">{{ currentTask.scheduled_date }}</el-descriptions-item>
        <el-descriptions-item label="到期日期">{{ currentTask.due_date }}</el-descriptions-item>
        <el-descriptions-item label="开始时间" v-if="currentTask.started_at">
          {{ formatDateTime(currentTask.started_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="完成时间" v-if="currentTask.completed_at">
          {{ formatDateTime(currentTask.completed_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="实际工时" v-if="currentTask.actual_hours">
          {{ currentTask.actual_hours }}h
        </el-descriptions-item>
        <el-descriptions-item label="进度">
          {{ currentTask.completed_count }}/{{ currentTask.item_count }}
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2" v-if="currentTask.remark">
          {{ currentTask.remark }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="showDetailDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getMaintenanceTasks,
  startMaintenance,
  getMaintenanceStatistics,
  getStatusName,
  getStatusType,
  type MaintenanceTask,
  type MaintenanceTaskQuery,
  type MaintenanceStatistics
} from '@/api/maintenance'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const tasks = ref<MaintenanceTask[]>([])
const stats = ref<MaintenanceStatistics>({
  total_plans: 0,
  total_tasks: 0,
  pending_tasks: 0,
  in_progress_tasks: 0,
  completed_tasks: 0,
  overdue_tasks: 0,
  today_completed: 0,
  completion_rate: 0
})

const users = ref<any[]>([]) // In real app, fetch from user API

// Filter
const filterForm = reactive<MaintenanceTaskQuery>({
  status: '',
  assigned_to: undefined,
  date_from: '',
  date_to: '',
  page: 1,
  page_size: 20
})

const dateRange = ref<string[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// Dialogs
const showStartDialog = ref(false)
const showDetailDialog = ref(false)
const currentTask = ref<MaintenanceTask | null>(null)
const starting = ref(false)

const startForm = reactive({
  task_id: 0
})

// Watch date range
watch(dateRange, (val) => {
  if (val && val.length === 2) {
    filterForm.date_from = val[0]
    filterForm.date_to = val[1]
  } else {
    filterForm.date_from = ''
    filterForm.date_to = ''
  }
})

// Methods
const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadStats = async () => {
  try {
    const res = await getMaintenanceStatistics()
    stats.value = res.data || {
      total_plans: 0,
      total_tasks: 0,
      pending_tasks: 0,
      in_progress_tasks: 0,
      completed_tasks: 0,
      overdue_tasks: 0,
      today_completed: 0,
      completion_rate: 0
    }
  } catch (err) {
    console.error('Failed to load stats:', err)
  }
}

const loadTasks = async () => {
  loading.value = true
  try {
    filterForm.page = pagination.page
    filterForm.page_size = pagination.pageSize
    const res = await getMaintenanceTasks(filterForm)
    tasks.value = res.data.items
    pagination.total = res.data.total
  } catch (err) {
    ElMessage.error('加载任务列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadTasks()
}

const handleReset = () => {
  filterForm.status = ''
  filterForm.assigned_to = undefined
  dateRange.value = []
  pagination.page = 1
  loadTasks()
}

const handleStart = (task: MaintenanceTask) => {
  currentTask.value = task
  startForm.task_id = task.id
  showStartDialog.value = true
}

const confirmStart = async () => {
  starting.value = true
  try {
    await startMaintenance({
      task_id: startForm.task_id
    })
    ElMessage.success('保养开始')
    showStartDialog.value = false
    loadTasks()
    loadStats()
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '操作失败')
  } finally {
    starting.value = false
  }
}

const handleComplete = (task: MaintenanceTask) => {
  // Redirect to mobile execution page or show completion dialog
  ElMessage.info('请在移动端完成保养任务')
}

const viewDetail = (task: MaintenanceTask) => {
  currentTask.value = task
  showDetailDialog.value = true
}

onMounted(() => {
  loadStats()
  loadTasks()
})
</script>

<style scoped>
.maintenance-task-view {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.filter-card {
  margin-bottom: 20px;
}

.filter-form {
  margin-bottom: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  padding: 10px;
}

.stat-card.pending {
  border-left: 3px solid #e6a23c;
}

.stat-card.progress {
  border-left: 3px solid #409eff;
}

.stat-card.success {
  border-left: 3px solid #67c23a;
}

.stat-card.danger {
  border-left: 3px solid #f56c6c;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}

.table-card {
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.text-secondary {
  font-size: 12px;
  color: #909399;
}
</style>
