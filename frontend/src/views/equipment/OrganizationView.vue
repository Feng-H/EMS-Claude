<template>
  <div class="organization">
    <el-row :gutter="20">
      <!-- Bases -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>基地列表</span>
              <el-button type="primary" size="small" @click="showBaseDialog = true">
                <el-icon><Plus /></el-icon>
                新增
              </el-button>
            </div>
          </template>
          <el-table :data="bases" border max-height="400">
            <el-table-column prop="code" label="编号" width="100" />
            <el-table-column prop="name" label="名称" />
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button size="small" text type="primary" @click="editBase(row)">编辑</el-button>
                <el-button size="small" text type="danger" @click="deleteBase(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- Factories -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>工厂列表</span>
              <el-button type="primary" size="small" @click="showFactoryDialog = true">
                <el-icon><Plus /></el-icon>
                新增
              </el-button>
            </div>
          </template>
          <el-table :data="factories" border max-height="400">
            <el-table-column prop="code" label="编号" width="100" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="base_name" label="所属基地" />
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button size="small" text type="primary" @click="editFactory(row)">编辑</el-button>
                <el-button size="small" text type="danger" @click="deleteFactory(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- Workshops -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>车间列表</span>
              <el-button type="primary" size="small" @click="showWorkshopDialog = true">
                <el-icon><Plus /></el-icon>
                新增
              </el-button>
            </div>
          </template>
          <el-table :data="workshops" border max-height="400">
            <el-table-column prop="code" label="编号" width="100" />
            <el-table-column prop="name" label="名称" />
            <el-table-column prop="factory_name" label="所属工厂" />
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button size="small" text type="primary" @click="editWorkshop(row)">编辑</el-button>
                <el-button size="small" text type="danger" @click="deleteWorkshop(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- Base Dialog -->
    <el-dialog v-model="showBaseDialog" :title="editingBase ? '编辑基地' : '新增基地'" width="500px">
      <el-form ref="baseFormRef" :model="baseForm" :rules="baseRules" label-width="100px">
        <el-form-item label="基地编号" prop="code">
          <el-input v-model="baseForm.code" placeholder="如: BASE01" />
        </el-form-item>
        <el-form-item label="基地名称" prop="name">
          <el-input v-model="baseForm.name" placeholder="如: 华东基地" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBaseDialog = false">取消</el-button>
        <el-button type="primary" @click="saveBase">确定</el-button>
      </template>
    </el-dialog>

    <!-- Factory Dialog -->
    <el-dialog v-model="showFactoryDialog" :title="editingFactory ? '编辑工厂' : '新增工厂'" width="500px">
      <el-form ref="factoryFormRef" :model="factoryForm" :rules="factoryRules" label-width="100px">
        <el-form-item label="所属基地" prop="base_id">
          <el-select v-model="factoryForm.base_id" placeholder="请选择基地" style="width: 100%">
            <el-option v-for="base in bases" :key="base.id" :label="base.name" :value="base.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="工厂编号" prop="code">
          <el-input v-model="factoryForm.code" placeholder="如: FAC01" />
        </el-form-item>
        <el-form-item label="工厂名称" prop="name">
          <el-input v-model="factoryForm.name" placeholder="如: 苏州工厂" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFactoryDialog = false">取消</el-button>
        <el-button type="primary" @click="saveFactory">确定</el-button>
      </template>
    </el-dialog>

    <!-- Workshop Dialog -->
    <el-dialog v-model="showWorkshopDialog" :title="editingWorkshop ? '编辑车间' : '新增车间'" width="500px">
      <el-form ref="workshopFormRef" :model="workshopForm" :rules="workshopRules" label-width="100px">
        <el-form-item label="所属工厂" prop="factory_id">
          <el-select v-model="workshopForm.factory_id" placeholder="请选择工厂" style="width: 100%">
            <el-option v-for="factory in factories" :key="factory.id" :label="factory.name" :value="factory.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="车间编号" prop="code">
          <el-input v-model="workshopForm.code" placeholder="如: WS01" />
        </el-form-item>
        <el-form-item label="车间名称" prop="name">
          <el-input v-model="workshopForm.name" placeholder="如: 机加车间" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showWorkshopDialog = false">取消</el-button>
        <el-button type="primary" @click="saveWorkshop">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { orgApi, type Base, Factory, Workshop } from '@/api/equipment'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const bases = ref<Base[]>([])
const factories = ref<Factory[]>([])
const workshops = ref<Workshop[]>([])

// Base
const showBaseDialog = ref(false)
const editingBase = ref<Base | null>(null)
const baseFormRef = ref<FormInstance>()
const baseForm = reactive({ code: '', name: '' })
const baseRules: FormRules = {
  code: [{ required: true, message: '请输入编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

// Factory
const showFactoryDialog = ref(false)
const editingFactory = ref<Factory | null>(null)
const factoryFormRef = ref<FormInstance>()
const factoryForm = reactive({ base_id: undefined as number | undefined, code: '', name: '' })
const factoryRules: FormRules = {
  base_id: [{ required: true, message: '请选择基地', trigger: 'change' }],
  code: [{ required: true, message: '请输入编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

// Workshop
const showWorkshopDialog = ref(false)
const editingWorkshop = ref<Workshop | null>(null)
const workshopFormRef = ref<FormInstance>()
const workshopForm = reactive({ factory_id: undefined as number | undefined, code: '', name: '' })
const workshopRules: FormRules = {
  factory_id: [{ required: true, message: '请选择工厂', trigger: 'change' }],
  code: [{ required: true, message: '请输入编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

async function loadBases() {
  try {
    bases.value = await orgApi.getBases()
  } catch (error) {
    // Error handled by interceptor
  }
}

async function loadFactories() {
  try {
    factories.value = await orgApi.getFactories()
  } catch (error) {
    // Error handled by interceptor
  }
}

async function loadWorkshops() {
  try {
    workshops.value = await orgApi.getWorkshops()
  } catch (error) {
    // Error handled by interceptor
  }
}

// Base CRUD
function editBase(row: Base) {
  editingBase.value = row
  baseForm.code = row.code
  baseForm.name = row.name
  showBaseDialog.value = true
}

async function saveBase() {
  if (!baseFormRef.value) return

  await baseFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (editingBase.value) {
        await orgApi.updateBase(editingBase.value.id, baseForm)
        ElMessage.success('更新成功')
      } else {
        await orgApi.createBase(baseForm)
        ElMessage.success('创建成功')
      }

      showBaseDialog.value = false
      resetBaseForm()
      loadBases()
    } catch (error) {
      // Error handled by interceptor
    }
  })
}

function resetBaseForm() {
  editingBase.value = null
  baseForm.code = ''
  baseForm.name = ''
  baseFormRef.value?.clearValidate()
}

async function deleteBase(row: Base) {
  try {
    await ElMessageBox.confirm(`确定要删除基地 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await orgApi.deleteBase(row.id)
    ElMessage.success('删除成功')
    loadBases()
    loadFactories()
    loadWorkshops()
  } catch {
    // User cancelled
  }
}

// Factory CRUD
function editFactory(row: Factory) {
  editingFactory.value = row
  factoryForm.base_id = row.base_id
  factoryForm.code = row.code
  factoryForm.name = row.name
  showFactoryDialog.value = true
}

async function saveFactory() {
  if (!factoryFormRef.value) return

  await factoryFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (editingFactory.value) {
        await orgApi.updateFactory(editingFactory.value.id, factoryForm)
        ElMessage.success('更新成功')
      } else {
        await orgApi.createFactory(factoryForm)
        ElMessage.success('创建成功')
      }

      showFactoryDialog.value = false
      resetFactoryForm()
      loadFactories()
    } catch (error) {
      // Error handled by interceptor
    }
  })
}

function resetFactoryForm() {
  editingFactory.value = null
  factoryForm.base_id = undefined
  factoryForm.code = ''
  factoryForm.name = ''
  factoryFormRef.value?.clearValidate()
}

async function deleteFactory(row: Factory) {
  try {
    await ElMessageBox.confirm(`确定要删除工厂 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await orgApi.deleteFactory(row.id)
    ElMessage.success('删除成功')
    loadFactories()
    loadWorkshops()
  } catch {
    // User cancelled
  }
}

// Workshop CRUD
function editWorkshop(row: Workshop) {
  editingWorkshop.value = row
  workshopForm.factory_id = row.factory_id
  workshopForm.code = row.code
  workshopForm.name = row.name
  showWorkshopDialog.value = true
}

async function saveWorkshop() {
  if (!workshopFormRef.value) return

  await workshopFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (editingWorkshop.value) {
        await orgApi.updateWorkshop(editingWorkshop.value.id, workshopForm)
        ElMessage.success('更新成功')
      } else {
        await orgApi.createWorkshop(workshopForm)
        ElMessage.success('创建成功')
      }

      showWorkshopDialog.value = false
      resetWorkshopForm()
      loadWorkshops()
    } catch (error) {
      // Error handled by interceptor
    }
  })
}

function resetWorkshopForm() {
  editingWorkshop.value = null
  workshopForm.factory_id = undefined
  workshopForm.code = ''
  workshopForm.name = ''
  workshopFormRef.value?.clearValidate()
}

async function deleteWorkshop(row: Workshop) {
  try {
    await ElMessageBox.confirm(`确定要删除车间 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await orgApi.deleteWorkshop(row.id)
    ElMessage.success('删除成功')
    loadWorkshops()
  } catch {
    // User cancelled
  }
}

onMounted(() => {
  loadBases()
  loadFactories()
  loadWorkshops()
})
</script>

<style scoped>
.organization {
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
