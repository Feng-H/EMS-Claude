<template>
  <div class="h5-equipment-detail">
    <mobile-header title="设备详情" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <div v-if="loading && !refreshing" class="loading-state">
        <van-loading vertical>加载中...</van-loading>
      </div>
      
      <div v-else-if="equipment" class="detail-content">
        <!-- 基础信息卡片 -->
        <van-cell-group inset title="基础信息">
          <van-cell title="设备名称" :value="equipment.name" />
          <van-cell title="设备编号" :value="equipment.code" />
          <van-cell title="设备类型" :value="equipment.type_name" />
          <van-cell title="状态">
            <template #value>
              <van-tag :type="getStatusType(equipment.status)">{{ getStatusText(equipment.status) }}</van-tag>
            </template>
          </van-cell>
        </van-cell-group>

        <!-- 部署位置卡片 -->
        <van-cell-group inset title="位置信息">
          <van-cell title="所属工厂" :value="equipment.factory_name" />
          <van-cell title="所属车间" :value="equipment.workshop_name" />
        </van-cell-group>

        <!-- 技术与维护卡片 -->
        <van-cell-group inset title="技术与维护">
          <van-cell title="采购日期" :value="equipment.purchase_date || '-'" />
          <van-cell title="专属维修工" :value="equipment.dedicated_maintenance_name || '-'" />
          <van-cell title="技术参数" :label="equipment.spec || '暂无参数'" />
        </van-cell-group>

        <!-- 统计信息 -->
        <div class="stats-grid">
          <div class="stat-box">
            <span class="stat-label">点检次数</span>
            <span class="stat-value">24</span>
          </div>
          <div class="stat-box">
            <span class="stat-label">维修次数</span>
            <span class="stat-value">3</span>
          </div>
          <div class="stat-box">
            <span class="stat-label">运行天数</span>
            <span class="stat-value">128</span>
          </div>
        </div>

        <!-- 快捷操作按钮 -->
        <div class="action-buttons">
          <van-button type="primary" block round icon="setting-o" @click="reportRepair">故障报修</van-button>
          <van-button type="success" plain block round icon="scan" @click="showQRCode">查看二维码</van-button>
        </div>
      </div>
    </van-pull-refresh>

    <!-- 二维码弹窗 -->
    <van-dialog v-model:show="qrDialog.show" :title="equipment?.name" theme="round-button">
      <div class="qr-popup-content">
        <div class="qr-image-placeholder">
          <van-icon name="qr" size="180" color="#333" />
        </div>
        <p class="qr-code-val">{{ equipment?.code }}</p>
        <p class="qr-tip">现场扫描此码可快速报修或点检</p>
      </div>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast } from 'vant'
import { equipmentApi, type Equipment } from '@/api/equipment'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const refreshing = ref(false)
const equipment = ref<Equipment | null>(null)

const qrDialog = reactive({
  show: false
})

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    running: 'success',
    stopped: 'warning',
    maintenance: 'danger',
    scrapped: 'default'
  }
  return map[status] || 'default'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停机',
    maintenance: '维修中',
    scrapped: '已报废'
  }
  return map[status] || status
}

const loadData = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await equipmentApi.getById(id)
    equipment.value = res.data
  } catch (err) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const onRefresh = () => {
  loadData()
}

const reportRepair = () => {
  if (equipment.value) {
    router.push(`/h5/repair/report?equipmentId=${equipment.value.id}`)
  }
}

const showQRCode = () => {
  qrDialog.show = true
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.h5-equipment-detail {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 30px;
}

.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
}

.detail-content {
  padding-top: 10px;
}

.detail-content :deep(.van-cell-group__title) {
  padding: 16px 20px 8px;
  color: var(--color-text-tertiary);
  font-size: 13px;
}

.detail-content :deep(.van-cell) {
  background: var(--color-bg-card);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin: 20px;
}

.stat-box {
  background: var(--color-bg-card);
  padding: 12px;
  border-radius: var(--radius-md);
  text-align: center;
  border: 1px solid var(--color-border);
}

.stat-label {
  display: block;
  font-size: 11px;
  color: var(--color-text-tertiary);
  margin-bottom: 4px;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-terracotta);
}

.action-buttons {
  margin: 30px 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.qr-popup-content {
  padding: 30px;
  text-align: center;
}

.qr-code-val {
  font-size: 18px;
  font-weight: bold;
  margin: 15px 0 5px;
  color: var(--color-terracotta);
}

.qr-tip {
  font-size: 12px;
  color: var(--color-text-muted);
}
</style>
