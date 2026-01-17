<template>
  <div class="repair-report-view">
    <div class="page-header">
      <el-button link @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h1>报修申请</h1>
    </div>

    <!-- 扫码选择设备 -->
    <div class="scan-section">
      <el-card @click="showScanDialog = true" class="scan-card" shadow="hover">
        <div class="scan-content">
          <el-icon size="48"><Grid /></el-icon>
          <div class="scan-text">
            <div class="scan-title">扫描设备二维码</div>
            <div class="scan-subtitle">点击开始扫码</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 设备信息 -->
    <div v-if="selectedEquipment" class="equipment-info">
      <el-card>
        <div class="selected-equipment">
          <el-icon color="#67c23a" size="24"><CircleCheck /></el-icon>
          <div class="equipment-detail">
            <div class="equipment-name">{{ selectedEquipment.name }}</div>
            <div class="equipment-code">{{ selectedEquipment.code }}</div>
          </div>
          <el-button size="small" @click="selectedEquipment = null; formData.equipment_id = 0">
            重新选择
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- 报修表单 -->
    <el-form :model="formData" :rules="formRules" ref="formRef" label-width="80px" class="report-form">
      <el-form-item label="故障描述" prop="fault_description">
        <el-input
          v-model="formData.fault_description"
          type="textarea"
          :rows="4"
          placeholder="请详细描述设备故障情况..."
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="故障代码">
        <el-input v-model="formData.fault_code" placeholder="可选，填写故障代码" />
      </el-form-item>

      <el-form-item label="优先级">
        <el-radio-group v-model="formData.priority">
          <el-radio :label="1">高</el-radio>
          <el-radio :label="2">中</el-radio>
          <el-radio :label="3">低</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="现场照片">
        <div class="photo-upload">
          <div v-for="(photo, index) in formData.photos" :key="index" class="photo-item">
            <img :src="photo" @click="previewPhoto(index)" />
            <el-icon class="photo-remove" @click="removePhoto(index)"><Close /></el-icon>
          </div>
          <el-upload
            :action="uploadAction"
            :headers="uploadHeaders"
            :on-success="onUploadSuccess"
            :show-file-list="false"
            accept="image/*"
            :before-upload="beforeUpload"
          >
            <div class="upload-placeholder">
              <el-icon><Plus /></el-icon>
            </div>
          </el-upload>
        </div>
      </el-form-item>
    </el-form>

    <div class="submit-section">
      <el-button
        type="primary"
        size="large"
        :disabled="!selectedEquipment"
        :loading="submitting"
        @click="submitReport"
        class="submit-btn"
      >
        提交报修
      </el-button>
    </div>

    <!-- 扫码对话框 -->
    <el-dialog v-model="showScanDialog" title="扫码选择设备" width="90%" :close-on-click-modal="false">
      <div class="scan-dialog">
        <el-input
          v-model="manualCode"
          placeholder="请输入或扫描设备二维码"
          @keyup.enter="confirmScan"
        >
          <template #append>
            <el-button @click="confirmScan">确定</el-button>
          </template>
        </el-input>

        <div class="manual-input-hint">
          或手动输入设备编号查询
          <el-input
            v-model="equipmentSearch"
            placeholder="输入设备编号搜索"
            @keyup.enter="searchEquipment"
            style="margin-top: 12px"
          >
            <template #append>
              <el-button @click="searchEquipment">搜索</el-button>
            </template>
          </el-input>
        </div>

        <!-- 搜索结果 -->
        <div v-if="searchResults.length > 0" class="search-results">
          <div
            v-for="equip in searchResults"
            :key="equip.id"
            class="equipment-option"
            @click="selectEquipment(equip)"
          >
            <div class="option-name">{{ equip.name }}</div>
            <div class="option-code">{{ equip.code }}</div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 照片预览 -->
    <el-image-viewer v-if="showImageViewer" :url-list="formData.photos" @close="showImageViewer = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, ElImageViewer } from 'element-plus'
import { ArrowLeft, Grid, CircleCheck, Plus, Close } from '@element-plus/icons-vue'
import { repairOrderApi, type CreateRepairRequest } from '@/api/repair'
import { equipmentApi, type Equipment } from '@/api/equipment'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const selectedEquipment = ref<Equipment | null>(null)
const showScanDialog = ref(false)
const showImageViewer = ref(false)
const submitting = ref(false)
const manualCode = ref('')
const equipmentSearch = ref('')
const searchResults = ref<Equipment[]>([])

const formRef = ref<FormInstance>()
const formData = reactive<CreateRepairRequest>({
  equipment_id: 0,
  fault_description: '',
  fault_code: '',
  photos: [],
  priority: 2
})

const formRules: FormRules = {
  fault_description: [{ required: true, message: '请描述故障情况', trigger: 'blur' }]
}

const uploadAction = computed(() => `${import.meta.env.VITE_API_BASE_URL}/upload`)
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`
}))

const confirmScan = async () => {
  if (!manualCode.value.trim()) {
    ElMessage.warning('请输入二维码内容')
    return
  }

  try {
    // 根据二维码查找设备
    const equipment = await equipmentApi.getByQRCode(manualCode.value.trim())
    selectEquipment(equipment)
    manualCode.value = ''
    showScanDialog.value = false
  } catch (error: any) {
    ElMessage.error('未找到对应设备')
  }
}

const searchEquipment = async () => {
  if (!equipmentSearch.value.trim()) {
    ElMessage.warning('请输入设备编号')
    return
  }

  try {
    // 搜索包含编号的设备
    const data = await equipmentApi.getList({ search: equipmentSearch.value })
    searchResults.value = data.slice(0, 5)
  } catch (error: any) {
    ElMessage.error('搜索失败')
  }
}

const selectEquipment = (equipment: Equipment) => {
  selectedEquipment.value = equipment
  formData.equipment_id = equipment.id
  showScanDialog.value = false
  searchResults.value = []
  equipmentSearch.value = ''
}

const beforeUpload = (file: File) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过5MB')
    return false
  }
  return true
}

const onUploadSuccess = (response: any) => {
  if (response.url) {
    formData.photos.push(response.url)
  }
}

const previewPhoto = (index: number) => {
  showImageViewer.value = true
}

const removePhoto = (index: number) => {
  formData.photos.splice(index, 1)
}

const submitReport = async () => {
  if (!formRef.value) return
  if (!selectedEquipment.value) {
    ElMessage.warning('请选择设备')
    return
  }

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      await repairOrderApi.createOrder(formData)
      ElMessage.success('报修申请已提交')
      router.back()
    } catch (error: any) {
      ElMessage.error(error.message || '提交失败')
    } finally {
      submitting.value = false
    }
  })
}
</script>

<style scoped>
.repair-report-view {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 80px;
}

.page-header {
  padding: 16px;
  background: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-header :deep(.el-icon) {
  font-size: 18px !important;
}


.page-header h1 {
  margin: 0;
  font-size: 18px;
}

.scan-section {
  padding: 16px;
}

.scan-card {
  cursor: pointer;
}

.scan-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.scan-text {
  flex: 1;
}

.scan-title {
  font-weight: 500;
  font-size: 16px;
}

.scan-subtitle {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.equipment-info {
  padding: 0 16px 16px;
}

.selected-equipment {
  display: flex;
  align-items: center;
  gap: 12px;
}

.equipment-detail {
  flex: 1;
}

.equipment-name {
  font-weight: 500;
}

.equipment-code {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.report-form {
  padding: 0 16px;
  background: #fff;
  margin-top: 16px;
}

.photo-upload {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.photo-item {
  position: relative;
  width: 80px;
  height: 80px;
}

.photo-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 4px;
}

.photo-remove {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #f56c6c;
  color: #fff;
  border-radius: 50%;
  padding: 2px;
  cursor: pointer;
}

.upload-placeholder {
  width: 80px;
  height: 80px;
  border: 2px dashed #dcdfe6;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #999;
}

.submit-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
}

.submit-btn {
  width: 100%;
}

.manual-input-hint {
  margin-top: 16px;
  font-size: 14px;
  color: #666;
}

.search-results {
  margin-top: 16px;
  max-height: 200px;
  overflow-y: auto;
}

.equipment-option {
  padding: 12px;
  background: #f5f5f5;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
}

.equipment-option:last-child {
  margin-bottom: 0;
}

.option-name {
  font-weight: 500;
}

.option-code {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>
