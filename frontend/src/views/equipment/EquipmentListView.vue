<template>
  <div class="equipment-list">
    <el-card shadow="never" class="list-card">
      <template #header>
        <div class="card-header">
          <div class="title-wrapper">
            <el-icon class="title-icon"><Monitor /></el-icon>
            <span class="title-text">设备台账</span>
          </div>
          <el-button type="primary" @click="showCreateDialog = true" v-if="canEdit">
            <el-icon><Plus /></el-icon>
            新增设备
          </el-button>
        </div>
      </template>

      <!-- Search Form -->
      <div class="search-wrapper">
        <el-form :model="queryParams" label-width="80px">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-form-item label="设备编号">
                <el-input v-model="queryParams.code" placeholder="输入编号搜索" clearable @keyup.enter="handleSearch" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="设备名称">
                <el-input v-model="queryParams.name" placeholder="输入名称搜索" clearable @keyup.enter="handleSearch" />
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="设备类型">
                <el-select v-model="queryParams.type_id" placeholder="全部类型" clearable class="full-width">
                  <el-option
                    v-for="type in equipmentTypes"
                    :key="type.id"
                    :label="type.name"
                    :value="type.id"
                  />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item label="设备状态">
                <el-select v-model="queryParams.status" placeholder="全部状态" clearable class="full-width">
                  <el-option label="运行中" value="running" />
                  <el-option label="已停机" value="stopped" />
                  <el-option label="维修中" value="maintenance" />
                  <el-option label="已报废" value="scrapped" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <div class="search-buttons">
            <el-button type="primary" @click="handleSearch">
              <el-icon><Search /></el-icon>
              查询
            </el-button>
            <el-button @click="handleReset">
              <el-icon><Refresh /></el-icon>
              重置
            </el-button>
          </div>
        </el-form>
      </div>

      <!-- Equipment Table -->
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        class="custom-table"
      >
        <el-table-column type="index" label="序号" width="60" align="center" />
        <el-table-column prop="code" label="设备编号" width="140" show-overflow-tooltip />
        <el-table-column prop="name" label="设备名称" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="equipment-name-link" @click="handleView(row)">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type_name" label="设备类型" width="120" />
        <el-table-column prop="factory_name" label="工厂" width="120" show-overflow-tooltip />
        <el-table-column prop="workshop_name" label="车间" width="120" show-overflow-tooltip />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" effect="light" round>{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" align="center" fixed="right">
          <template #default="{ row }">
            <div class="operation-buttons">
              <el-button size="small" type="primary" link @click="handleView(row)">
                <el-icon><View /></el-icon> 查看
              </el-button>
              <el-button size="small" type="primary" link @click="handleEdit(row)" v-if="canEdit">
                <el-icon><Edit /></el-icon> 编辑
              </el-button>
              <el-button size="small" type="warning" link @click="handleQRCode(row)">
                <el-icon><Grid /></el-icon> 码
              </el-button>
              <el-button size="small" type="danger" link @click="handleDelete(row)" v-if="canDelete">
                <el-icon><Delete /></el-icon> 删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- Pagination -->
      <el-pagination
        v-model:current-page="queryParams.page"
        v-model:page-size="queryParams.page_size"
        :total="total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="loadData"
        @current-change="loadData"
      />
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingEquipment ? '编辑设备' : '新增设备'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="设备编号" prop="code">
          <el-input v-model="form.code" placeholder="请输入设备编号" />
        </el-form-item>
        <el-form-item label="设备名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入设备名称" />
        </el-form-item>
        <el-form-item label="设备类型" prop="type_id">
          <el-select v-model="form.type_id" placeholder="请选择设备类型" style="width: 100%">
            <el-option
              v-for="type in equipmentTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="所属车间" prop="workshop_id">
          <el-cascader
            v-model="workshopValue"
            :options="organizationTree"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            placeholder="请选择所属车间"
            style="width: 100%"
            @change="handleWorkshopChange"
          />
        </el-form-item>
        <el-form-item label="技术参数">
          <el-input v-model="form.spec" type="textarea" :rows="3" placeholder="请输入技术参数" />
        </el-form-item>
        <el-form-item label="采购日期">
          <el-date-picker
            v-model="form.purchase_date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="运行中" value="running" />
            <el-option label="已停机" value="stopped" />
            <el-option label="维修中" value="maintenance" />
            <el-option label="已报废" value="scrapped" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- QR Code Dialog -->
    <el-dialog v-model="showQRDialog" title="设备二维码" width="400px">
      <div class="qr-content">
        <div class="qr-placeholder">
          <el-icon :size="150"><Grid /></el-icon>
        </div>
        <p class="qr-code-text">{{ currentEquipment?.qr_code }}</p>
        <p class="qr-name">{{ currentEquipment?.name }}</p>
      </div>
      <template #footer>
        <el-button @click="showQRDialog = false">关闭</el-button>
        <el-button type="primary" @click="handleDownloadQR">下载二维码</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { equipmentApi, equipmentTypeApi, orgApi, type Equipment, type EquipmentType } from '@/api/equipment'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const canEdit = computed(() => authStore.hasRole('admin', 'engineer'))
const canDelete = computed(() => authStore.hasRole('admin'))

const loading = ref(false)
const tableData = ref<Equipment[]>([])
const total = ref(0)
const equipmentTypes = ref<EquipmentType[]>([])

const queryParams = reactive({
  page: 1,
  page_size: 20,
  code: '',
  name: '',
  type_id: undefined as number | undefined,
  status: '',
})

const showCreateDialog = ref(false)
const showQRDialog = ref(false)
const editingEquipment = ref<Equipment | null>(null)
const currentEquipment = ref<Equipment | null>(null)
const submitting = ref(false)

const formRef = ref<FormInstance>()
const form = reactive({
  code: '',
  name: '',
  type_id: undefined as number | undefined,
  workshop_id: undefined as number | undefined,
  spec: '',
  purchase_date: '',
  status: 'running',
  dedicated_maintenance_id: undefined as number | undefined,
})

const workshopValue = ref<number[]>([])

const formRules: FormRules = {
  code: [{ required: true, message: '请输入设备编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  type_id: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
  workshop_id: [{ required: true, message: '请选择所属车间', trigger: 'change' }],
}

// Organization tree for cascader
interface TreeNode {
  id: number
  name: string
  code: string
  type: 'base' | 'factory' | 'workshop'
  children?: TreeNode[]
}

const organizationTree = ref<TreeNode[]>([])

function getStatusType(status: string) {
  const map: Record<string, any> = {
    running: 'success',
    stopped: 'info',
    maintenance: 'warning',
    scrapped: 'danger',
  }
  return map[status] || ''
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停机',
    maintenance: '维修中',
    scrapped: '已报废',
  }
  return map[status] || status
}

async function loadData() {
  loading.value = true
  try {
    const response = await equipmentApi.getList(queryParams)
    tableData.value = response.data.items
    total.value = response.data.total
  } catch (error) {
    // Error handled by interceptor
  } finally {
    loading.value = false
  }
}

async function loadEquipmentTypes() {
  try {
    const response = await equipmentTypeApi.getTypes()
    equipmentTypes.value = response.data
  } catch (error) {
    // Error handled by interceptor
  }
}

async function loadOrganization() {
  try {
    const bases = await orgApi.getBases()
    const tree: TreeNode[] = []

    for (const base of bases) {
      const baseNode: TreeNode = { id: base.id, name: base.name, code: base.code, type: 'base' }
      const factories = await orgApi.getFactories(base.id)

      if (factories.length > 0) {
        baseNode.children = []
        for (const factory of factories) {
          const factoryNode: TreeNode = {
            id: factory.id,
            name: factory.name,
            code: factory.code,
            type: 'factory',
          }
          const workshops = await orgApi.getWorkshops(factory.id)

          if (workshops.length > 0) {
            factoryNode.children = workshops.map((w) => ({
              id: w.id,
              name: w.name,
              code: w.code,
              type: 'workshop' as const,
            }))
          }
          baseNode.children.push(factoryNode)
        }
      }
      tree.push(baseNode)
    }

    organizationTree.value = tree
  } catch (error) {
    // Error handled by interceptor
  }
}

function handleSearch() {
  queryParams.page = 1
  loadData()
}

function handleReset() {
  queryParams.code = ''
  queryParams.name = ''
  queryParams.type_id = undefined
  queryParams.status = ''
  queryParams.page = 1
  loadData()
}

function handleView(row: Equipment) {
  router.push(`/equipment/detail/${row.id}`)
}

function handleEdit(row: Equipment) {
  editingEquipment.value = row
  form.code = row.code
  form.name = row.name
  form.type_id = row.type_id
  form.workshop_id = row.workshop_id
  form.spec = row.spec || ''
  form.purchase_date = row.purchase_date || ''
  form.status = row.status
  form.dedicated_maintenance_id = row.dedicated_maintenance_id

  // Set cascader value
  workshopValue.value = [row.factory_id!, row.workshop_id]

  showCreateDialog.value = true
}

function handleWorkshopChange(value: number[]) {
  if (value && value.length > 0) {
    form.workshop_id = value[value.length - 1]
  }
}

async function handleDelete(row: Equipment) {
  try {
    await ElMessageBox.confirm(`确定要删除设备 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await equipmentApi.delete(row.id)
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // User cancelled or error
  }
}

function handleQRCode(row: Equipment) {
  currentEquipment.value = row
  showQRDialog.value = true
}

function handleDownloadQR() {
  ElMessage.info('二维码下载功能开发中')
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      const submitData = { ...form }
      if (submitData.purchase_date && submitData.purchase_date.length === 10) {
        submitData.purchase_date = `${submitData.purchase_date}T00:00:00Z`
      }

      if (editingEquipment.value) {
        await equipmentApi.update(editingEquipment.value.id, submitData)
        ElMessage.success('更新成功')
      } else {
        await equipmentApi.create(submitData)
        ElMessage.success('创建成功')
      }

      showCreateDialog.value = false
      resetForm()
      loadData()
    } catch (error) {
      // Error handled by interceptor
    } finally {
      submitting.value = false
    }
  })
}

function resetForm() {
  editingEquipment.value = null
  form.code = ''
  form.name = ''
  form.type_id = undefined
  form.workshop_id = undefined
  form.spec = ''
  form.purchase_date = ''
  form.status = 'running'
  form.dedicated_maintenance_id = undefined
  workshopValue.value = []
  formRef.value?.clearValidate()
}

onMounted(() => {
  loadData()
  loadEquipmentTypes()
  loadOrganization()
})
</script>

<style scoped>
.equipment-list {
  padding: 20px;
}

.list-card {
  border-radius: 8px;
  border: 1px solid #ebeef5;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title-wrapper {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  font-size: 20px;
  color: var(--el-color-primary);
}

.title-text {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.search-wrapper {
  background-color: #f8f9fb;
  padding: 24px 20px 0;
  border-radius: 8px;
  margin-bottom: 24px;
}

.full-width {
  width: 100%;
}

.search-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding-bottom: 20px;
  margin-top: 10px;
}

.custom-table {
  width: 100%;
  margin-top: 8px;
  border-radius: 8px;
  overflow: hidden;
}

.equipment-name-link {
  color: var(--el-color-primary);
  cursor: pointer;
  font-weight: 500;
}

.equipment-name-link:hover {
  text-decoration: underline;
}

.operation-buttons {
  display: flex;
  justify-content: center;
  gap: 4px;
}

.operation-buttons :deep(.el-button) {
  padding: 4px 8px;
  margin-left: 0 !important;
}

.qr-content {
  text-align: center;
  padding: 20px 0;
}

.qr-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  background-color: #f5f7fa;
  border-radius: 8px;
  margin: 0 auto;
  width: 200px;
}

.qr-code-text {
  margin-top: 20px;
  font-size: 18px;
  font-weight: bold;
  color: var(--el-color-primary);
  letter-spacing: 1px;
}

.qr-name {
  margin-top: 10px;
  font-size: 14px;
  color: #909399;
}
</style>
