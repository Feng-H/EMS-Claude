<template>
  <div class="h5-equipment-list">
    <mobile-header title="设备台账" />

    <!-- 搜索栏 -->
    <van-search
      v-model="queryParams.name"
      placeholder="搜索设备名称或编号"
      @search="onSearch"
      @clear="onReset"
    />

    <!-- 筛选 -->
    <van-dropdown-menu active-color="#c96442">
      <van-dropdown-item v-model="queryParams.type_id" :options="typeOptions" title="设备类型" @change="onSearch" />
      <van-dropdown-item v-model="queryParams.status" :options="statusOptions" title="设备状态" @change="onSearch" />
    </van-dropdown-menu>

    <!-- 列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <div class="list-container">
        <template v-if="loading && !refreshing">
          <van-skeleton v-for="i in 5" :key="i" title avatar :row="2" style="margin-bottom: 20px" />
        </template>
        <van-list
          v-else
          v-model:loading="loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
        >
          <van-empty v-if="equipmentList.length === 0" description="暂无设备数据" />
          
          <div
            v-for="item in equipmentList"
            :key="item.id"
            class="equipment-card"
            @click="navigateToDetail(item.id)"
          >
            <div class="card-header">
              <span class="equipment-name">{{ item.name }}</span>
              <van-tag :type="getStatusType(item.status)" plain>{{ getStatusText(item.status) }}</van-tag>
            </div>
            <div class="card-body">
              <div class="info-row">
                <span class="label">编号：</span>
                <span class="value">{{ item.code }}</span>
              </div>
              <div class="info-row">
                <span class="label">类型：</span>
                <span class="value">{{ item.type_name }}</span>
              </div>
              <div class="info-row">
                <span class="label">位置：</span>
                <span class="value">{{ item.factory_name }} / {{ item.workshop_name }}</span>
              </div>
            </div>
            <div class="card-footer" @click.stop>
              <van-button size="small" icon="scan" @click="showQRCode(item)">二维码</van-button>
              <van-button size="small" type="primary" plain @click="reportRepair(item)">报修</van-button>
            </div>
          </div>
        </van-list>
      </div>
    </van-pull-refresh>

    <!-- 二维码弹窗 -->
    <van-dialog v-model:show="qrDialog.show" :title="qrDialog.name" theme="round-button">
      <div class="qr-popup-content">
        <van-icon name="qr" size="180" color="#333" />
        <p class="qr-code-val">{{ qrDialog.code }}</p>
        <p class="qr-tip">现场扫描此码可快速报修或点检</p>
      </div>
    </van-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { equipmentApi, equipmentTypeApi, type Equipment } from '@/api/equipment'
import MobileHeader from '@/components/mobile/MobileHeader.vue'

const router = useRouter()

const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const equipmentList = ref<Equipment[]>([])

const queryParams = reactive({
  page: 1,
  page_size: 10,
  name: '',
  type_id: 0,
  status: ''
})

const typeOptions = ref([{ text: '全部类型', value: 0 }])
const statusOptions = [
  { text: '全部状态', value: '' },
  { text: '运行中', value: 'running' },
  { text: '已停机', value: 'stopped' },
  { text: '维修中', value: 'maintenance' },
  { text: '已报废', value: 'scrapped' }
]

const qrDialog = reactive({
  show: false,
  name: '',
  code: ''
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

const loadTypes = async () => {
  try {
    const res = await equipmentTypeApi.getTypes()
    typeOptions.value = [
      { text: '全部类型', value: 0 },
      ...res.data.map(t => ({ text: t.name, value: t.id }))
    ]
  } catch (err) {
    console.error(err)
  }
}

const loadData = async () => {
  try {
    const res = await equipmentApi.getList({
      ...queryParams,
      type_id: queryParams.type_id || undefined
    })
    
    if (refreshing.value) {
      equipmentList.value = res.data.items
    } else {
      equipmentList.value.push(...res.data.items)
    }

    if (equipmentList.value.length >= res.data.total) {
      finished.value = true
    } else {
      queryParams.page++
    }
  } catch (err) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const onRefresh = () => {
  finished.value = false
  queryParams.page = 1
  loadData()
}

const onLoad = () => {
  loadData()
}

const onSearch = () => {
  equipmentList.value = []
  finished.value = false
  queryParams.page = 1
  loadData()
}

const onReset = () => {
  queryParams.name = ''
  queryParams.type_id = 0
  queryParams.status = ''
  onSearch()
}

const navigateToDetail = (id: number) => {
  router.push(`/h5/equipment/detail/${id}`)
}

const showQRCode = (item: Equipment) => {
  qrDialog.name = item.name
  qrDialog.code = item.code
  qrDialog.show = true
}

const reportRepair = (item: Equipment) => {
  router.push(`/h5/repair/report?equipmentId=${item.id}`)
}

onMounted(() => {
  loadTypes()
})
</script>

<style scoped>
.h5-equipment-list {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 60px;
}

.list-container {
  padding: var(--space-md);
}

.equipment-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  padding: 16px;
  margin-bottom: var(--space-md);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
  transition: all 0.2s;
}

.equipment-card:active {
  background: var(--color-bg-tertiary);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.equipment-name {
  font-weight: 600;
  font-size: 16px;
  color: var(--color-text-primary);
  font-family: var(--font-serif);
}

.card-body {
  margin-bottom: 16px;
}

.info-row {
  font-size: 13px;
  margin-bottom: 6px;
  display: flex;
}

.label {
  color: var(--color-text-tertiary);
  width: 45px;
  flex-shrink: 0;
}

.value {
  color: var(--color-text-secondary);
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  border-top: 1px solid var(--color-border);
  padding-top: 12px;
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

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .equipment-card {
    background: var(--color-bg-card);
  }
}
</style>
