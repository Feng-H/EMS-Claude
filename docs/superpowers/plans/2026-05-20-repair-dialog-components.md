# Standalone Repair Dialog Components Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create 4 standalone Element Plus dialog components for the repair management loop.

**Architecture:** Standalone Vue 3 components using Composition API and Element Plus. Components will be placed in `frontend/src/views/repair/components/`.

**Tech Stack:** Vue 3, TypeScript, Element Plus, Vite.

---

### Task 1: Create RepairReportDialog.vue

**Files:**
- Create: `frontend/src/views/repair/components/RepairReportDialog.vue`

- [ ] **Step 1: Implement the component**
- Implement a dialog with equipment selection (searchable select), fault description (textarea), fault code (input), priority (radio), and mock upload.
- Fetch equipment list from `equipmentApi.getList`.
- Call `repairOrderApi.createOrder` on submit.

```vue
<template>
  <el-dialog
    :model-value="visible"
    title="报修申请"
    width="500px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="设备" prop="equipment_id">
        <el-select
          v-model="form.equipment_id"
          placeholder="搜索设备"
          filterable
          remote
          :remote-method="searchEquipment"
          :loading="loadingEquipment"
        >
          <el-option
            v-for="item in equipmentList"
            :key="item.id"
            :label="`${item.name} (${item.code})`"
            :value="item.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="故障描述" prop="fault_description">
        <el-input v-model="form.fault_description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="故障代码" prop="fault_code">
        <el-input v-model="form.fault_code" />
      </el-form-item>
      <el-form-item label="优先级" prop="priority">
        <el-radio-group v-model="form.priority">
          <el-radio :label="1">高</el-radio>
          <el-radio :label="2">中</el-radio>
          <el-radio :label="3">低</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="故障照片">
        <el-upload
          action="#"
          list-type="picture-card"
          :auto-upload="false"
          :on-change="handlePhotoChange"
          :on-remove="handlePhotoRemove"
        >
          <el-icon><Plus /></el-icon>
        </el-upload>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { repairOrderApi } from '@/api/repair'
import { equipmentApi, type Equipment } from '@/api/equipment'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const loadingEquipment = ref(false)
const equipmentList = ref<Equipment[]>([])

const form = reactive({
  equipment_id: undefined as number | undefined,
  fault_description: '',
  fault_code: '',
  priority: 2,
  photos: [] as string[]
})

const rules: FormRules = {
  equipment_id: [{ required: true, message: '请选择设备', trigger: 'change' }],
  fault_description: [{ required: true, message: '请输入故障描述', trigger: 'blur' }]
}

const searchEquipment = async (query: string) => {
  if (query) {
    loadingEquipment.value = true
    try {
      const res = await equipmentApi.getList({ name: query, page: 1, page_size: 20 })
      equipmentList.value = res.data.items
    } finally {
      loadingEquipment.value = false
    }
  }
}

const handlePhotoChange = (file: any) => {
  // Simple mock: store file object or base64
  // In real app, you would upload to server
}

const handlePhotoRemove = (file: any) => {
}

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await repairOrderApi.createOrder({
        equipment_id: form.equipment_id!,
        fault_description: form.fault_description,
        fault_code: form.fault_code,
        priority: form.priority,
        photos: [] // Mocked for now
      })
      ElMessage.success('报修申请提交成功')
      emit('success')
      emit('update:visible', false)
    } catch (error: any) {
      ElMessage.error(error.message || '提交失败')
    } finally {
      submitting.value = false
    }
  })
}

onMounted(() => {
  searchEquipment('')
})
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/repair/components/RepairReportDialog.vue
git commit -m "feat(repair): add RepairReportDialog component"
```

---

### Task 2: Create RepairExecuteDialog.vue

**Files:**
- Create: `frontend/src/views/repair/components/RepairExecuteDialog.vue`

- [ ] **Step 1: Implement the component**
- If status is `assigned`, show "开始维修" button calling `startRepair`.
- If status is `in_progress`, show form for solution, hours, spare parts, and next status.
- Call `repairOrderApi.updateRepair` on submit.

```vue
<template>
  <el-dialog
    :model-value="visible"
    title="维修执行"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <div v-if="order">
      <div v-if="order.status === 'assigned'" class="start-section">
        <p>工单状态：已派单</p>
        <el-button type="primary" @click="handleStart" :loading="processing">开始维修</el-button>
      </div>
      
      <el-form v-else-if="order.status === 'in_progress'" :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="解决方案" prop="solution">
          <el-input v-model="form.solution" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="实际工时" prop="actual_hours">
          <el-input-number v-model="form.actual_hours" :precision="1" :step="0.5" :min="0" />
        </el-form-item>
        <el-form-item label="使用备件" prop="spare_parts">
          <el-input v-model="form.spare_parts" placeholder="请输入备件名称及数量" />
        </el-form-item>
        <el-form-item label="后续状态" prop="next_status">
          <el-select v-model="form.next_status">
            <el-option label="待测试 (维修完成，需验证)" value="testing" />
            <el-option label="待审核 (已修复，直接提交)" value="confirmed" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
    <template #footer v-if="order?.status === 'in_progress'">
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="submit" :loading="processing">提交完成</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { repairOrderApi, type RepairOrder } from '@/api/repair'

const props = defineProps<{
  visible: boolean
  order: RepairOrder | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const processing = ref(false)

const form = reactive({
  solution: '',
  actual_hours: 1,
  spare_parts: '',
  next_status: 'testing'
})

const rules: FormRules = {
  solution: [{ required: true, message: '请输入解决方案', trigger: 'blur' }],
  next_status: [{ required: true, message: '请选择后续状态', trigger: 'change' }]
}

watch(() => props.order, (newOrder) => {
  if (newOrder) {
    form.solution = newOrder.solution || ''
    form.actual_hours = newOrder.actual_hours || 1
    form.spare_parts = newOrder.spare_parts || ''
  }
}, { immediate: true })

const handleStart = async () => {
  if (!props.order) return
  processing.value = true
  try {
    await repairOrderApi.startRepair(props.order.id, {})
    ElMessage.success('已开始维修')
    emit('success')
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    processing.value = false
  }
}

const submit = async () => {
  if (!formRef.value || !props.order) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    processing.value = true
    try {
      await repairOrderApi.updateRepair(props.order.id, {
        solution: form.solution,
        actual_hours: form.actual_hours,
        spare_parts: form.spare_parts,
        next_status: form.next_status
      })
      ElMessage.success('维修记录已更新')
      emit('success')
      emit('update:visible', false)
    } catch (error: any) {
      ElMessage.error(error.message || '更新失败')
    } finally {
      processing.value = false
    }
  })
}
</script>

<style scoped>
.start-section {
  text-align: center;
  padding: 20px 0;
}
</style>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/repair/components/RepairExecuteDialog.vue
git commit -m "feat(repair): add RepairExecuteDialog component"
```

---

### Task 3: Create RepairAuditDialog.vue

**Files:**
- Create: `frontend/src/views/repair/components/RepairAuditDialog.vue`

- [ ] **Step 1: Implement the component**
- Fields: `approved` (Pass/Reject), `comment`, `actual_hours`.
- Call `repairOrderApi.auditRepair` on submit.

```vue
<template>
  <el-dialog
    :model-value="visible"
    title="维修审核"
    width="500px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="审核结果" prop="approved">
        <el-radio-group v-model="form.approved">
          <el-radio :label="true">通过</el-radio>
          <el-radio :label="false">驳回</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="审核意见" prop="comment">
        <el-input v-model="form.comment" type="textarea" :rows="3" placeholder="请输入审核意见" />
      </el-form-item>
      <el-form-item label="调整工时" prop="actual_hours">
        <el-input-number v-model="form.actual_hours" :precision="1" :step="0.5" :min="0" />
        <span class="unit">小时</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { repairOrderApi, type RepairOrder } from '@/api/repair'

const props = defineProps<{
  visible: boolean
  order: RepairOrder | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  approved: true,
  comment: '',
  actual_hours: 0
})

const rules: FormRules = {
  comment: [{ required: true, message: '请输入审核意见', trigger: 'blur' }]
}

watch(() => props.order, (newOrder) => {
  if (newOrder) {
    form.actual_hours = newOrder.actual_hours || 0
  }
}, { immediate: true })

const submit = async () => {
  if (!formRef.value || !props.order) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await repairOrderApi.auditRepair(props.order.id, {
        approved: form.approved,
        comment: form.comment,
        actual_hours: form.actual_hours
      })
      ElMessage.success('审核完成')
      emit('success')
      emit('update:visible', false)
    } catch (error: any) {
      ElMessage.error(error.message || '审核提交失败')
    } finally {
      submitting.value = false
    }
  })
}
</script>

<style scoped>
.unit {
  margin-left: 10px;
  color: #909399;
}
</style>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/repair/components/RepairAuditDialog.vue
git commit -m "feat(repair): add RepairAuditDialog component"
```

---

### Task 4: Create RepairToKnowledgeDialog.vue

**Files:**
- Create: `frontend/src/views/repair/components/RepairToKnowledgeDialog.vue`

- [ ] **Step 1: Implement the component**
- Fields: `title`, `fault_phenomenon`, `cause_analysis`, `tags`.
- Call `convertFromRepair` from knowledge API.

```vue
<template>
  <el-dialog
    :model-value="visible"
    title="转入知识库"
    width="600px"
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="标题" prop="title">
        <el-input v-model="form.title" placeholder="输入知识库条目标题" />
      </el-form-item>
      <el-form-item label="故障现象" prop="fault_phenomenon">
        <el-input v-model="form.fault_phenomenon" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="原因分析" prop="cause_analysis">
        <el-input v-model="form.cause_analysis" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="标签">
        <el-select
          v-model="form.tags"
          multiple
          filterable
          allow-create
          default-first-option
          placeholder="请输入标签"
        >
          <el-option label="机械故障" value="机械故障" />
          <el-option label="电气故障" value="电气故障" />
          <el-option label="软件异常" value="软件异常" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="$emit('update:visible', false)">取消</el-button>
      <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { convertFromRepair } from '@/api/knowledge'
import type { RepairOrder } from '@/api/repair'

const props = defineProps<{
  visible: boolean
  order: RepairOrder | null
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  order_id: 0,
  title: '',
  fault_phenomenon: '',
  cause_analysis: '',
  tags: [] as string[]
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  fault_phenomenon: [{ required: true, message: '请输入故障现象', trigger: 'blur' }]
}

watch(() => props.order, (newOrder) => {
  if (newOrder) {
    form.order_id = newOrder.id
    form.title = `${newOrder.equipment_name} 故障处理记录 - ${newOrder.id}`
    form.fault_phenomenon = newOrder.fault_description
    form.cause_analysis = ''
    form.tags = []
  }
}, { immediate: true })

const submit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await convertFromRepair(form)
      ElMessage.success('成功转入知识库')
      emit('success')
      emit('update:visible', false)
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '转换失败')
    } finally {
      submitting.value = false
    }
  })
}
</script>
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/views/repair/components/RepairToKnowledgeDialog.vue
git commit -m "feat(repair): add RepairToKnowledgeDialog component"
```
