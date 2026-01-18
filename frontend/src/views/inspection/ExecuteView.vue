<template>
  <div class="inspection-execute-view">
    <!-- 扫码页面 -->
    <div v-if="step === 'scan'" class="scan-page">
      <!-- Header -->
      <mobile-header title="设备点检" :show-back="false" />

      <div class="scan-content">
        <div class="scan-instructions">
          <p>请扫描设备二维码开始点检</p>
        </div>

        <div class="scan-area">
          <div v-if="!scanning" class="scan-placeholder" @click="startScan">
            <van-icon name="scan" size="64" color="#409eff" />
            <p>点击开始扫码</p>
          </div>
          <video v-show="scanning" ref="videoRef" class="scan-video" autoplay playsinline></video>
          <canvas v-show="scanning" ref="canvasRef" class="scan-canvas"></canvas>
        </div>

        <van-button type="primary" plain block @click="manualInput">
          手动输入设备编号
        </van-button>

        <!-- 我今天的任务 -->
        <div class="my-tasks">
          <div class="section-title">我的今日任务</div>
          <van-cell-group inset>
            <van-cell
              v-for="task in myTasks"
              :key="task.id"
              class="task-cell"
              :title="task.equipment_name || task.equipment_code"
              :label="`模板：${task.template_name}`"
            >
              <template #right-icon>
                <van-tag :type="getTaskStatusType(task.status)">
                  {{ getTaskStatusText(task.status) }}
                </van-tag>
              </template>
            </van-cell>
          </van-cell-group>
          <div v-if="myTasks.length === 0" class="empty-tasks">
            <van-empty description="暂无任务" />
          </div>

          <!-- 任务操作按钮 -->
          <div v-for="task in myTasks" :key="`action-${task.id}`" class="task-action">
            <van-button
              v-if="task.status === 'pending' || task.status === 'in_progress'"
              type="primary"
              size="small"
              block
              @click="resumeTask(task)"
            >
              {{ task.status === 'in_progress' ? '继续' : '开始' }}
            </van-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 点检执行页面 -->
    <div v-else-if="step === 'inspect'" class="inspect-page">
      <!-- Header -->
      <mobile-header
        :title="equipment?.name || equipment?.code"
        :show-back="true"
        @click-left="goBack"
      />

      <div class="template-info">{{ templateName }}</div>

      <!-- 进度指示 -->
      <div class="progress-section">
        <van-progress :percentage="progress" :color="progress === 100 ? '#07c160' : '#1989fa'" />
        <div class="progress-text">
          <span>{{ currentItemIndex + 1 }} / {{ inspectionItems.length }}</span>
          <span>完成度 {{ progress }}%</span>
        </div>
      </div>

      <!-- 点检项目列表 -->
      <div class="items-list">
        <div
          v-for="(item, index) in inspectionItems"
          :key="item.id"
          class="item-card"
          :class="{ active: index === currentItemIndex, completed: itemsResult[item.id] }"
          @click="selectItem(index)"
        >
          <div class="item-number">{{ index + 1 }}</div>
          <div class="item-content">
            <div class="item-name">{{ item.name }}</div>
            <div v-if="item.method" class="item-method">{{ item.method }}</div>
          </div>
          <div class="item-status">
            <van-icon v-if="itemsResult[item.id] === 'OK'" name="success" color="#07c160" size="24" />
            <van-icon v-else-if="itemsResult[item.id] === 'NG'" name="fail" color="#ee0a24" size="24" />
            <van-icon v-else name="circle" color="#dcdee0" size="24" />
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <mobile-action-bar :actions="[
        {
          text: '提交点检结果',
          type: 'primary',
          disabled: !canSubmit,
          loading: submitting,
          onClick: submitInspection
        }
      ]" />
    </div>

    <!-- 点检项目详情弹窗 -->
    <van-popup
      v-model:show="showItemDrawer"
      position="bottom"
      :style="{ height: '70%' }"
      round
    >
      <div class="item-detail" v-if="currentItem">
        <div class="detail-header">
          <span class="item-index">{{ currentItemIndex + 1 }}.</span>
          <span class="item-title">{{ currentItem.name }}</span>
        </div>

        <van-cell-group inset class="detail-info">
          <van-cell title="检查方法" :value="currentItem.method || '无特殊要求'" />
          <van-cell title="判定标准" :value="currentItem.criteria || '无特殊要求'" />
        </van-cell-group>

        <van-divider>点检结果</van-divider>

        <van-radio-group v-model="currentItemResult" class="result-group">
          <van-radio name="OK" class="result-option">
            <div class="option-content ok">
              <van-icon name="success" />
              <span>正常 (OK)</span>
            </div>
          </van-radio>
          <van-radio name="NG" class="result-option">
            <div class="option-content ng">
              <van-icon name="fail" />
              <span>异常 (NG)</span>
            </div>
          </van-radio>
        </van-radio-group>

        <div v-if="currentItemResult === 'NG'" class="ng-section">
          <van-field
            v-model="currentItemRemark"
            type="textarea"
            label="异常描述"
            placeholder="请描述异常情况..."
            rows="3"
            maxlength="200"
            show-word-limit
          />
          <van-field label="照片上传">
            <template #input>
              <div class="photo-upload">
                <van-uploader
                  :after-read="onPhotoRead"
                  accept="image/*"
                  :max-size="5 * 1024 * 1024"
                  @oversize="onPhotoOversize"
                />
                <img v-if="currentItemPhoto" :src="currentItemPhoto" class="preview-photo" />
              </div>
            </template>
          </van-field>
        </div>

        <div class="drawer-actions">
          <van-button plain @click="showItemDrawer = false">取消</van-button>
          <van-button
            type="primary"
            :disabled="!currentItemResult"
            @click="confirmItemResult"
          >
            确认
          </van-button>
        </div>
      </div>
    </van-popup>

    <!-- 手动输入对话框 -->
    <van-dialog
      v-model:show="showManualDialog"
      title="手动输入设备编号"
      show-cancel-button
      @confirm="confirmManualInput"
    >
      <van-field
        v-model="manualGrid"
        placeholder="请输入设备二维码内容"
        @keyup.enter="confirmManualInput"
      />
    </van-dialog>

    <!-- 完成结果弹窗 -->
    <van-dialog
      v-model:show="showResultDialog"
      :title="ngCount > 0 ? '点检完成，发现异常' : '点检完成，全部正常'"
      :show-cancel-button="false"
      confirm-button-text="完成"
      @confirm="goBackToScan"
    >
      <div class="result-content">
        <van-icon
          :name="ngCount > 0 ? 'warning-o' : 'success'"
          :color="ngCount > 0 ? '#ff976a' : '#07c160'"
          size="64"
        />
        <div class="result-summary">
          <p>共检查 {{ totalCount }} 项</p>
          <p>正常 {{ okCount }} 项，异常 {{ ngCount }} 项</p>
        </div>

        <div v-if="ngCount > 0" class="ng-items">
          <div class="section-label">异常项目：</div>
          <van-tag v-for="itemId in ngItemIds" :key="itemId" type="danger" style="margin: 4px">
            {{ getItemName(itemId) }}
          </van-tag>
        </div>
      </div>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import {
  inspectionTaskApi,
  type InspectionTask,
  type InspectionItem,
  type StartInspectionResponse,
  type Equipment
} from '@/api/inspection'
import { equipmentApi } from '@/api/equipment'
import { useAuthStore } from '@/stores/auth'
import MobileHeader from '@/components/mobile/MobileHeader.vue'
import MobileActionBar from '@/components/mobile/MobileActionBar.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const step = ref<'scan' | 'inspect'>('scan')
const scanning = ref(false)
const submitting = ref(false)
const videoRef = ref<HTMLVideoElement>()
const canvasRef = ref<HTMLCanvasElement>()
const scanStreamRef = ref<MediaStream | null>(null)

const myTasks = ref<InspectionTask[]>([])
const equipment = ref<Equipment | null>(null)
const taskId = ref<number>(0)
const inspectionItems = ref<InspectionItem[]>([])
const itemsResult = ref<Record<number, 'OK' | 'NG'>>({})
const itemsRemark = ref<Record<number, string>>({})
const itemsPhoto = ref<Record<number, string>>({})

const currentItemIndex = ref(0)
const showItemDrawer = ref(false)
const currentItemResult = ref<'OK' | 'NG' | ''>('')
const currentItemRemark = ref('')
const currentItemPhoto = ref('')

const showManualDialog = ref(false)
const manualGrid = ref('')
const showResultDialog = ref(false)

const templateName = computed(() => {
  return myTasks.value.find(t => t.id === taskId.value)?.template_name || ''
})

const progress = computed(() => {
  if (inspectionItems.value.length === 0) return 0
  return Math.round((Object.keys(itemsResult.value).length / inspectionItems.value.length) * 100)
})

const currentItem = computed(() => {
  return inspectionItems.value[currentItemIndex.value]
})

const canSubmit = computed(() => {
  return Object.keys(itemsResult.value).length === inspectionItems.value.length && inspectionItems.value.length > 0
})

const totalCount = computed(() => inspectionItems.value.length)
const okCount = computed(() => Object.values(itemsResult.value).filter(r => r === 'OK').length)
const ngCount = computed(() => Object.values(itemsResult.value).filter(r => r === 'NG').length)
const ngItemIds = computed(() => {
  return Object.entries(itemsResult.value)
    .filter(([_, result]) => result === 'NG')
    .map(([itemId]) => parseInt(itemId))
})

const getTaskStatusType = (status: string) => {
  const map: Record<string, any> = {
    pending: 'primary',
    in_progress: 'warning',
    completed: 'success'
  }
  return map[status] || 'default'
}

const getTaskStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待执行',
    in_progress: '进行中',
    completed: '已完成'
  }
  return map[status] || status
}

const getItemName = (itemId: number) => {
  const item = inspectionItems.value.find(i => i.id === itemId)
  return item?.name || ''
}

// 加载我的任务
const loadMyTasks = async () => {
  try {
    const data = await inspectionTaskApi.getMyTasks()
    myTasks.value = data
  } catch (error: any) {
    showToast(error.message || '加载任务失败')
  }
}

// 启动扫码
const startScan = async () => {
  try {
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      showToast('当前浏览器不支持扫码，请使用手动输入')
      return
    }

    scanning.value = true
    const stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' }
    })
    scanStreamRef.value = stream

    if (videoRef.value) {
      videoRef.value.srcObject = stream
      videoRef.value.onloadedmetadata = () => {
        startScanLoop()
      }
    }
  } catch (error: any) {
    showToast(error.message || '无法访问摄像头')
    scanning.value = false
  }
}

// 停止扫码
const stopScan = () => {
  if (scanStreamRef.value) {
    scanStreamRef.value.getTracks().forEach(track => track.stop())
    scanStreamRef.value = null
  }
  scanning.value = false
}

// 扫码循环
let scanInterval: number | null = null
const startScanLoop = () => {
  if (scanInterval) clearInterval(scanInterval)
  scanInterval = window.setInterval(() => {
    // 实际项目中应集成 jsQR 或 html5-qrcode 库
    // 这里简化处理
  }, 500)
}

// 扫码成功处理
const onScanSuccess = async (qrCode: string) => {
  stopScan()
  try {
    // 先获取设备信息拿到ID
    const res = await equipmentApi.getByQRCode(qrCode)
    if (res && res.id) {
      await startInspection(qrCode, res.id)
    } else {
      await startInspection(qrCode)
    }
  } catch (error) {
    // 如果获取失败，仍然尝试启动，让后端处理错误
    console.warn('获取设备信息失败:', error)
    await startInspection(qrCode)
  }
}

// 开始点检
const startInspection = async (qrCode: string, equipmentId?: number) => {
  try {
    let latitude: number | undefined
    let longitude: number | undefined
    try {
      const position = await getCurrentPosition()
      latitude = position.coords.latitude
      longitude = position.coords.longitude
    } catch {
      // GPS获取失败，继续
    }

    const response: StartInspectionResponse = await inspectionTaskApi.start({
      equipment_id: equipmentId || 0,
      qr_code: qrCode,
      latitude,
      longitude,
      timestamp: Math.floor(Date.now() / 1000)
    })

    taskId.value = response.task_id
    equipment.value = response.equipment || null
    inspectionItems.value = response.items || []

    itemsResult.value = {}
    itemsRemark.value = {}
    itemsPhoto.value = {}
    currentItemIndex.value = 0

    step.value = 'inspect'
  } catch (error: any) {
    showToast(error.message || '启动点检失败')
  }
}

// 获取GPS位置
const getCurrentPosition = (): Promise<GeolocationPosition> => {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('设备不支持定位'))
      return
    }
    navigator.geolocation.getCurrentPosition(resolve, reject, {
      enableHighAccuracy: true,
      timeout: 10000,
      maximumAge: 0
    })
  })
}

// 手动输入
const manualInput = () => {
  showManualDialog.value = true
}

// 确认手动输入
const confirmManualInput = () => {
  if (!manualGrid.value.trim()) {
    showToast('请输入设备二维码')
    return
  }
  showManualDialog.value = false
  onScanSuccess(manualGrid.value.trim())
  manualGrid.value = ''
}

// 恢复任务
const resumeTask = async (task: InspectionTask) => {
  await startInspection(task.equipment_code || '', task.equipment_id)
}

// 选择项目
const selectItem = (index: number) => {
  currentItemIndex.value = index
  const item = inspectionItems.value[index]

  currentItemResult.value = itemsResult.value[item.id] || ''
  currentItemRemark.value = itemsRemark.value[item.id] || ''
  currentItemPhoto.value = itemsPhoto.value[item.id] || ''

  showItemDrawer.value = true
}

// 确认项目结果
const confirmItemResult = () => {
  if (!currentItem.value || !currentItemResult.value) return

  itemsResult.value[currentItem.value.id] = currentItemResult.value
  itemsRemark.value[currentItem.value.id] = currentItemRemark.value
  itemsPhoto.value[currentItem.value.id] = currentItemPhoto.value

  showItemDrawer.value = false

  // 自动跳到下一项
  if (currentItemIndex.value < inspectionItems.value.length - 1) {
    setTimeout(() => {
      selectItem(currentItemIndex.value + 1)
    }, 300)
  }
}

// 照片读取
const onPhotoRead = (file: any) => {
  currentItemPhoto.value = file.content
}

// 照片过大
const onPhotoOversize = () => {
  showToast('图片大小不能超过5MB')
}

// 提交点检
const submitInspection = async () => {
  if (!canSubmit.value) return

  submitting.value = true
  try {
    const records = Object.entries(itemsResult.value).map(([itemId, result]) => ({
      item_id: parseInt(itemId),
      result,
      remark: itemsRemark.value[itemId],
      photo_url: itemsPhoto.value[itemId]
    }))

    let latitude: number | undefined
    let longitude: number | undefined
    try {
      const position = await getCurrentPosition()
      latitude = position.coords.latitude
      longitude = position.coords.longitude
    } catch {
      // GPS获取失败
    }

    await inspectionTaskApi.complete({
      task_id: taskId.value,
      records,
      latitude,
      longitude
    })

    showResultDialog.value = true
  } catch (error: any) {
    showToast(error.message || '提交失败')
  } finally {
    submitting.value = false
  }
}

// 返回扫码页
const goBackToScan = () => {
  showResultDialog.value = false
  step.value = 'scan'
  loadMyTasks()
}

// 返回
const goBack = async () => {
  if (progress.value > 0) {
    try {
      await showDialog({
        title: '提示',
        message: '返回将丢失当前点检进度，确定返回？',
        showCancelButton: true
      })
      step.value = 'scan'
    } catch {
      // 用户取消
    }
  } else {
    step.value = 'scan'
  }
}

onMounted(async () => {
  // 按照顺序加载任务，避免并发竞争
  await loadMyTasks()

  const taskParam = route.params.taskId || route.query.taskId
  if (taskParam) {
    const id = Number(taskParam)
    const task = myTasks.value.find(t => t.id === id)
    if (task) {
      await resumeTask(task)
    } else {
      // 深度查找：如果今日任务列表中没找到（可能是从列表直接跳转的其他日期任务）
      try {
        const fullTask = await inspectionTaskApi.getTask(id)
        if (fullTask) {
          await resumeTask(fullTask)
        }
      } catch (error) {
        console.warn('获取任务详情失败:', error)
      }
    }
  }
})

onUnmounted(() => {
  stopScan()
  if (scanInterval) clearInterval(scanInterval)
})
</script>

<style scoped>
.inspection-execute-view {
  min-height: 100vh;
  background: #f5f5f5;
}

/* 扫码页面 */
.scan-page {
  min-height: 100vh;
  padding-top: 46px; /* NavBar 高度 */
}

.scan-content {
  padding: 16px;
  padding-bottom: 80px;
}

.scan-instructions {
  text-align: center;
  margin-bottom: 20px;
}

.scan-instructions p {
  color: #666;
  font-size: 16px;
  margin: 0;
}

.scan-area {
  position: relative;
  width: 100%;
  height: 300px;
  background: #000;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 20px;
}

.scan-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #fff;
  cursor: pointer;
}

.scan-placeholder p {
  margin-top: 12px;
  font-size: 14px;
}

.scan-video,
.scan-canvas {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scan-canvas {
  position: absolute;
  top: 0;
  left: 0;
}

.my-tasks {
  margin-top: 30px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 12px;
}

.task-cell {
  margin-bottom: 8px;
  border-radius: 8px;
}

.task-action {
  margin-top: -4px;
  margin-bottom: 12px;
}

.empty-tasks {
  margin-top: 20px;
}

/* 点检执行页面 */
.inspect-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-top: 46px; /* NavBar 高度 */
  padding-bottom: 80px;
}

.template-info {
  padding: 12px 16px;
  text-align: center;
  color: #666;
  font-size: 14px;
  background: #fff;
}

.progress-section {
  padding: 16px;
  background: #fff;
  margin-bottom: 12px;
}

.progress-text {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  font-size: 14px;
  color: #666;
}

.items-list {
  padding: 0 16px;
}

.item-card {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  padding: 12px;
  background: #fff;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.item-card.active {
  border: 1px solid #1989fa;
  background: #ecf9ff;
}

.item-card.completed {
  opacity: 0.7;
}

.item-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #1989fa;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  flex-shrink: 0;
}

.item-card.completed .item-number {
  background: #07c160;
}

.item-content {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-weight: 500;
  margin-bottom: 4px;
  font-size: 15px;
}

.item-method {
  font-size: 13px;
  color: #999;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 项目详情 */
.item-detail {
  padding: 20px 16px;
  max-height: 70vh;
  overflow-y: auto;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.item-index {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #1989fa;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  flex-shrink: 0;
}

.item-title {
  font-size: 18px;
  font-weight: 500;
  flex: 1;
}

.detail-info {
  margin-bottom: 16px;
}

.result-group {
  margin-bottom: 16px;
}

.result-option {
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid #ebedf0;
  border-radius: 8px;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
}

.option-content.ok {
  color: #07c160;
}

.option-content.ng {
  color: #ee0a24;
}

.ng-section {
  margin-top: 16px;
}

.photo-upload {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-photo {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 4px;
}

.drawer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

/* 结果弹窗 */
.result-content {
  text-align: center;
  padding: 20px;
}

.result-summary {
  margin: 16px 0;
}

.result-summary p {
  margin: 4px 0;
  font-size: 15px;
}

.ng-items {
  margin-top: 16px;
  text-align: left;
}

.section-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}
</style>
