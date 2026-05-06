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
  loadingEquipment.value = true
  try {
    const res = await equipmentApi.getList({ name: query, page: 1, page_size: 20 })
    equipmentList.value = res.data.items
  } catch (error) {
    console.error('搜索设备失败', error)
  } finally {
    loadingEquipment.value = false
  }
}

const handlePhotoChange = (file: any) => {
  // Simple mock: store file object or base64
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
