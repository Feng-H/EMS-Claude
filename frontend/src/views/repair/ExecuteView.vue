<template>
  <div class="repair-execute-view">
    <!-- Header -->
    <mobile-header title="维修执行" :show-back="true" />

    <div v-if="loading" class="loading-container">
      <van-loading size="24">加载中...</van-loading>
    </div>

    <div v-else-if="order" class="content">
      <!-- 工单信息 -->
      <van-cell-group inset class="order-info">
        <van-cell>
          <template #title>
            <div class="info-header">
              <span class="order-id">#{{ order.id }}</span>
              <van-tag :type="getStatusType(order.status)">
                {{ getStatusText(order.status) }}
              </van-tag>
            </div>
          </template>
        </van-cell>
        <van-cell title="设备名称" :value="order.equipment_name" />
        <van-cell title="设备编号" :value="order.equipment_code" />
        <van-cell title="故障描述">
          <template #value>
            <div class="fault-desc">{{ order.fault_description }}</div>
            <div v-if="order.fault_code" class="fault-code">
              故障代码: {{ order.fault_code }}
            </div>
          </template>
        </van-cell>
        <van-cell
          v-if="order.photos && order.photos.length > 0"
          title="故障照片"
        >
          <template #value>
            <div class="photos-grid">
              <van-image
                v-for="(photo, index) in order.photos"
                :key="index"
                :src="photo"
                fit="cover"
                class="photo-item"
                @click="previewFaultPhotos(index)"
              />
            </div>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 维修操作 -->
      <div v-if="canExecute" class="action-section">
        <van-cell-group inset>
          <van-cell title="维修操作" />
        </van-cell-group>

        <van-cell-group inset class="action-form">
          <!-- 开始维修 -->
          <div v-if="order.status === 'assigned'" class="start-repair">
            <van-button
              type="primary"
              block
              size="large"
              @click="startRepair"
            >
              <van-icon name="play" />
              开始维修
            </van-button>
            <div class="action-hint">点击开始维修将记录开始时间</div>
          </div>

          <!-- 维修进行中 -->
          <div v-else class="repair-form">
            <van-field
              v-model="updateForm.solution"
              label="解决方案"
              type="textarea"
              rows="4"
              placeholder="请描述维修过程、更换的部件等..."
              maxlength="500"
              show-word-limit
            />

            <!-- 备件选择器 -->
            <div class="spare-part-section">
              <spare-part-selector
                v-model="selectedSpareParts"
                :factory-id="order.factory_id"
                @change="onSparePartsChange"
              />
            </div>

            <van-field
              v-model.number="updateForm.actual_hours"
              label="实际工时"
              type="number"
              placeholder="请输入实际工时"
              input-align="right"
            >
              <template #button>
                <span>小时</span>
              </template>
            </van-field>

            <van-field label="维修后照片">
              <template #input>
                <van-uploader
                  v-model="repairPhotoFiles"
                  multiple
                  :max-count="6"
                  :after-read="onRepairPhotoRead"
                  :before-delete="onRepairPhotoDelete"
                />
              </template>
            </van-field>

            <div class="action-buttons">
              <van-button
                type="primary"
                block
                :loading="updating"
                :disabled="!updateForm.solution"
                @click="updateRepair('testing')"
              >
                完成维修，待确认
              </van-button>
              <van-button
                v-if="order.status === 'testing'"
                type="success"
                block
                :loading="updating"
                :disabled="!updateForm.solution"
                @click="updateRepair('confirmed')"
              >
                直接完成
              </van-button>
            </div>
          </div>
        </van-cell-group>
      </div>

      <!-- 确认操作 (报修人) -->
      <div v-if="order.status === 'testing' && isReporter" class="confirm-section">
        <van-cell-group inset>
          <van-cell title="维修确认" />
        </van-cell-group>

        <van-cell-group inset>
          <div class="confirm-actions">
            <van-button
              type="success"
              block
              size="large"
              @click="confirmRepair(true)"
            >
              <van-icon name="success" />
              确认修好
            </van-button>
            <van-button
              type="danger"
              block
              size="large"
              @click="showRejectDialog = true"
            >
              <van-icon name="cross" />
              未修好
            </van-button>
          </div>
          <div class="confirm-hint">请现场测试后确认</div>
        </van-cell-group>
      </div>

      <!-- 维修信息展示 -->
      <div v-if="order.solution" class="solution-section">
        <van-cell-group inset>
          <van-cell title="维修记录" />
        </van-cell-group>

        <van-cell-group inset>
          <van-cell title="解决方案" :value="order.solution" />
          <van-cell
            v-if="order.spare_parts"
            title="使用备件"
            :value="order.spare_parts"
          />
          <van-cell
            v-if="order.actual_hours"
            title="实际工时"
            :value="`${order.actual_hours} 小时`"
          />
          <van-cell title="开始时间" :value="formatDateTime(order.started_at)" />
          <van-cell
            v-if="order.completed_at"
            title="完成时间"
            :value="formatDateTime(order.completed_at)"
          />
        </van-cell-group>
      </div>
    </div>

    <!-- 拒绝确认对话框 -->
    <van-dialog
      v-model:show="showRejectDialog"
      title="未修好说明"
      show-cancel-button
      confirm-button-text="确定"
      confirm-button-color="#ee0a24"
      @confirm="confirmRepair(false)"
    >
      <van-field
        v-model="rejectForm.comment"
        type="textarea"
        placeholder="请说明设备还有什么问题..."
        rows="4"
        maxlength="200"
        show-word-limit
      />
    </van-dialog>

    <!-- 照片预览 -->
    <van-image-preview
      v-model:show="showFaultImageViewer"
      :images="order?.photos || []"
      :start-position="faultPreviewIndex"
    />

    <van-image-preview
      v-model:show="showRepairImageViewer"
      :images="repairPhotos"
      :start-position="repairPreviewIndex"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import {
  repairOrderApi,
  getStatusText,
  getStatusType,
  type RepairOrder,
  type UpdateRepairRequest,
  type ConfirmRepairRequest
} from '@/api/repair'
import { sparePartApi } from '@/api/sparepart'
import { useAuthStore } from '@/stores/auth'
import MobileHeader from '@/components/mobile/MobileHeader.vue'
import SparePartSelector from '@/components/sparepart/SparePartSelector.vue'

interface PhotoFile {
  url: string
  file?: File
}

interface SelectedSparePart {
  id: number
  code: string
  name: string
  specification?: string
  unit: string
  stock: number
  quantity: number
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(true)
const updating = ref(false)
const order = ref<RepairOrder | null>(null)
const showRejectDialog = ref(false)
const showFaultImageViewer = ref(false)
const showRepairImageViewer = ref(false)
const faultPreviewIndex = ref(0)
const repairPreviewIndex = ref(0)
const repairPhotos = ref<string[]>([])
const repairPhotoFiles = ref<PhotoFile[]>([])
const selectedSpareParts = ref<SelectedSparePart[]>([])

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
    showToast(error.message || '加载工单失败')
    router.back()
  } finally {
    loading.value = false
  }
}

const startRepair = async () => {
  if (!order.value) return

  try {
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
    showToast('已开始维修')
    loadOrder()
  } catch (error: any) {
    showToast(error.message || '操作失败')
  }
}

const updateRepair = async (nextStatus: string) => {
  if (!order.value || !updateForm.solution) {
    showToast('请填写解决方案')
    return
  }

  try {
    await showConfirmDialog({
      title: '确认提交',
      message: '确认提交维修记录？'
    })
  } catch {
    return
  }

  updating.value = true
  try {
    // 先记录备件消耗
    if (selectedSpareParts.value.length > 0) {
      for (const part of selectedSpareParts.value) {
        await sparePartApi.createConsumption({
          spare_part_id: part.id,
          quantity: part.quantity,
          order_id: order.value.id
        })
      }
    }

    // 然后更新维修工单
    await repairOrderApi.updateRepair(order.value.id, {
      ...updateForm,
      photos: repairPhotos.value,
      next_status
    })
    showToast('操作成功')
    loadOrder()
  } catch (error: any) {
    showToast(error.message || '操作失败')
  } finally {
    updating.value = false
  }
}

const confirmRepair = async (accepted: boolean) => {
  if (!order.value) return
  if (!accepted && !rejectForm.comment) {
    showToast('请说明问题')
    return
  }

  try {
    await repairOrderApi.confirmRepair(order.value.id, {
      accepted,
      comment: rejectForm.comment,
      photos: rejectForm.photos
    })
    showToast(accepted ? '已确认完成' : '已反馈问题')
    showRejectDialog.value = false
    loadOrder()
  } catch (error: any) {
    showToast(error.message || '操作失败')
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

const onRepairPhotoRead = (file: any) => {
  repairPhotos.value.push(file.content)
}

const onRepairPhotoDelete = (file: any, detail: any) => {
  const index = detail.index
  repairPhotos.value.splice(index, 1)
  return true
}

const onSparePartsChange = (parts: SelectedSparePart[]) => {
  // 更新文本格式的备件记录（用于提交）
  updateForm.spare_parts = parts
    .map(p => `${p.name}(${p.code}) x${p.quantity}${p.unit}`)
    .join(', ')
}

const previewFaultPhotos = (index: number) => {
  faultPreviewIndex.value = index
  showFaultImageViewer.value = true
}

const previewRepairPhotos = (index: number) => {
  repairPreviewIndex.value = index
  showRepairImageViewer.value = true
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped>
.repair-execute-view {
  min-height: 100vh;
  background: #f5f5f5;
  padding-top: 46px; /* NavBar 高度 */
  padding-bottom: 20px;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
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
  width: 100%;
}

.order-id {
  font-weight: bold;
  font-size: 16px;
}

.fault-desc {
  font-size: 15px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.fault-code {
  font-size: 13px;
  color: #ee0a24;
  margin-top: 8px;
}

.photos-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.photo-item {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 4px;
  overflow: hidden;
}

.action-section {
  margin-bottom: 16px;
}

.action-form {
  margin-top: -1px;
}

.start-repair {
  padding: 16px;
  text-align: center;
}

.action-hint {
  font-size: 12px;
  color: #999;
  margin-top: 12px;
}

.repair-form {
  padding: 16px;
}

.spare-part-section {
  padding: 0 16px;
  margin-bottom: 16px;
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 16px;
}

.confirm-section {
  margin-bottom: 16px;
}

.confirm-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
}

.confirm-hint {
  text-align: center;
  font-size: 12px;
  color: #999;
  padding: 0 16px 16px;
}

.solution-section {
  margin-bottom: 16px;
}
</style>
