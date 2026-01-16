<template>
  <div class="order-list-view">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>维修工单管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="$router.push('/repair/create')">
              <el-icon><Plus /></el-icon>
              报修申请
            </el-button>
            <el-button @click="loadOrders" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 筛选条件 -->
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable @change="loadOrders">
            <el-option label="待派单" value="pending" />
            <el-option label="已派单" value="assigned" />
            <el-option label="维修中" value="in_progress" />
            <el-option label="待测试" value="testing" />
            <el-option label="待审核" value="confirmed" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="filterForm.priority" placeholder="全部" clearable @change="loadOrders">
            <el-option label="高" :value="1" />
            <el-option label="中" :value="2" />
            <el-option label="低" :value="3" />
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
      </el-form>

      <!-- 统计卡片 -->
      <el-row :gutter="16" class="stats-row">
        <el-col :span="4">
          <el-statistic title="总工单" :value="statistics.total_orders" />
        </el-col>
        <el-col :span="4">
          <el-statistic title="待派单" :value="statistics.pending_orders">
            <template #suffix>
              <span class="stat-pending">件</span>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="维修中" :value="statistics.in_progress_orders">
            <template #suffix>
              <span class="stat-progress">件</span>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="今日完成" :value="statistics.today_completed">
            <template #suffix>
              <span class="stat-completed">件</span>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="今日新增" :value="statistics.today_created">
            <template #suffix>
              <span class="stat-created">件</span>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="平均响应" :value="statistics.avg_response_time" :precision="0">
            <template #suffix>分钟</template>
          </el-statistic>
        </el-col>
      </el-row>

      <!-- 工单表格 -->
      <el-table :data="orders" v-loading="loading" stripe class="order-table">
        <el-table-column prop="id" label="工单号" width="80" />
        <el-table-column prop="equipment_code" label="设备编号" width="120" />
        <el-table-column prop="equipment_name" label="设备名称" width="150" />
        <el-table-column prop="fault_description" label="故障描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="priority" label="优先级" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="assignee_name" label="维修工" width="100" />
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="报修时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="viewDetail(row)">
              详情
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="warning"
              size="small"
              link
              @click="showAssignDialog(row)"
            >
              派单
            </el-button>
            <el-button
              v-if="row.status === 'assigned' || row.status === 'in_progress'"
              type="success"
              size="small"
              link
              @click="executeOrder(row)"
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
          @size-change="loadOrders"
          @current-change="loadOrders"
        />
      </div>
    </el-card>

    <!-- 工单详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="工单详情"
      width="700px"
    >
      <div v-if="currentOrder" class="order-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="工单号">
            #{{ currentOrder.id }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentOrder.status)">
              {{ getStatusText(currentOrder.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="设备编号">
            {{ currentOrder.equipment_code }}
          </el-descriptions-item>
          <el-descriptions-item label="设备名称">
            {{ currentOrder.equipment_name }}
          </el-descriptions-item>
          <el-descriptions-item label="故障代码" v-if="currentOrder.fault_code">
            {{ currentOrder.fault_code }}
          </el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentOrder.priority)">
              {{ getPriorityText(currentOrder.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="报修人">
            {{ currentOrder.reporter_name }}
          </el-descriptions-item>
          <el-descriptions-item label="维修工">
            {{ currentOrder.assignee_name || '未分配' }}
          </el-descriptions-item>
          <el-descriptions-item label="报修时间" :span="2">
            {{ formatDateTime(currentOrder.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="故障描述" :span="2">
            {{ currentOrder.fault_description }}
          </el-descriptions-item>
          <el-descriptions-item label="解决方案" v-if="currentOrder.solution" :span="2">
            {{ currentOrder.solution }}
          </el-descriptions-item>
          <el-descriptions-item label="备件使用" v-if="currentOrder.spare_parts" :span="2">
            {{ currentOrder.spare_parts }}
          </el-descriptions-item>
          <el-descriptions-item label="实际工时" v-if="currentOrder.actual_hours">
            {{ currentOrder.actual_hours }} 小时
          </el-descriptions-item>
        </el-descriptions>

        <!-- 照片 -->
        <div v-if="currentOrder.photos && currentOrder.photos.length > 0" class="photos-section">
          <div class="section-label">故障照片</div>
          <div class="photos-grid">
            <el-image
              v-for="(photo, index) in currentOrder.photos"
              :key="index"
              :src="photo"
              fit="cover"
              class="photo-item"
              :preview-src-list="currentOrder.photos"
            />
          </div>
        </div>

        <!-- 维修日志 -->
        <el-divider>处理记录</el-divider>
        <el-timeline v-if="currentOrder.logs && currentOrder.logs.length > 0">
          <el-timeline-item
            v-for="log in currentOrder.logs"
            :key="log.id"
            :timestamp="formatDateTime(log.created_at)"
          >
            <strong>{{ log.user_name }}</strong> {{ log.content }}
          </el-timeline-item>
        </el-timeline>
        <el-empty v-else description="暂无处理记录" :image-size="60" />
      </div>
    </el-dialog>

    <!-- 派单对话框 -->
    <el-dialog v-model="showAssignDialog" title="派单" width="500px">
      <el-form :model="assignForm" :rules="assignRules" ref="assignFormRef" label-width="80px">
        <el-form-item label="维修工" prop="assign_to">
          <el-select v-model="assignForm.assign_to" placeholder="请选择维修工" filterable>
            <el-option
              v-for="user in technicians"
              :key="user.id"
              :label="user.name"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAssignDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmAssign" :loading="assigning">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import {
  repairOrderApi,
  getStatusText,
  getStatusType,
  getPriorityText,
  getPriorityType,
  type RepairOrder,
  type RepairStatistics,
  type AssignRepairRequest
} from '@/api/repair'

const router = useRouter()

const loading = ref(false)
const orders = ref<RepairOrder[]>([])
const currentOrder = ref<RepairOrder | null>(null)
const statistics = ref<RepairStatistics>({
  total_orders: 0,
  pending_orders: 0,
  in_progress_orders: 0,
  completed_orders: 0,
  today_completed: 0,
  today_created: 0,
  avg_repair_time: 0,
  avg_response_time: 0,
})

const dateRange = ref<[string, string] | null>(null)
const showDetailDialog = ref(false)
const showAssignDialog = ref(false)
const assigning = ref(false)

const filterForm = reactive({
  status: '',
  priority: undefined as number | undefined,
  date_from: '',
  date_to: ''
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const assignFormRef = ref<FormInstance>()
const assignForm = reactive<AssignRepairRequest>({
  assign_to: undefined as unknown as number
})

const assignRules: FormRules = {
  assign_to: [{ required: true, message: '请选择维修工', trigger: 'change' }]
}

// 模拟维修工列表
const technicians = ref([
  { id: 1, name: '维修工1' },
  { id: 2, name: '维修工2' },
  { id: 3, name: '维修工3' },
])

const formatDateTime = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const onDateChange = (dates: [string, string] | null) => {
  if (dates) {
    filterForm.date_from = dates[0]
    filterForm.date_to = dates[1]
  } else {
    filterForm.date_from = ''
    filterForm.date_to = ''
  }
  loadOrders()
}

const loadOrders = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (filterForm.status) params.status = filterForm.status
    if (filterForm.priority) params.priority = filterForm.priority
    if (filterForm.date_from) params.date_from = filterForm.date_from
    if (filterForm.date_to) params.date_to = filterForm.date_to

    const data = await repairOrderApi.getOrders(params)
    orders.value = data.items
    pagination.total = data.total
  } catch (error: any) {
    ElMessage.error(error.message || '加载工单列表失败')
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    statistics.value = await repairOrderApi.getStatistics()
  } catch (error: any) {
    console.error('加载统计数据失败', error)
  }
}

const viewDetail = async (order: RepairOrder) => {
  currentOrder.value = await repairOrderApi.getOrder(order.id)
  showDetailDialog.value = true
}

const showAssignDialogFunc = (order: RepairOrder) => {
  currentOrder.value = order
  assignForm.assign_to = undefined as unknown as number
  showAssignDialog.value = true
}

// 需要重新命名避免冲突
const showAssignDialogInner = (order: RepairOrder) => {
  showAssignDialogFunc(order)
}

const confirmAssign = async () => {
  if (!assignFormRef.value) return
  await assignFormRef.value.validate(async (valid) => {
    if (!valid || !currentOrder.value) return
    assigning.value = true
    try {
      await repairOrderApi.assignOrder(currentOrder.value.id, assignForm)
      ElMessage.success('派单成功')
      showAssignDialog.value = false
      loadOrders()
      loadStatistics()
    } catch (error: any) {
      ElMessage.error(error.message || '派单失败')
    } finally {
      assigning.value = false
    }
  })
}

const executeOrder = (order: RepairOrder) => {
  router.push({
    name: 'repair-execute',
    params: { orderId: order.id }
  })
}

onMounted(() => {
  loadOrders()
  loadStatistics()
})
</script>

<style scoped>
.order-list-view {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filter-form {
  margin-bottom: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.stat-pending { color: #909399; }
.stat-progress { color: #e6a23c; }
.stat-completed { color: #67c23a; }
.stat-created { color: #409eff; }

.order-table {
  margin-bottom: 16px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
}

.order-detail .photos-section {
  margin-top: 16px;
}

.section-label {
  font-weight: 500;
  margin-bottom: 8px;
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.photo-item {
  width: 100%;
  height: 100px;
  border-radius: 4px;
  cursor: pointer;
}
</style>
