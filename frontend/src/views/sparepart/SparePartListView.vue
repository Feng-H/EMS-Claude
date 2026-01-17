<template>
  <div class="sparepart-view">
    <div class="header">
      <h2>备件管理</h2>
      <el-button type="primary" @click="showCreateDialog = true" v-if="canManage">
        <el-icon><Plus /></el-icon>
        新增备件
      </el-button>
    </div>

    <!-- Statistics -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card>
          <stat-item icon="Goods" label="备件总数" :value="stats.total_parts" color="#409eff" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <stat-item icon="Warning" label="库存预警" :value="stats.low_stock_count" color="#f56c6c" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <stat-item icon="ShoppingCart" label="本月消耗" :value="stats.monthly_consumption" color="#e6a23c" />
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <stat-item icon="Coin" label="库存总值" :value="formatCurrency(stats.total_stock_value)" color="#67c23a" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Tabs -->
    <el-tabs v-model="activeTab" class="tabs">
      <!-- Parts List -->
      <el-tab-pane label="备件列表" name="parts">
        <el-card>
          <el-form :inline="true" :model="filterForm" class="filter-form">
            <el-form-item label="编码">
              <el-input v-model="filterForm.code" placeholder="请输入编码" clearable />
            </el-form-item>
            <el-form-item label="名称">
              <el-input v-model="filterForm.name" placeholder="请输入名称" clearable />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSearch">查询</el-button>
              <el-button @click="handleReset">重置</el-button>
            </el-form-item>
          </el-form>

          <el-table :data="parts" v-loading="loading" stripe>
            <el-table-column prop="code" label="编码" width="120" />
            <el-table-column prop="name" label="名称" min-width="150" />
            <el-table-column prop="specification" label="规格" min-width="120" />
            <el-table-column prop="unit" label="单位" width="80" />
            <el-table-column prop="factory_name" label="所属工厂" width="120" />
            <el-table-column prop="safety_stock" label="安全库存" width="100" align="right" />
            <el-table-column prop="current_stock" label="当前库存" width="100" align="right">
              <template #default="{ row }">
                <span :class="{ 'low-stock': row.current_stock < row.safety_stock }">
                  {{ row.current_stock || 0 }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right" v-if="canManage">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
                <el-button link type="primary" size="small" @click="showStockDialog(row)">入库</el-button>
                <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.id)">
                  <template #reference>
                    <el-button link type="danger" size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="pagination.page"
              v-model:page-size="pagination.pageSize"
              :total="pagination.total"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next"
              @size-change="loadParts"
              @current-change="loadParts"
            />
          </div>
        </el-card>
      </el-tab-pane>

      <!-- Inventory -->
      <el-tab-pane label="库存管理" name="inventory">
        <el-card>
          <el-table :data="inventory" v-loading="loading" stripe>
            <el-table-column prop="spare_part_code" label="备件编码" width="120" />
            <el-table-column prop="spare_part_name" label="备件名称" min-width="150" />
            <el-table-column prop="factory_name" label="工厂" width="120" />
            <el-table-column prop="quantity" label="库存数量" width="100" align="right">
              <template #default="{ row }">
                <span :class="{ 'low-stock': row.is_low_stock }">{{ row.quantity }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.is_low_stock" type="danger" size="small">低库存</el-tag>
                <el-tag v-else type="success" size="small">正常</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="updated_at" label="更新时间" width="160">
              <template #default="{ row }">{{ formatDateTime(row.updated_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- Consumption Records -->
      <el-tab-pane label="消耗记录" name="consumption">
        <el-card>
          <el-table :data="consumptions" v-loading="loading" stripe>
            <el-table-column prop="spare_part_code" label="备件编码" width="120" />
            <el-table-column prop="spare_part_name" label="备件名称" min-width="150" />
            <el-table-column prop="quantity" label="消耗数量" width="100" align="right" />
            <el-table-column prop="user_name" label="操作人" width="100" />
            <el-table-column prop="created_at" label="消耗时间" width="160">
              <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="150" />
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- Alerts -->
      <el-tab-pane label="库存预警" name="alerts">
        <el-card>
          <el-table :data="alerts" v-loading="loading" stripe>
            <el-table-column prop="spare_part_code" label="备件编码" width="120" />
            <el-table-column prop="spare_part_name" label="备件名称" min-width="150" />
            <el-table-column prop="factory_name" label="工厂" width="120" />
            <el-table-column prop="current_stock" label="当前库存" width="100" align="right">
              <template #default="{ row }"><span class="low-stock">{{ row.current_stock }}</span></template>
            </el-table-column>
            <el-table-column prop="safety_stock" label="安全库存" width="100" align="right" />
            <el-table-column prop="shortage" label="缺货数量" width="100" align="right">
              <template #default="{ row }"><span class="shortage">{{ row.shortage }}</span></template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingPart ? '编辑备件' : '新增备件'"
      width="500px"
    >
      <el-form :model="partForm" :rules="partRules" ref="partFormRef" label-width="100px">
        <el-form-item label="编码" prop="code">
          <el-input v-model="partForm.code" placeholder="请输入编码" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="partForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="规格">
          <el-input v-model="partForm.specification" placeholder="请输入规格" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="partForm.unit" placeholder="如：个、件、套" />
        </el-form-item>
        <el-form-item label="安全库存">
          <el-input-number v-model="partForm.safety_stock" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">确定</el-button>
      </template>
    </el-dialog>

    <!-- Stock In Dialog -->
    <el-dialog v-model="showStockInDialog" title="入库" width="400px">
      <el-form :model="stockForm" ref="stockFormRef" label-width="100px">
        <el-form-item label="备件">
          <span>{{ currentPart?.code }} - {{ currentPart?.name }}</span>
        </el-form-item>
        <el-form-item label="工厂" prop="factory_id">
          <el-select v-model="stockForm.factory_id" placeholder="请选择工厂" style="width: 100%">
            <el-option
              v-for="factory in factories"
              :key="factory.id"
              :label="factory.name"
              :value="factory.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="数量" prop="quantity">
          <el-input-number v-model="stockForm.quantity" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="stockForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showStockInDialog = false">取消</el-button>
        <el-button type="primary" @click="handleStockIn" :loading="stocking">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import {
  getSpareParts,
  createSparePart,
  updateSparePart,
  deleteSparePart,
  getInventory,
  stockIn,
  getConsumptions,
  getLowStockAlerts,
  getSparePartStatistics,
  type SparePart,
  type CreateSparePartRequest,
  type SparePartStatistics
} from '@/api/sparepart'
import { equipmentApi } from '@/api/equipment'
import type { Factory } from '@/api/equipment'
import StatItem from '@/components/StatItem.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const canManage = computed(() => authStore.hasRole('admin', 'engineer'))

const loading = ref(false)
const saving = ref(false)
const stocking = ref(false)
const activeTab = ref('parts')

const parts = ref<SparePart[]>([])
const inventory = ref<any[]>([])
const consumptions = ref<any[]>([])
const alerts = ref<any[]>([])
const factories = ref<Factory[]>([])

const stats = ref<SparePartStatistics>({
  total_parts: 0,
  low_stock_count: 0,
  total_stock_value: 0,
  monthly_consumption: 0
})

const filterForm = reactive({
  code: '',
  name: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// Dialogs
const showCreateDialog = ref(false)
const showStockInDialog = ref(false)
const editingPart = ref<SparePart | null>(null)
const currentPart = ref<SparePart | null>(null)

const partFormRef = ref<FormInstance>()
const partForm = reactive<CreateSparePartRequest>({
  code: '',
  name: '',
  specification: '',
  unit: '',
  safety_stock: 10
})

const stockForm = reactive({
  spare_part_id: 0,
  factory_id: 0,
  quantity: 1,
  remark: ''
})

const partRules: FormRules = {
  code: [{ required: true, message: '请输入编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
}

// Methods
const formatCurrency = (val: number) => {
  return new Intl.NumberFormat('zh-CN').format(val)
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadStats = async () => {
  try {
    const res = await getSparePartStatistics()
    stats.value = res.data || {
      total_parts: 0,
      low_stock_count: 0,
      total_stock_value: 0,
      monthly_consumption: 0
    }
  } catch (err) {
    console.error('Failed to load stats:', err)
  }
}

const loadParts = async () => {
  loading.value = true
  try {
    const res = await getSpareParts({
      code: filterForm.code || undefined,
      name: filterForm.name || undefined,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    parts.value = res.data.items
    pagination.total = res.data.total
  } catch (err) {
    ElMessage.error('加载备件列表失败')
  } finally {
    loading.value = false
  }
}

const loadInventory = async () => {
  loading.value = true
  try {
    const res = await getInventory({ page: 1, page_size: 100 })
    inventory.value = res.data.items
  } catch (err) {
    ElMessage.error('加载库存失败')
  } finally {
    loading.value = false
  }
}

const loadConsumptions = async () => {
  loading.value = true
  try {
    const res = await getConsumptions({ page: 1, page_size: 50 })
    consumptions.value = res.data.items
  } catch (err) {
    ElMessage.error('加载消耗记录失败')
  } finally {
    loading.value = false
  }
}

const loadAlerts = async () => {
  loading.value = true
  try {
    const res = await getLowStockAlerts()
    alerts.value = res.data
  } catch (err) {
    ElMessage.error('加载预警失败')
  } finally {
    loading.value = false
  }
}

const loadFactories = async () => {
  try {
    const res = await equipmentApi.getFactories()
    factories.value = res.data
  } catch (err) {
    console.error('Failed to load factories:', err)
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadParts()
}

const handleReset = () => {
  filterForm.code = ''
  filterForm.name = ''
  pagination.page = 1
  loadParts()
}

const handleEdit = (part: SparePart) => {
  editingPart.value = part
  Object.assign(partForm, {
    code: part.code,
    name: part.name,
    specification: part.specification || '',
    unit: part.unit || '',
    safety_stock: part.safety_stock
  })
  showCreateDialog.value = true
}

const handleSave = async () => {
  if (!partFormRef.value) return
  await partFormRef.value.validate(async (valid) => {
    if (!valid) return
    saving.value = true
    try {
      if (editingPart.value) {
        await updateSparePart(editingPart.value.id, partForm)
        ElMessage.success('更新成功')
      } else {
        await createSparePart(partForm)
        ElMessage.success('创建成功')
      }
      showCreateDialog.value = false
      loadParts()
      loadStats()
    } catch (err: any) {
      ElMessage.error(err.response?.data?.error || '操作失败')
    } finally {
      saving.value = false
    }
  })
}

const handleDelete = async (id: number) => {
  try {
    await deleteSparePart(id)
    ElMessage.success('删除成功')
    loadParts()
    loadStats()
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '删除失败')
  }
}

const showStockDialog = (part: SparePart) => {
  currentPart.value = part
  stockForm.spare_part_id = part.id
  stockForm.factory_id = 0
  stockForm.quantity = 1
  stockForm.remark = ''
  showStockInDialog.value = true
}

const handleStockIn = async () => {
  if (stockForm.factory_id === 0) {
    ElMessage.warning('请选择工厂')
    return
  }
  stocking.value = true
  try {
    await stockIn(stockForm)
    ElMessage.success('入库成功')
    showStockInDialog.value = false
    loadInventory()
    loadStats()
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '入库失败')
  } finally {
    stocking.value = false
  }
}

onMounted(() => {
  loadStats()
  loadParts()
  loadFactories()
})
</script>

<style scoped>
.sparepart-view {
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

.tabs {
  margin-bottom: 20px;
}

.filter-form {
  margin-bottom: 16px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.low-stock {
  color: #f56c6c;
  font-weight: bold;
}

.shortage {
  color: #f56c6c;
  font-weight: bold;
}
</style>
