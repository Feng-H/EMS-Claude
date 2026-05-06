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
    form.comment = ''
    form.approved = true
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
