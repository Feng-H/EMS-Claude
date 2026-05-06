<template>
  <el-dialog
    v-model="visible"
    title="转入知识库"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
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
          style="width: 100%"
        >
          <el-option label="机械故障" value="机械故障" />
          <el-option label="电气故障" value="电气故障" />
          <el-option label="软件异常" value="软件异常" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { convertFromRepair } from '@/api/knowledge'
import type { RepairOrder } from '@/api/repair'

const props = defineProps<{
  modelValue: boolean
  order: RepairOrder | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

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

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(() => visible.value, (val) => {
  emit('update:modelValue', val)
})

watch(() => props.order, (order) => {
  if (order) {
    form.order_id = order.id
    form.title = `${order.equipment_name} 故障处理记录 - ${order.id}`
    form.fault_phenomenon = order.fault_description
    form.cause_analysis = ''
    form.tags = []
  }
}, { immediate: true })

const handleClose = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      await convertFromRepair(form)
      ElMessage.success('成功转入知识库')
      visible.value = false
      emit('success')
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '转换失败')
    } finally {
      submitting.value = false
    }
  })
}
</script>
