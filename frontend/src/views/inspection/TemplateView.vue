<template>
  <div class="template-view">
    <el-row :gutter="20">
      <!-- 模板列表 -->
      <el-col :span="8">
        <el-card class="template-list-card">
          <template #header>
            <div class="card-header">
              <span>点检模板</span>
              <el-button type="primary" size="small" @click="showCreateDialog = true">
                <el-icon><Plus /></el-icon>
                新建模板
              </el-button>
            </div>
          </template>

          <el-table
            :data="templates"
            @row-click="selectTemplate"
            highlight-current-row
            v-loading="loading"
          >
            <el-table-column prop="name" label="模板名称" />
            <el-table-column prop="equipment_type_name" label="设备类型" width="100" />
            <el-table-column prop="item_count" label="项目数" width="60" align="center" />
            <el-table-column label="操作" width="80" align="center">
              <template #default="{ row }">
                <el-button
                  type="danger"
                  size="small"
                  link
                  @click.stop="deleteTemplate(row)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- 模板详情 -->
      <el-col :span="16">
        <el-card class="template-detail-card" v-if="currentTemplate">
          <template #header>
            <div class="card-header">
              <span>{{ currentTemplate.name }} - 点检项目</span>
              <el-button type="primary" size="small" @click="showAddItemDialog = true">
                <el-icon><Plus /></el-icon>
                添加项目
              </el-button>
            </div>
          </template>

          <el-table :data="templateItems" v-loading="itemsLoading">
            <el-table-column prop="sequence_order" label="序号" width="60" align="center" />
            <el-table-column prop="name" label="项目名称" />
            <el-table-column prop="method" label="检查方法" />
            <el-table-column prop="criteria" label="判定标准" />
            <el-table-column label="操作" width="100" align="center">
              <template #default="{ row }">
                <el-button type="primary" size="small" link @click="editItem(row)">
                  编辑
                </el-button>
                <el-button type="danger" size="small" link @click="deleteItem(row)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-empty v-else description="请选择一个模板查看详情" />
      </el-col>
    </el-row>

    <!-- 创建模板对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="新建点检模板"
      width="500px"
    >
      <el-form :model="templateForm" :rules="templateRules" ref="templateFormRef" label-width="100px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="设备类型" prop="equipment_type_id">
          <el-select
            v-model="templateForm.equipment_type_id"
            placeholder="请选择设备类型"
            style="width: 100%"
          >
            <el-option
              v-for="type in equipmentTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createTemplate" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加/编辑项目对话框 -->
    <el-dialog
      v-model="showAddItemDialog"
      :title="editingItem ? '编辑点检项目' : '添加点检项目'"
      width="600px"
    >
      <el-form :model="itemForm" :rules="itemRules" ref="itemFormRef" label-width="100px">
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="itemForm.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="检查序号" prop="sequence_order">
          <el-input-number v-model="itemForm.sequence_order" :min="1" />
        </el-form-item>
        <el-form-item label="检查方法" prop="method">
          <el-input
            v-model="itemForm.method"
            type="textarea"
            :rows="2"
            placeholder="请输入检查方法，如：目视检查、听诊等"
          />
        </el-form-item>
        <el-form-item label="判定标准" prop="criteria">
          <el-input
            v-model="itemForm.criteria"
            type="textarea"
            :rows="2"
            placeholder="请输入判定标准，如：无异常、无泄漏等"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddItemDialog = false">取消</el-button>
        <el-button type="primary" @click="saveItem" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  inspectionTemplateApi,
  inspectionItemApi,
  type InspectionTemplate,
  type InspectionTemplateDetail,
  type InspectionItem
} from '@/api/inspection'
import { equipmentApi, type EquipmentType } from '@/api/equipment'

const loading = ref(false)
const itemsLoading = ref(false)
const submitting = ref(false)
const templates = ref<InspectionTemplate[]>([])
const templateItems = ref<InspectionItem[]>([])
const currentTemplate = ref<InspectionTemplateDetail | null>(null)
const equipmentTypes = ref<EquipmentType[]>([])

const showCreateDialog = ref(false)
const showAddItemDialog = ref(false)
const editingItem = ref<InspectionItem | null>(null)

const templateFormRef = ref<FormInstance>()
const itemFormRef = ref<FormInstance>()

const templateForm = reactive({
  name: '',
  equipment_type_id: undefined as number | undefined
})

const itemForm = reactive({
  id: undefined as number | undefined,
  name: '',
  method: '',
  criteria: '',
  sequence_order: 1
})

const templateRules: FormRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  equipment_type_id: [{ required: true, message: '请选择设备类型', trigger: 'change' }]
}

const itemRules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  sequence_order: [{ required: true, message: '请输入检查序号', trigger: 'blur' }]
}

// 加载模板列表
const loadTemplates = async () => {
  loading.value = true
  try {
    const data = await inspectionTemplateApi.getTemplates()
    templates.value = data
  } catch (error: any) {
    ElMessage.error(error.message || '加载模板列表失败')
  } finally {
    loading.value = false
  }
}

// 加载设备类型
const loadEquipmentTypes = async () => {
  try {
    const data = await equipmentApi.getTypes()
    equipmentTypes.value = data
  } catch (error: any) {
    ElMessage.error(error.message || '加载设备类型失败')
  }
}

// 选择模板
const selectTemplate = async (row: InspectionTemplate) => {
  itemsLoading.value = true
  try {
    const data = await inspectionTemplateApi.getTemplate(row.id)
    currentTemplate.value = data
    templateItems.value = data.items || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载模板详情失败')
  } finally {
    itemsLoading.value = false
  }
}

// 创建模板
const createTemplate = async () => {
  if (!templateFormRef.value) return
  await templateFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await inspectionTemplateApi.createTemplate({
        name: templateForm.name,
        equipment_type_id: templateForm.equipment_type_id!
      })
      ElMessage.success('模板创建成功')
      showCreateDialog.value = false
      templateForm.name = ''
      templateForm.equipment_type_id = undefined
      loadTemplates()
    } catch (error: any) {
      ElMessage.error(error.message || '创建模板失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除模板
const deleteTemplate = async (row: InspectionTemplate) => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗？', '提示', {
      type: 'warning'
    })
    // TODO: 实现删除API
    ElMessage.info('删除功能待实现')
  } catch {
    // User cancelled
  }
}

// 编辑项目
const editItem = (row: InspectionItem) => {
  editingItem.value = row
  itemForm.id = row.id
  itemForm.name = row.name
  itemForm.method = row.method || ''
  itemForm.criteria = row.criteria || ''
  itemForm.sequence_order = row.sequence_order
  showAddItemDialog.value = true
}

// 保存项目
const saveItem = async () => {
  if (!itemFormRef.value) return
  if (!currentTemplate.value) return

  await itemFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (editingItem.value) {
        // TODO: 实现更新API
        ElMessage.info('更新功能待实现')
      } else {
        await inspectionItemApi.createItem({
          template_id: currentTemplate.value.id,
          name: itemForm.name,
          method: itemForm.method,
          criteria: itemForm.criteria,
          sequence_order: itemForm.sequence_order
        })
        ElMessage.success('添加成功')
      }
      showAddItemDialog.value = false
      resetItemForm()
      // 重新加载模板详情
      selectTemplate(currentTemplate.value)
    } catch (error: any) {
      ElMessage.error(error.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除项目
const deleteItem = async (row: InspectionItem) => {
  try {
    await ElMessageBox.confirm('确定要删除该项目吗？', '提示', {
      type: 'warning'
    })
    // TODO: 实现删除API
    ElMessage.info('删除功能待实现')
  } catch {
    // User cancelled
  }
}

const resetItemForm = () => {
  editingItem.value = null
  itemForm.id = undefined
  itemForm.name = ''
  itemForm.method = ''
  itemForm.criteria = ''
  itemForm.sequence_order = templateItems.value.length + 1
}

onMounted(() => {
  loadTemplates()
  loadEquipmentTypes()
})
</script>

<style scoped>
.template-view {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.template-list-card,
.template-detail-card {
  height: calc(100vh - 200px);
}

.template-detail-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

.template-list-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

:deep(.el-table__row) {
  cursor: pointer;
}
</style>
