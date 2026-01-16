<template>
  <div class="equipment-detail">
    <el-page-header @back="() => $router.back()">
      <template #content>
        <span class="page-title">设备详情</span>
      </template>
    </el-page-header>

    <el-card v-loading="loading" style="margin-top: 16px">
      <el-descriptions v-if="equipment" :column="2" border>
        <el-descriptions-item label="设备编号">{{ equipment.code }}</el-descriptions-item>
        <el-descriptions-item label="设备名称">{{ equipment.name }}</el-descriptions-item>
        <el-descriptions-item label="设备类型">{{ equipment.type_name }}</el-descriptions-item>
        <el-descriptions-item label="所属工厂">{{ equipment.factory_name }}</el-descriptions-item>
        <el-descriptions-item label="所属车间">{{ equipment.workshop_name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(equipment.status)">
            {{ getStatusText(equipment.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="二维码">{{ equipment.qr_code }}</el-descriptions-item>
        <el-descriptions-item label="采购日期">
          {{ equipment.purchase_date || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="专属维修工">
          {{ equipment.dedicated_maintenance_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">
          {{ formatDate(equipment.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="技术参数" :span="2">
          {{ equipment.spec || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-divider />

      <div class="detail-actions">
        <el-button v-if="canEdit" type="primary" @click="handleEdit">
          <el-icon><Edit /></el-icon>
          编辑设备
        </el-button>
        <el-button @click="handleQRCode">
          <el-icon><Grid /></el-icon>
          查看二维码
        </el-button>
        <el-button @click="$router.push('/inspection')">
          <el-icon><CircleCheck /></el-icon>
          点检记录
        </el-button>
        <el-button @click="$router.push('/repair')">
          <el-icon><Tools /></el-icon>
          维修记录
        </el-button>
      </div>
    </el-card>

    <!-- QR Code Dialog -->
    <el-dialog v-model="showQRDialog" title="设备二维码" width="400px">
      <div class="qr-content">
        <div class="qr-placeholder">
          <el-icon :size="150"><Grid /></el-icon>
        </div>
        <p class="qr-code-text">{{ equipment?.qr_code }}</p>
        <p class="qr-name">{{ equipment?.name }}</p>
      </div>
      <template #footer>
        <el-button @click="showQRDialog = false">关闭</el-button>
        <el-button type="primary" @click="handleDownloadQR">下载二维码</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { equipmentApi, type Equipment } from '@/api/equipment'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const canEdit = computed(() => authStore.hasRole('admin', 'engineer'))

const loading = ref(false)
const equipment = ref<Equipment | null>(null)
const showQRDialog = ref(false)

function getStatusType(status: string) {
  const map: Record<string, any> = {
    running: 'success',
    stopped: 'info',
    maintenance: 'warning',
    scrapped: 'danger',
  }
  return map[status] || ''
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停机',
    maintenance: '维修中',
    scrapped: '已报废',
  }
  return map[status] || status
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const id = Number(route.params.id)
    equipment.value = await equipmentApi.getById(id)
  } catch (error) {
    ElMessage.error('获取设备详情失败')
    router.back()
  } finally {
    loading.value = false
  }
}

function handleEdit() {
  router.push('/equipment')
}

function handleQRCode() {
  showQRDialog.value = true
}

function handleDownloadQR() {
  ElMessage.info('二维码下载功能开发中')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.equipment-detail {
  height: 100%;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
}

.detail-actions {
  display: flex;
  gap: 12px;
}

.qr-content {
  text-align: center;
}

.qr-placeholder {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
}

.qr-code-text {
  margin-top: 16px;
  font-size: 16px;
  font-weight: bold;
  color: #409eff;
}

.qr-name {
  margin-top: 8px;
  color: #606266;
}
</style>
