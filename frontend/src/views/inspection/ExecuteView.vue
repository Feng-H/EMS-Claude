<template>
  <div class="execute-view">
    <!-- 扫码页面 -->
    <div v-if="step === 'scan'" class="scan-page">
      <div class="scan-header">
        <h1>设备点检</h1>
        <p>请扫描设备二维码开始点检</p>
      </div>

      <div class="scan-area">
        <div v-if="!scanning" class="scan-placeholder" @click="startScan">
          <el-icon size="64"><Grid /></el-icon>
          <p>点击开始扫码</p>
        </div>
        <video v-show="scanning" ref="videoRef" class="scan-video" autoplay playsinline></video>
        <canvas v-show="scanning" ref="canvasRef" class="scan-canvas"></canvas>
      </div>

      <div class="scan-actions">
        <el-button @click="manualInput" type="primary" plain>
          手动输入设备编号
        </el-button>
      </div>

      <!-- 我今天的任务 -->
      <div class="my-tasks">
        <div class="section-title">我的今日任务</div>
        <el-card v-for="task in myTasks" :key="task.id" class="task-card" shadow="hover">
          <div class="task-info">
            <span class="task-equipment">{{ task.equipment_name || task.equipment_code }}</span>
            <el-tag :type="getTaskStatusType(task.status)" size="small">
              {{ getTaskStatusText(task.status) }}
            </el-tag>
          </div>
          <el-button
            v-if="task.status === 'pending' || task.status === 'in_progress'"
            type="primary"
            size="small"
            @click="resumeTask(task)"
          >
            {{ task.status === 'in_progress' ? '继续' : '开始' }}
          </el-button>
        </el-card>
        <el-empty v-if="myTasks.length === 0" description="暂无任务" :image-size="60" />
      </div>
    </div>

    <!-- 点检执行页面 -->
    <div v-else-if="step === 'inspect'" class="inspect-page">
      <div class="inspect-header">
        <el-button link @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h2>{{ equipment?.name || equipment?.code }}</h2>
        <p class="template-name">{{ templateName }}</p>
      </div>

      <!-- 进度指示 -->
      <div class="progress-indicator">
        <el-progress
          :percentage="progress"
          :color="progress === 100 ? '#67c23a' : '#409eff'"
        />
        <span class="progress-text">{{ currentItemIndex + 1 }} / {{ inspectionItems.length }}</span>
      </div>

      <!-- 点检项目列表 -->
      <div class="items-container">
        <el-card
          v-for="(item, index) in inspectionItems"
          :key="item.id"
          :class="['item-card', { active: index === currentItemIndex, completed: itemsResult[item.id] }]"
          shadow="hover"
          @click="selectItem(index)"
        >
          <div class="item-number">{{ index + 1 }}</div>
          <div class="item-content">
            <div class="item-name">{{ item.name }}</div>
            <div class="item-method" v-if="item.method">{{ item.method }}</div>
          </div>
          <div class="item-status">
            <el-icon v-if="itemsResult[item.id] === 'OK'" color="#67c23a" size="24"><CircleCheck /></el-icon>
            <el-icon v-else-if="itemsResult[item.id] === 'NG'" color="#f56c6c" size="24"><CircleClose /></el-icon>
            <el-icon v-else color="#dcdfe6" size="24"><CircleCheck /></el-icon>
          </div>
        </el-card>
      </div>

      <!-- 操作按钮 -->
      <div class="inspect-actions">
        <el-button
          @click="submitInspection"
          type="primary"
          size="large"
          :disabled="!canSubmit"
          :loading="submitting"
        >
          提交点检结果
        </el-button>
      </div>
    </div>

    <!-- 点检项目详情弹窗 -->
    <el-drawer
      v-model="showItemDrawer"
      :title="`${currentItemIndex + 1}. ${currentItem?.name}`"
      direction="btt"
      size="70%"
    >
      <div class="item-detail" v-if="currentItem">
        <div class="detail-section">
          <div class="section-label">检查方法</div>
          <div class="section-value">{{ currentItem.method || '无特殊要求' }}</div>
        </div>
        <div class="detail-section">
          <div class="section-label">判定标准</div>
          <div class="section-value">{{ currentItem.criteria || '无特殊要求' }}</div>
        </div>

        <el-divider>点检结果</el-divider>

        <el-radio-group v-model="currentItemResult" class="result-group" size="large">
          <el-radio-button label="OK">
            <el-icon color="#67c23a"><CircleCheck /></el-icon>
            正常
          </el-radio-button>
          <el-radio-button label="NG">
            <el-icon color="#f56c6c"><CircleClose /></el-icon>
            异常
          </el-radio-button>
        </el-radio-group>

        <div v-if="currentItemResult === 'NG'" class="ng-section">
          <el-input
            v-model="currentItemRemark"
            type="textarea"
            placeholder="请描述异常情况..."
            :rows="3"
            maxlength="200"
            show-word-limit
          />
          <div class="photo-upload">
            <el-upload
              :action="uploadAction"
              :headers="uploadHeaders"
              :on-success="onUploadSuccess"
              :show-file-list="false"
              accept="image/*"
              :before-upload="beforeUpload"
            >
              <el-button type="primary" plain>
                <el-icon><Camera /></el-icon>
                拍照上传
              </el-button>
            </el-upload>
            <img v-if="currentItemPhoto" :src="currentItemPhoto" class="preview-photo" />
          </div>
        </div>

        <div class="drawer-actions">
          <el-button @click="showItemDrawer = false">取消</el-button>
          <el-button type="primary" @click="confirmItemResult" :disabled="!currentItemResult">
            确认
          </el-button>
        </div>
      </div>
    </el-drawer>

    <!-- 手动输入对话框 -->
    <el-dialog v-model="showManualDialog" title="手动输入设备编号" width="90%">
      <el-input
        v-model="manualGrid"
        placeholder="请输入设备二维码内容"
        @keyup.enter="confirmManualInput"
      />
      <template #footer>
        <el-button @click="showManualDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmManualInput">确定</el-button>
      </template>
    </el-dialog>

    <!-- 完成结果弹窗 -->
    <el-dialog v-model="showResultDialog" title="点检完成" width="90%" :close-on-click-modal="false">
      <div class="result-content">
        <el-result :icon="ngCount > 0 ? 'warning' : 'success'">
          <template #title>
            <h3>{{ ngCount > 0 ? '点检完成，发现异常' : '点检完成，全部正常' }}</h3>
          </template>
          <template #sub-title>
            <p>共检查 {{ totalCount }} 项，正常 {{ okCount }} 项，异常 {{ ngCount }} 项</p>
          </template>
        </el-result>

        <div v-if="ngCount > 0" class="ng-items">
          <div class="section-label">异常项目：</div>
          <el-tag v-for="itemId in ngItemIds" :key="itemId" type="danger" style="margin: 4px">
            {{ getItemName(itemId) }}
          </el-tag>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="goBackToScan">完成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Grid,
  ArrowLeft,
  CircleCheck,
  CircleClose,
  Camera
} from '@element-plus/icons-vue'
import {
  inspectionTaskApi,
  type InspectionTask,
  type InspectionItem,
  type StartInspectionResponse,
  type Equipment
} from '@/api/inspection'
import { useAuthStore } from '@/stores/auth'

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

// 上传配置
const uploadAction = computed(() => `${import.meta.env.VITE_API_BASE_URL}/upload`)
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`
}))

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
    pending: 'info',
    in_progress: 'warning',
    completed: 'success'
  }
  return map[status] || 'info'
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
    ElMessage.error(error.message || '加载任务失败')
  }
}

// 启动扫码
const startScan = async () => {
  try {
    // 检查是否支持摄像头
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      ElMessage.warning('当前浏览器不支持扫码，请使用手动输入')
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
    ElMessage.error(error.message || '无法访问摄像头')
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

// 扫码循环（使用简单的扫码检测）
let scanInterval: number | null = null
const startScanLoop = () => {
  if (scanInterval) clearInterval(scanInterval)
  scanInterval = window.setInterval(() => {
    // 实际项目中应集成 jsQR 或 html5-qrcode 库
    // 这里简化处理，等待用户手动输入或点击任务
  }, 500)
}

// 扫码成功处理
const onScanSuccess = async (qrCode: string) => {
  stopScan()
  await startInspection(qrCode)
}

// 开始点检
const startInspection = async (qrCode: string, equipmentId?: number) => {
  try {
    // 获取GPS位置
    let latitude: number | undefined
    let longitude: number | undefined
    try {
      const position = await getCurrentPosition()
      latitude = position.coords.latitude
      longitude = position.coords.longitude
    } catch {
      // GPS获取失败，继续但可能触发警告
    }

    const response: StartInspectionResponse = await inspectionTaskApi.start({
      equipment_id: equipmentId || 0,
      qr_code: qrCode,
      latitude,
      longitude,
      timestamp: Date.now()
    })

    taskId.value = response.task_id
    equipment.value = response.equipment || null
    inspectionItems.value = response.items || []

    // 重置状态
    itemsResult.value = {}
    itemsRemark.value = {}
    itemsPhoto.value = {}
    currentItemIndex.value = 0

    step.value = 'inspect'
  } catch (error: any) {
    ElMessage.error(error.message || '启动点检失败')
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
    ElMessage.warning('请输入设备二维码')
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

  // 加载已保存的结果
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

// 上传前
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

// 上传成功
const onUploadSuccess = (response: any) => {
  if (response.url) {
    currentItemPhoto.value = response.url
  }
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

    await inspectionTaskApi.complete({
      task_id: taskId.value,
      records,
      latitude,
      longitude
    })

    showResultDialog.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '提交失败')
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
const goBack = () => {
  if (progress.value > 0) {
    ElMessageBox.confirm('返回将丢失当前点检进度，确定返回？', '提示', {
      type: 'warning'
    }).then(() => {
      step.value = 'scan'
    }).catch(() => {})
  } else {
    step.value = 'scan'
  }
}

onMounted(() => {
  loadMyTasks()

  // 检查是否有taskId参数（从任务列表进入）
  const taskParam = route.params.taskId
  if (taskParam) {
    const task = myTasks.value.find(t => t.id === Number(taskParam))
    if (task) {
      resumeTask(task)
    }
  }
})

onUnmounted(() => {
  stopScan()
  if (scanInterval) clearInterval(scanInterval)
})
</script>

<style scoped>
.execute-view {
  min-height: 100vh;
  background: #f5f5f5;
}

/* 扫码页面 */
.scan-page {
  padding: 20px;
}

.scan-header {
  text-align: center;
  margin-bottom: 20px;
}

.scan-header h1 {
  font-size: 24px;
  margin: 0 0 8px;
}

.scan-header p {
  color: #666;
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

.scan-actions {
  display: flex;
  justify-content: center;
  margin-bottom: 30px;
}

.my-tasks {
  margin-top: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 12px;
}

.task-card {
  margin-bottom: 12px;
}

.task-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.task-equipment {
  font-weight: 500;
}

/* 点检执行页面 */
.inspect-page {
  min-height: 100vh;
  background: #fff;
}

.inspect-header {
  padding: 16px;
  background: #409eff;
  color: #fff;
}

.inspect-header :deep(.el-icon) {
  font-size: 18px !important;
}

.inspect-header h2 {
  margin: 8px 0 4px;
  font-size: 20px;
}

.template-name {
  margin: 0;
  opacity: 0.8;
  font-size: 14px;
}

.progress-indicator {
  padding: 16px;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  gap: 16px;
}

.progress-text {
  white-space: nowrap;
  font-weight: 500;
}

.items-container {
  padding: 16px;
}

.item-card {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.3s;
}

.item-card.active {
  border-color: #409eff;
  background: #ecf5ff;
}

.item-card.completed {
  opacity: 0.7;
}

.item-number {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #409eff;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
}

.item-card.completed .item-number {
  background: #67c23a;
}

.item-content {
  flex: 1;
}

.item-name {
  font-weight: 500;
  margin-bottom: 4px;
}

.item-method {
  font-size: 12px;
  color: #999;
}

.inspect-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
}

.inspect-actions .el-button {
  width: 100%;
}

/* 项目详情 */
.item-detail {
  padding: 16px;
}

.detail-section {
  margin-bottom: 16px;
}

.section-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.section-value {
  font-size: 16px;
  color: #333;
}

.result-group {
  display: flex;
  width: 100%;
  margin-bottom: 16px;
}

.result-group .el-radio-button {
  flex: 1;
}

.result-group .el-radio-button__inner {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.ng-section {
  margin-top: 16px;
}

.photo-upload {
  margin-top: 12px;
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
}

.ng-items {
  margin-top: 16px;
  text-align: left;
}
</style>
