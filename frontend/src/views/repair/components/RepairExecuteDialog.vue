<template>
  <el-dialog
    :model-value="modelValue"
    title="维修执行"
    width="600px"
    @update:model-value="$emit('update:modelValue', $event)"
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
          <span class="unit">小时</span>
        </el-form-item>
        <el-form-item label="备件消耗" prop="spare_parts">
          <el-input v-model="form.spare_parts" placeholder="请输入备件名称及数量" />
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
        <el-form-item label="后续状态" prop="next_status">
          <el-select v-model="form.next_status">
            <el-option label="待测试 (维修完成，需验证)" value="testing" />
            <el-option label="待审核 (已修复，直接提交)" value="confirmed" />
          </el-select>
        </el-form-item>
      </el-form>
      <div v-else class="status-tip">
        当前工单状态 ({{ order.status }}) 不支持执行操作。
      </div>
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
  modelValue: boolean
  order: RepairOrder | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const formRef = ref<FormInstance>()
const processing = ref(false)

const form = reactive({
  solution: '',
  actual_hours: 1,
  spare_parts: '',
  spare_part_cost: 0,
  labor_cost: 0,
  other_cost: 0,
  downtime_loss: 0,
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
    form.spare_part_cost = newOrder.spare_part_cost || 0
    form.labor_cost = newOrder.labor_cost || 0
    form.other_cost = newOrder.other_cost || 0
    form.downtime_loss = newOrder.downtime_loss || 0
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
        spare_part_cost: form.spare_part_cost,
        labor_cost: form.labor_cost,
        other_cost: form.other_cost,
        downtime_loss: form.downtime_loss,
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
.status-tip {
  text-align: center;
  padding: 20px 0;
  color: #909399;
}
</style>
