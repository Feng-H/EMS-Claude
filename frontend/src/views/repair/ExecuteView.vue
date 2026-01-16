<template>
  <div class="repair-execute-view">
    <div class="page-header">
      <el-button link @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h1>维修执行</h1>
    </div>

    <div v-if="loading" v-loading="true" class="loading-container"></div>

    <div v-else-if="order" class="content">
      <!-- 工单信息 -->
      <el-card class="order-info">
        <div class="info-header">
          <span class="order-id">#{{ order.id }}</span>
          <el-tag :type="getStatusType(order.status)" size="small">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>
        <div class="equipment-info">
          <div class="equipment-name">{{ order.equipment_name }}</div>
          <div class="equipment-code">{{ order.equipment_code }}</div>
        </div>
        <el-divider />
        <div class="fault-section">
          <div class="section-label">故障描述</div>
          <div class="fault-desc">{{ order.fault_description }}</div>
          <div v-if="order.fault_code" class="fault-code">
            故障代码: {{ order.fault_code }}
          </div>
        </div>
        <div v-if="order.photos && order.photos.length > 0" class="fault-photos">
          <div class="section-label">故障照片</div>
          <div class="photos-grid">
            <el-image
              v-for="(photo, index) in order.photos"
              :key="index"
              :src="photo"
              fit="cover"
              class="photo-item"
              :preview-src-list="order.photos"
            />
          </div>
        </div>
      </el-card>

      <!-- 维修操作 -->
      <el-card class="action-card" v-if="canExecute">
        <div class="section-title">维修操作</div>

        <!-- 开始维修 -->
        <div v-if="order.status === 'assigned'" class="action-section">
          <el-button type="primary" size="large" @click="startRepair" class="action-btn">
            <el-icon><VideoPlay /></el-icon>
            开始维修
          </el-button>
          <div class="action-hint">点击开始维修将记录开始时间</div>
        </div>

        <!-- 维修进行中 -->
        <div v-else class="repair-form">
          <el-form label-position="top">
            <el-form-item label="解决方案">
              <el-input
                v-model="updateForm.solution"
                type="textarea"
                :rows="4"
                placeholder="请描述维修过程、更换的部件等..."
              />
            </el-form-item>

            <el-form-item label="使用备件">
              <el-input
                v-model="updateForm.spare_parts"
                placeholder="记录使用的备件名称、数量等"
              />
            </el-form-item>

            <el-form-item label="实际工时(小时)">
              <el-input-number
                v-model="updateForm.actual_hours"
                :min="0"
                :max="999"
                :precision="1"
                :step="0.5"
              />
            </el-form-item>

            <el-form-item label="维修后照片">
              <div class="photo-upload">
                <div v-for="(photo, index) in repairPhotos" :key="index" class="photo-item">
                  <img :src="photo" @click="previewPhoto(index)" />
                  <el-icon class="photo-remove" @click="removeRepairPhoto(index)"><Close /></el-icon>
                </div>
                <el-upload
                  :action="uploadAction"
                  :headers="uploadHeaders"
                  :on-success="onRepairPhotoUpload"
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

          <div class="action-buttons">
            <el-button
              type="primary"
              size="large"
              @click="updateRepair('testing')"
              :loading="updating"
              :disabled="!updateForm.solution"
            >
              完成维修，待确认
            </el-button>
            <el-button
              v-if="order.status === 'testing'"
              type="success"
              size="large"
              @click="updateRepair('confirmed')"
              :loading="updating"
              :disabled="!updateForm.solution"
            >
              直接完成
            </el-button>
          </div>
        </div>
      </el-card>

      <!-- 确认操作 (报修人) -->
      <el-card class="action-card" v-if="order.status === 'testing' && isReporter">
        <div class="section-title">维修确认</div>
        <div class="confirm-actions">
          <el-button
            type="success"
            size="large"
            @click="confirmRepair(true)"
            class="confirm-btn accept"
          >
            <el-icon><Select /></el-icon>
            确认修好
          </el-button>
          <el-button
            type="danger"
            size="large"
            @click="showRejectDialog = true"
            class="confirm-btn reject"
          >
            <el-icon><CloseBold /></el-icon>
            未修好
          </el-button>
        </div>
        <div class="confirm-hint">请现场测试后确认</div>
      </el-card>

      <!-- 维修信息展示 -->
      <el-card class="solution-card" v-if="order.solution">
        <div class="section-title">维修记录</div>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="解决方案">
            {{ order.solution }}
          </el-descriptions-item>
          <el-descriptions-item label="使用备件" v-if="order.spare_parts">
            {{ order.spare_parts }}
          </el-descriptions-item>
          <el-descriptions-item label="实际工时" v-if="order.actual_hours">
            {{ order.actual_hours }} 小时
          </el-descriptions-item>
          <el-descriptions-item label="开始时间">
            {{ formatDateTime(order.started_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="完成时间" v-if="order.completed_at">
            {{ formatDateTime(order.completed_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>
    </div>

    <!-- 拒绝确认对话框 -->
    <el-dialog v-model="showRejectDialog" title="未修好说明" width="90%">
      <el-input
        v-model="rejectForm.comment"
        type="textarea"
        :rows="4"
        placeholder="请说明设备还有什么问题..."
        maxlength="200"
        show-word-limit
      />
      <template #footer>
        <el-button @click="showRejectDialog = false">取消</el-button>
        <el-button type="danger" @click="confirmRepair(false)">确定</el-button>
      </template>
    </el-dialog>

    <!-- 照片预览 -->
    <el-image-viewer v-if="showImageViewer" :url-list="repairPhotos" @close="showImageViewer = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElImageViewer } from 'element-plus'
import {
  ArrowLeft,
  VideoPlay,
  Plus,
  Close,
  Select,
  CloseBold
} from '@element-plus/icons-vue'
import {
  repairOrderApi,
  getStatusText,
  getStatusType,
  type RepairOrder,
  type UpdateRepairRequest,
  type ConfirmRepairRequest
} from '@/api/repair'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const updating = ref(false)
const order = ref<RepairOrder | null>(null)
const showRejectDialog = ref(false)
const showImageViewer = ref(false)
const repairPhotos = ref<string[]>([])

const updateForm = reactive<UpdateRepairRequest>({
  solution: '',
  spare_parts: '',
  actual_hours: 0,
  photos: []
})

const rejectForm = reactive({
  comment: '',
  photos: []
})

const uploadAction = computed(() => `${import.meta.env.VITE_API_BASE_URL}/upload`)
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`
}))

const canExecute = computed(() => {
  return order.value &&
    ['assigned', 'in_progress', 'testing'].includes(order.value.status)
})

const isReporter = computed(() => {
  return order.value && order.value.reporter_id === authStore.userId
})

const formatDateTime = (dateStr?: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const loadOrder = async () => {
  const orderId = Number(route.params.orderId)
  loading.value = true
  try {
    order.value = await repairOrderApi.getOrder(orderId)
    // 预填充已有数据
    if (order.value.solution) {
      updateForm.solution = order.value.solution
    }
    if (order.value.spare_parts) {
      updateForm.spare_parts = order.value.spare_parts
    }
    if (order.value.actual_hours) {
      updateForm.actual_hours = order.value.actual_hours
    }
    if (order.value.photos) {
      repairPhotos.value = [...order.value.photos]
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载工单失败')
    router.back()
  } finally {
    loading.value = false
  }
}

const startRepair = async () => {
  if (!order.value) return

  try {
    // 获取GPS位置
    let latitude: number | undefined
    let longitude: number | undefined
    try {
      const position = await getCurrentPosition()
      latitude = position.coords.latitude
      longitude = position.coords.longitude
    } catch {
      // GPS获取失败
    }

    await repairOrderApi.startRepair(order.value.id, { latitude, longitude })
    ElMessage.success('已开始维修')
    loadOrder()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

const updateRepair = async (nextStatus: string) => {
  if (!order.value || !updateForm.solution) {
    ElMessage.warning('请填写解决方案')
    return
  }

  updating.value = true
  try {
    await repairOrderApi.updateRepair(order.value.id, {
      ...updateForm,
      photos: repairPhotos.value,
      next_status
    })
    ElMessage.success('操作成功')
    loadOrder()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    updating.value = false
  }
}

const confirmRepair = async (accepted: boolean) => {
  if (!order.value) return
  if (!accepted && !rejectForm.comment) {
    ElMessage.warning('请说明问题')
    return
  }

  try {
    await repairOrderApi.confirmRepair(order.value.id, {
      accepted,
      comment: rejectForm.comment,
      photos: rejectForm.photos
    })
    ElMessage.success(accepted ? '已确认完成' : '已反馈问题')
    showRejectDialog.value = false
    loadOrder()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

const getCurrentPosition = (): Promise<GeolocationPosition> => {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('设备不支持定位'))
      return
    }
    navigator.geolocation.getCurrentPosition(resolve, reject, {
      enableHighAccuracy: true,
      timeout: 10000
    })
  })
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

const onRepairPhotoUpload = (response: any) => {
  if (response.url) {
    repairPhotos.value.push(response.url)
  }
}

const previewPhoto = (index: number) => {
  showImageViewer.value = true
}

const removeRepairPhoto = (index: number) => {
  repairPhotos.value.splice(index, 1)
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped>
.repair-execute-view {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 20px;
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

.loading-container {
  height: 300px;
}

.content {
  padding: 16px;
}

.order-info {
  margin-bottom: 16px;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.order-id {
  font-weight: bold;
  font-size: 18px;
}

.equipment-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.equipment-name {
  font-weight: 500;
  font-size: 16px;
}

.equipment-code {
  font-size: 12px;
  color: #999;
}

.fault-section {
  margin-top: 12px;
}

.section-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.fault-desc {
  font-size: 16px;
  line-height: 1.6;
}

.fault-code {
  font-size: 12px;
  color: #f56c6c;
  margin-top: 8px;
}

.fault-photos {
  margin-top: 12px;
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.photo-item {
  width: 100%;
  aspect-ratio: 1;
}

.action-card {
  margin-bottom: 16px;
}

.section-title {
  font-weight: 500;
  margin-bottom: 16px;
}

.action-section {
  text-align: center;
  padding: 16px 0;
}

.action-btn {
  width: 100%;
}

.action-hint {
  font-size: 12px;
  color: #999;
  margin-top: 12px;
}

.repair-form {
  padding: 16px 0;
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

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 16px;
}

.action-buttons .el-button {
  width: 100%;
}

.confirm-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px 0;
}

.confirm-btn {
  width: 100%;
}

.confirm-btn.accept {
  background: #67c23a;
  border-color: #67c23a;
}

.confirm-btn.reject {
  background: #f56c6c;
  border-color: #f56c6c;
}

.confirm-hint {
  text-align: center;
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

.solution-card {
  margin-bottom: 16px;
}
</style>
