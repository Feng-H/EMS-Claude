<template>
  <div class="repair-report-view">
    <!-- Header -->
    <mobile-header title="报修申请" :show-back="true" />

    <div class="content">
      <!-- 扫码选择设备 -->
      <div class="scan-section">
        <van-cell-group inset>
          <van-cell
            title="扫描设备二维码"
            label="点击开始扫码"
            is-link
            @click="showScanDialog = true"
          >
            <template #icon>
              <van-icon name="scan" size="24" color="#1989fa" />
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 设备信息 -->
      <div v-if="selectedEquipment" class="equipment-info">
        <van-cell-group inset>
          <van-cell>
            <template #title>
              <div class="selected-equipment">
                <van-icon name="success" color="#07c160" size="20" />
                <div class="equipment-detail">
                  <div class="equipment-name">{{ selectedEquipment.name }}</div>
                  <div class="equipment-code">{{ selectedEquipment.code }}</div>
                </div>
                <van-button size="small" type="primary" plain @click="clearEquipment">
                  重新选择
                </van-button>
              </div>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <!-- 报修表单 -->
      <van-cell-group inset class="form-section">
        <van-field
          v-model="formData.fault_description"
          label="故障描述"
          type="textarea"
          rows="4"
          placeholder="请详细描述设备故障情况..."
          maxlength="500"
          show-word-limit
          required
          :rules="[{ required: true, message: '请描述故障情况' }]"
        />

        <van-field
          v-model="formData.fault_code"
          label="故障代码"
          placeholder="可选，填写故障代码"
        />

        <van-field label="优先级">
          <template #input>
            <van-radio-group v-model="formData.priority" direction="horizontal">
              <van-radio :name="1">高</van-radio>
              <van-radio :name="2">中</van-radio>
              <van-radio :name="3">低</van-radio>
            </van-radio-group>
          </template>
        </van-field>

        <van-field label="现场照片">
          <template #input>
            <van-uploader
              v-model="photoFiles"
              multiple
              :max-count="6"
              :after-read="onPhotoRead"
              :before-delete="onPhotoDelete"
            />
          </template>
        </van-field>
      </van-cell-group>
    </div>

    <!-- 操作按钮 -->
    <mobile-action-bar :actions="[
      {
        text: '提交报修',
        type: 'primary',
        disabled: !selectedEquipment,
        loading: submitting,
        onClick: submitReport
      }
    ]" />

    <!-- 扫码对话框 -->
    <van-dialog
      v-model:show="showScanDialog"
      title="扫码选择设备"
      :show-cancel-button="true"
      confirm-button-text="确定"
      @confirm="confirmScan"
    >
      <div class="scan-dialog">
        <van-field
          v-model="manualCode"
          label="二维码"
          placeholder="请输入设备二维码"
          :border="false"
        />

        <div class="manual-input-hint">
          <div class="hint-text">或手动输入设备编号搜索</div>
          <van-field
            v-model="equipmentSearch"
            label="设备编号"
            placeholder="输入设备编号搜索"
            :border="false"
          >
            <template #button>
              <van-button size="small" type="primary" @click="searchEquipment">
                搜索
              </van-button>
            </template>
          </van-field>
        </div>

        <!-- 搜索结果 -->
        <div v-if="searchResults.length > 0" class="search-results">
          <van-cell
            v-for="equip in searchResults"
            :key="equip.id"
            :title="equip.name"
            :label="equip.code"
            is-link
            @click="selectEquipment(equip)"
          />
        </div>
      </div>
    </van-dialog>

    <!-- 照片预览 -->
    <van-image-preview
      v-model:show="showImageViewer"
      :images="formData.photos"
      :start-position="previewIndex"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { repairOrderApi, type CreateRepairRequest } from '@/api/repair'
import { equipmentApi, type Equipment } from '@/api/equipment'
import { useAuthStore } from '@/stores/auth'
import MobileHeader from '@/components/mobile/MobileHeader.vue'
import MobileActionBar from '@/components/mobile/MobileActionBar.vue'

interface PhotoFile {
  url: string
  file?: File
}

const router = useRouter()
const authStore = useAuthStore()

const selectedEquipment = ref<Equipment | null>(null)
const showScanDialog = ref(false)
const showImageViewer = ref(false)
const submitting = ref(false)
const manualCode = ref('')
const equipmentSearch = ref('')
const searchResults = ref<Equipment[]>([])
const photoFiles = ref<PhotoFile[]>([])
const previewIndex = ref(0)

const formData = reactive<CreateRepairRequest & { photos: string[] }>({
  equipment_id: 0,
  fault_description: '',
  fault_code: '',
  photos: [],
  priority: 2
})

const confirmScan = async () => {
  if (!manualCode.value.trim()) {
    showToast('请输入二维码内容')
    return
  }

  try {
    const equipment = await equipmentApi.getByQRCode(manualCode.value.trim())
    selectEquipment(equipment)
    manualCode.value = ''
    showScanDialog.value = false
  } catch (error: any) {
    showToast('未找到对应设备')
  }
}

const searchEquipment = async () => {
  if (!equipmentSearch.value.trim()) {
    showToast('请输入设备编号')
    return
  }

  try {
    const data = await equipmentApi.getList({ search: equipmentSearch.value })
    searchResults.value = data.slice(0, 5)
  } catch (error: any) {
    showToast('搜索失败')
  }
}

const selectEquipment = (equipment: Equipment) => {
  selectedEquipment.value = equipment
  formData.equipment_id = equipment.id
  showScanDialog.value = false
  searchResults.value = []
  equipmentSearch.value = ''
}

const clearEquipment = () => {
  selectedEquipment.value = null
  formData.equipment_id = 0
}

const onPhotoRead = (file: any, detail: any) => {
  // 添加到photos数组
  formData.photos.push(file.content)

  // 验证文件
  if (file.file) {
    const isImage = file.file.type.startsWith('image/')
    if (!isImage) {
      showToast('只能上传图片')
      return
    }
    const isLt5M = file.file.size / 1024 / 1024 < 5
    if (!isLt5M) {
      showToast('图片大小不能超过5MB')
      return
    }
  }
}

const onPhotoDelete = (file: any, detail: any) => {
  const index = detail.index
  formData.photos.splice(index, 1)
  return true
}

const previewPhoto = (index: number) => {
  previewIndex.value = index
  showImageViewer.value = true
}

const submitReport = async () => {
  if (!selectedEquipment.value) {
    showToast('请选择设备')
    return
  }

  if (!formData.fault_description.trim()) {
    showToast('请描述故障情况')
    return
  }

  try {
    await showConfirmDialog({
      title: '确认提交',
      message: '确认提交报修申请？'
    })
  } catch {
    return
  }

  submitting.value = true
  try {
    await repairOrderApi.createOrder({
      equipment_id: formData.equipment_id,
      fault_description: formData.fault_description,
      fault_code: formData.fault_code || undefined,
      priority: formData.priority
    })
    showToast('报修申请已提交')
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (error: any) {
    showToast(error.message || '提交失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.repair-report-view {
  min-height: 100vh;
  background: #f5f5f5;
  padding-top: 46px; /* NavBar 高度 */
  padding-bottom: 80px;
}

.content {
  padding: 16px;
}

.scan-section {
  margin-bottom: 16px;
}

.equipment-info {
  margin-bottom: 16px;
}

.selected-equipment {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.equipment-detail {
  flex: 1;
}

.equipment-name {
  font-weight: 500;
  font-size: 15px;
}

.equipment-code {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.form-section {
  margin-bottom: 16px;
}

.scan-dialog {
  padding: 16px;
}

.manual-input-hint {
  margin-top: 16px;
}

.hint-text {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.search-results {
  margin-top: 16px;
  max-height: 200px;
  overflow-y: auto;
}
</style>
