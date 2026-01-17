<template>
  <div class="maintenance-plan-view">
    <div class="header">
      <h2>保养计划配置</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        新增计划
      </el-button>
    </div>

    <!-- Statistics -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card>
          <stat-item icon="Calendar" label="计划总数" :value="stats.totalPlans" color="#409eff" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <stat-item icon="Clock" label="待执行任务" :value="stats.pendingTasks" color="#e6a23c" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <stat-item icon="CircleCheck" label="已完成任务" :value="stats.completedTasks" color="#67c23a" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <stat-item icon="TrendCharts" label="完成率" :value="stats.completionRate + '%'" color="#909399" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Plans List -->
    <el-card class="table-card">
      <el-table :data="plans" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="计划名称" min-width="180" />
        <el-table-column prop="equipment_type_name" label="设备类型" width="150" />
        <el-table-column label="保养级别" width="120">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)">{{ row.level_name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="周期" width="150">
          <template #default="{ row }">
            {{ row.cycle_days }}天 / 弹性{{ row.flexible_days }}天
          </template>
        </el-table-column>
        <el-table-column prop="work_hours" label="工时定额" width="100">
          <template #default="{ row }">{{ row.work_hours }}h</template>
        </el-table-column>
        <el-table-column prop="item_count" label="保养项数" width="100" align="center" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewItems(row)">
              查看项目
            </el-button>
            <el-button link type="primary" size="small" @click="generateTasks(row)">
              生成任务
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create Plan Dialog -->
    <el-dialog v-model="showCreateDialog" title="新增保养计划" width="500px">
      <el-form :model="planForm" :rules="planRules" ref="planFormRef" label-width="100px">
        <el-form-item label="计划名称" prop="name">
          <el-input v-model="planForm.name" placeholder="请输入计划名称" />
        </el-form-item>
        <el-form-item label="设备类型" prop="equipment_type_id">
          <el-select v-model="planForm.equipment_type_id" placeholder="请选择设备类型" style="width: 100%">
            <el-option
              v-for="type in equipmentTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="保养级别" prop="level">
          <el-select v-model="planForm.level" placeholder="请选择保养级别" style="width: 100%">
            <el-option label="一级保养" :value="1" />
            <el-option label="二级保养" :value="2" />
            <el-option label="精度保养" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="周期天数" prop="cycle_days">
          <el-input-number v-model="planForm.cycle_days" :min="1" :max="365" style="width: 100%" />
        </el-form-item>
        <el-form-item label="弹性窗口" prop="flexible_days">
          <el-input-number v-model="planForm.flexible_days" :min="0" :max="30" style="width: 100%" />
        </el-form-item>
        <el-form-item label="工时定额" prop="work_hours">
          <el-input-number v-model="planForm.work_hours" :min="0" :max="100" :precision="1" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreatePlan" :loading="creating">确定</el-button>
      </template>
    </el-dialog>

    <!-- View Items Dialog -->
    <el-dialog v-model="showItemsDialog" title="保养项目" width="600px">
      <div class="items-header">
        <el-button type="primary" size="small" @click="showAddItemDialog = true">
          <el-icon><Plus /></el-icon>
          添加项目
        </el-button>
      </div>
      <el-table :data="planItems" stripe max-height="400">
        <el-table-column prop="sequence_order" label="序号" width="80" align="center" />
        <el-table-column prop="name" label="项目名称" min-width="150" />
        <el-table-column prop="method" label="保养方法" min-width="150" />
        <el-table-column prop="criteria" label="验收标准" min-width="150" />
      </el-table>

      <!-- Add Item Dialog -->
      <el-dialog v-model="showAddItemDialog" title="添加保养项目" width="500px" append-to-body>
        <el-form :model="itemForm" :rules="itemRules" ref="itemFormRef" label-width="100px">
          <el-form-item label="项目名称" prop="name">
            <el-input v-model="itemForm.name" placeholder="请输入项目名称" />
          </el-form-item>
          <el-form-item label="保养方法" prop="method">
            <el-input v-model="itemForm.method" type="textarea" :rows="3" placeholder="请输入保养方法" />
          </el-form-item>
          <el-form-item label="验收标准" prop="criteria">
            <el-input v-model="itemForm.criteria" type="textarea" :rows="3" placeholder="请输入验收标准" />
          </el-form-item>
          <el-form-item label="序号" prop="sequence_order">
            <el-input-number v-model="itemForm.sequence_order" :min="1" style="width: 100%" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showAddItemDialog = false">取消</el-button>
          <el-button type="primary" @click="handleAddItem" :loading="addingItem">确定</el-button>
        </template>
      </el-dialog>
    </el-dialog>

    <!-- Generate Tasks Dialog -->
    <el-dialog v-model="showGenerateDialog" title="生成保养任务" width="500px">
      <el-form :model="generateForm" ref="generateFormRef" label-width="100px">
        <el-form-item label="基准日期">
          <el-date-picker
            v-model="generateForm.date"
            type="date"
            placeholder="选择基准日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="选择设备">
          <el-select
            v-model="generateForm.equipment_ids"
            multiple
            filterable
            placeholder="请选择设备"
            style="width: 100%"
          >
            <el-option
              v-for="eq in filteredEquipment"
              :key="eq.id"
              :label="`${eq.code} - ${eq.name}`"
              :value="eq.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showGenerateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleGenerateTasks" :loading="generating">生成任务</el-button>
      </template>
    </el-dialog>

    <!-- Result Dialog -->
    <el-dialog v-model="showResultDialog" title="生成结果" width="400px">
      <el-result :icon="generateResult.created_count > 0 ? 'success' : 'warning'">
        <template #title>
          成功生成 {{ generateResult.created_count }} 个任务
        </template>
        <template #sub-title v-if="generateResult.errors?.length">
          <el-alert type="warning" :closable="false">
            <template #title>
              <div v-for="(err, i) in generateResult.errors" :key="i">{{ err }}</div>
            </template>
          </el-alert>
        </template>
        <template #extra>
          <el-button type="primary" @click="showResultDialog = false">关闭</el-button>
        </template>
      </el-result>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import {
  getMaintenancePlans,
  createMaintenancePlan,
  createMaintenanceItem,
  generateMaintenanceTasks,
  getMaintenanceStatistics,
  getLevelName,
  type MaintenancePlan,
  type MaintenanceItem,
  type CreateMaintenancePlanRequest,
  type CreateMaintenanceItemRequest,
  type GenerateMaintenanceTasksRequest,
  type GenerateMaintenanceTasksResponse,
  type MaintenanceStatistics
} from '@/api/maintenance'
import { equipmentTypeApi, equipmentApi } from '@/api/equipment'
import type { EquipmentType, Equipment } from '@/api/equipment'
import StatItem from '@/components/StatItem.vue'

const loading = ref(false)
const plans = ref<MaintenancePlan[]>([])
const equipmentTypes = ref<EquipmentType[]>([])
const equipment = ref<Equipment[]>([])
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

// Dialogs
const showCreateDialog = ref(false)
const showItemsDialog = ref(false)
const showAddItemDialog = ref(false)
const showGenerateDialog = ref(false)
const showResultDialog = ref(false)

// Forms
const planFormRef = ref<FormInstance>()
const itemFormRef = ref<FormInstance>()
const generating = ref(false)
const creating = ref(false)
const addingItem = ref(false)

const planForm = reactive<CreateMaintenancePlanRequest>({
  name: '',
  equipment_type_id: 0,
  level: 1,
  cycle_days: 30,
  flexible_days: 3,
  work_hours: 4
})

const itemForm = reactive<CreateMaintenanceItemRequest>({
  plan_id: 0,
  name: '',
  method: '',
  criteria: '',
  sequence_order: 1
})

const generateForm = reactive<{
  plan_id: number
  date: string
  equipment_ids: number[]
}>({
  plan_id: 0,
  date: new Date().toISOString().split('T')[0],
  equipment_ids: []
})

const generateResult = ref<GenerateMaintenanceTasksResponse>({
  created_count: 0,
  task_ids: [],
  errors: []
})

const planItems = ref<MaintenanceItem[]>([])
const currentPlan = ref<MaintenancePlan | null>(null)

const filteredEquipment = computed(() => {
  if (!currentPlan.value) return []
  return equipment.value.filter(eq => eq.type_id === currentPlan.value!.equipment_type_id)
})

// Validation rules
const planRules: FormRules = {
  name: [{ required: true, message: '请输入计划名称', trigger: 'blur' }],
  equipment_type_id: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
  level: [{ required: true, message: '请选择保养级别', trigger: 'change' }],
  cycle_days: [{ required: true, message: '请输入周期天数', trigger: 'blur' }]
}

const itemRules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  sequence_order: [{ required: true, message: '请输入序号', trigger: 'blur' }]
}

// Methods
const getLevelType = (level: number) => {
  const types: Record<number, string> = { 1: '', 2: 'warning', 3: 'success' }
  return types[level] || ''
}

const loadData = async () => {
  loading.value = true
  try {
    const [plansRes, typesRes, eqRes, statsRes] = await Promise.all([
      getMaintenancePlans(),
      equipmentTypeApi.getTypes(),
      equipmentApi.getList({ page: 1, page_size: 100 }), // Reduced page size to avoid 400 error
      getMaintenanceStatistics()
    ])
    plans.value = plansRes.data
    equipmentTypes.value = typesRes.data
    equipment.value = eqRes.data.items
    stats.value = statsRes.data || {
      total_plans: 0,
      total_tasks: 0,
      pending_tasks: 0,
      in_progress_tasks: 0,
      completed_tasks: 0,
      overdue_tasks: 0,
      today_completed: 0,
      completion_rate: 0
    } // Safe fallback
  } catch (err) {
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

const handleCreatePlan = async () => {
  if (!planFormRef.value) return
  await planFormRef.value.validate(async (valid) => {
    if (!valid) return
    creating.value = true
    try {
      await createMaintenancePlan(planForm)
      ElMessage.success('创建成功')
      showCreateDialog.value = false
      loadData()
      // Reset form
      Object.assign(planForm, {
        name: '',
        equipment_type_id: 0,
        level: 1,
        cycle_days: 30,
        flexible_days: 3,
        work_hours: 4
      })
    } catch (err: any) {
      ElMessage.error(err.response?.data?.error || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

const viewItems = async (plan: MaintenancePlan) => {
  currentPlan.value = plan
  // In real app, fetch items by plan_id
  planItems.value = []
  showItemsDialog.value = true
}

const handleAddItem = async () => {
  if (!itemFormRef.value || !currentPlan.value) return
  await itemFormRef.value.validate(async (valid) => {
    if (!valid) return
    addingItem.value = true
    try {
      itemForm.plan_id = currentPlan.value!.id
      const res = await createMaintenanceItem(itemForm)
      planItems.value.push(res.data)
      ElMessage.success('添加成功')
      showAddItemDialog.value = false
      // Reset form
      Object.assign(itemForm, {
        plan_id: 0,
        name: '',
        method: '',
        criteria: '',
        sequence_order: planItems.value.length + 1
      })
    } catch (err: any) {
      ElMessage.error(err.response?.data?.error || '添加失败')
    } finally {
      addingItem.value = false
    }
  })
}

const generateTasks = (plan: MaintenancePlan) => {
  currentPlan.value = plan
  generateForm.plan_id = plan.id
  generateForm.equipment_ids = []
  showGenerateDialog.value = true
}

const handleGenerateTasks = async () => {
  if (generateForm.equipment_ids.length === 0) {
    ElMessage.warning('请选择至少一个设备')
    return
  }
  generating.value = true
  try {
    const res = await generateMaintenanceTasks(generateForm)
    generateResult.value = res.data
    showGenerateDialog.value = false
    showResultDialog.value = true
    if (res.data.created_count > 0) {
      loadData()
    }
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '生成失败')
  } finally {
    generating.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.maintenance-plan-view {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.stats-row {
  margin-bottom: 20px;
}

.table-card {
  margin-bottom: 20px;
}

.items-header {
  margin-bottom: 16px;
}
</style>
