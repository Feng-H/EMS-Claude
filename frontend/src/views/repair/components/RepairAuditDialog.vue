<template>
  <el-dialog
    :model-value="modelValue"
    title="维修审核"
    width="500px"
    @update:model-value="$emit('update:modelValue', $event)"
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
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="备件费用" prop="spare_part_cost">
            <el-input-number v-model="form.spare_part_cost" :precision="2" :min="0" style="width: 100%" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="人工费用" prop="labor_cost">
            <el-input-number v-model="form.labor_cost" :precision="2" :min="0" style="width: 100%" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="其他费用" prop="other_cost">
            <el-input-number v-model="form.other_cost" :precision="2" :min="0" style="width: 100%" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="停机损失" prop="downtime_loss">
            <el-input-number v-model="form.downtime_loss" :precision="2" :min="0" style="width: 100%" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <template #footer>
      <el-button @click="$emit('update:modelValue', false)">取消</el-button>
      <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { repairOrderApi, type RepairOrder } from '@/api/repair'

const props = defineProps<{
  modelValue: boolean
  order: RepairOrder | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = reactive({
  approved: true,
  comment: '',
  actual_hours: 0,
  spare_part_cost: 0,
  labor_cost: 0,
  other_cost: 0,
  downtime_loss: 0
})

const rules: FormRules = {
  comment: [{ required: true, message: '请输入审核意见', trigger: 'blur' }]
}

watch(() => props.order, (newOrder) => {
  if (newOrder) {
    form.actual_hours = newOrder.actual_hours || 0
    form.spare_part_cost = newOrder.spare_part_cost || 0
    form.labor_cost = newOrder.labor_cost || 0
    form.other_cost = newOrder.other_cost || 0
    form.downtime_loss = newOrder.downtime_loss || 0
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
        actual_hours: form.actual_hours,
        spare_part_cost: form.spare_part_cost,
        labor_cost: form.labor_cost,
        other_cost: form.other_cost,
        downtime_loss: form.downtime_loss
      })
      ElMessage.success('审核完成')
      emit('success')
      emit('update:modelValue', false)
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
