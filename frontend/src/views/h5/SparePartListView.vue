<template>
  <div class="h5-sparepart-list">
    <mobile-header title="备件库存" />

    <!-- 数据概览 -->
    <div class="stats-cards">
      <div class="stat-card blue">
        <div class="stat-val">{{ stats.total_parts }}</div>
        <div class="stat-label">备件总数</div>
      </div>
      <div class="stat-card red">
        <div class="stat-val">{{ stats.low_stock_count }}</div>
        <div class="stat-label">库存预警</div>
      </div>
      <div class="stat-card orange">
        <div class="stat-val">{{ stats.monthly_consumption }}</div>
        <div class="stat-label">本月消耗</div>
      </div>
    </div>

    <!-- 标签页 -->
    <van-tabs v-model:active="activeTab" class="sparepart-tabs" sticky offset-top="46">
      <van-tab title="备件列表" name="parts">
        <!-- 搜索 -->
        <van-search
          v-model="filterForm.name"
          placeholder="搜索备件名称或编码"
          @search="onSearch"
          @clear="onReset"
        />

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
              <van-empty v-if="parts.length === 0" description="暂无备件数据" />
              
              <div
                v-for="item in parts"
                :key="item.id"
                class="sparepart-card"
              >
                <div class="card-header">
                  <span class="part-name">{{ item.name }}</span>
                  <van-tag :type="item.current_stock < item.safety_stock ? 'danger' : 'success'" plain>
                    库存: {{ item.current_stock || 0 }}
                  </van-tag>
                </div>
                <div class="card-body">
                  <div class="info-row">
                    <span class="label">编码：</span>
                    <span class="value">{{ item.code }}</span>
                  </div>
                  <div class="info-row">
                    <span class="label">规格：</span>
                    <span class="value">{{ item.specification || '-' }}</span>
                  </div>
                  <div class="info-row">
                    <span class="label">工厂：</span>
                    <span class="value">{{ item.factory_name || '-' }}</span>
                  </div>
                </div>
                <div class="card-footer" v-if="canManage">
                  <van-button size="small" type="primary" plain @click="showStockDialog(item)">入库</van-button>
                </div>
              </div>
            </van-list>
          </div>
        </van-pull-refresh>
      </van-tab>

      <van-tab title="消耗记录" name="consumption">
        <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
          <div class="list-container">
            <van-list
              v-model:loading="loading"
              :finished="finished"
              finished-text="没有更多了"
              @load="onLoad"
            >
              <van-empty v-if="consumptions.length === 0" description="暂无消耗记录" />
              
              <div
                v-for="item in consumptions"
                :key="item.id"
                class="record-card"
              >
                <div class="card-header">
                  <span class="part-name">{{ item.spare_part_name }}</span>
                  <span class="qty-out">-{{ item.quantity }}</span>
                </div>
                <div class="card-meta">
                  <span>{{ item.user_name }}</span>
                  <span>{{ formatDateTime(item.created_at) }}</span>
                </div>
                <div class="card-remark" v-if="item.remark">
                  备注: {{ item.remark }}
                </div>
              </div>
            </van-list>
          </div>
        </van-pull-refresh>
      </van-tab>
    </van-tabs>

    <!-- 入库弹窗 -->
    <van-dialog v-model:show="showStockInDialog" title="备件入库" show-cancel-button @confirm="handleStockIn">
      <div class="stock-form">
        <van-cell title="备件" :value="currentPart?.name" label-class="form-label" />
        <van-field
          v-model="stockForm.factory_name"
          is-link
          readonly
          label="工厂"
          placeholder="点击选择工厂"
          @click="showFactoryPicker = true"
        />
        <van-field name="stepper" label="入库数量">
          <template #input>
            <van-stepper v-model="stockForm.quantity" min="1" />
          </template>
        </van-field>
        <van-field
          v-model="stockForm.remark"
          rows="2"
          autosize
          label="备注"
          type="textarea"
          placeholder="请输入备注"
        />
      </div>
    </van-dialog>

    <!-- 工厂选择器 -->
    <van-popup v-model:show="showFactoryPicker" position="bottom">
      <van-picker
        :columns="factoryOptions"
        @confirm="onFactoryConfirm"
        @cancel="showFactoryPicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showNotify } from 'vant'
import {
  getSpareParts,
  getConsumptions,
  getSparePartStatistics,
  stockIn,
  type SparePart,
  type SparePartStatistics
} from '@/api/sparepart'
import { equipmentApi, type Factory } from '@/api/equipment'
import MobileHeader from '@/components/mobile/MobileHeader.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const canManage = computed(() => authStore.hasRole('admin', 'engineer'))

const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const activeTab = ref('parts')

const parts = ref<SparePart[]>([])
const consumptions = ref<any[]>([])
const stats = ref<SparePartStatistics>({
  total_parts: 0,
  low_stock_count: 0,
  total_stock_value: 0,
  monthly_consumption: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 10
})

const filterForm = reactive({
  name: ''
})

// 入库相关
const showStockInDialog = ref(false)
const showFactoryPicker = ref(false)
const currentPart = ref<SparePart | null>(null)
const factories = ref<Factory[]>([])
const stockForm = reactive({
  spare_part_id: 0,
  factory_id: 0,
  factory_name: '',
  quantity: 1,
  remark: ''
})

const factoryOptions = computed(() => 
  factories.value.map(f => ({ text: f.name, value: f.id }))
)

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}-${date.getDate()} ${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
}

const loadStats = async () => {
  try {
    const res = await getSparePartStatistics()
    stats.value = res.data
  } catch (err) {}
}

const loadParts = async () => {
  try {
    const res = await getSpareParts({
      name: filterForm.name || undefined,
      page: pagination.page,
      page_size: pagination.pageSize
    })
    
    if (refreshing.value) {
      parts.value = res.data.items
    } else {
      parts.value.push(...res.data.items)
    }

    if (parts.value.length >= res.data.total) {
      finished.value = true
    } else {
      pagination.page++
    }
  } catch (err) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const loadConsumptionsList = async () => {
  try {
    const res = await getConsumptions({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    
    if (refreshing.value) {
      consumptions.value = res.data.items
    } else {
      consumptions.value.push(...res.data.items)
    }

    if (consumptions.value.length >= res.data.total) {
      finished.value = true
    } else {
      pagination.page++
    }
  } catch (err) {
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const loadFactories = async () => {
  try {
    const res = await equipmentApi.getFactories()
    factories.value = res.data
  } catch (err) {}
}

const onLoad = () => {
  if (activeTab.value === 'parts') {
    loadParts()
  } else {
    loadConsumptionsList()
  }
}

const onRefresh = () => {
  finished.value = false
  pagination.page = 1
  onLoad()
}

const onSearch = () => {
  parts.value = []
  onRefresh()
}

const onReset = () => {
  filterForm.name = ''
  onSearch()
}

const showStockDialog = (item: SparePart) => {
  currentPart.value = item
  stockForm.spare_part_id = item.id
  stockForm.factory_id = 0
  stockForm.factory_name = ''
  stockForm.quantity = 1
  stockForm.remark = ''
  showStockInDialog.value = true
}

const onFactoryConfirm = ({ selectedOptions }: any) => {
  stockForm.factory_id = selectedOptions[0].value
  stockForm.factory_name = selectedOptions[0].text
  showFactoryPicker.value = false
}

const handleStockIn = async () => {
  if (stockForm.factory_id === 0) {
    showNotify({ type: 'warning', message: '请选择工厂' })
    return
  }
  try {
    await stockIn({
      spare_part_id: stockForm.spare_part_id,
      factory_id: stockForm.factory_id,
      quantity: stockForm.quantity,
      remark: stockForm.remark
    })
    showToast('入库成功')
    onRefresh()
    loadStats()
  } catch (err: any) {
    showNotify({ type: 'danger', message: err.response?.data?.error || '操作失败' })
  }
}

watch(activeTab, () => {
  onRefresh()
})

onMounted(() => {
  loadStats()
  loadFactories()
})
</script>

<style scoped>
.h5-sparepart-list {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 60px;
}

.stats-cards {
  display: flex;
  gap: 10px;
  padding: 16px;
}

.stat-card {
  flex: 1;
  padding: 12px 8px;
  border-radius: var(--radius-high);
  text-align: center;
  color: #fff;
  box-shadow: var(--shadow-sm);
}

.stat-card.blue { background: #409eff; }
.stat-card.red { background: #f56c6c; }
.stat-card.orange { background: #e6a23c; }

.stat-val {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 11px;
  opacity: 0.9;
}

.list-container {
  padding: 12px;
}

.sparepart-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--color-border);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.part-name {
  font-weight: 600;
  font-size: 15px;
  color: var(--color-text-primary);
}

.info-row {
  font-size: 13px;
  margin-bottom: 6px;
  display: flex;
}

.label {
  color: var(--color-text-tertiary);
  width: 45px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
}

.record-card {
  background: var(--color-bg-card);
  border-radius: var(--radius-very);
  padding: 12px 16px;
  margin-bottom: 10px;
  border: 1px solid var(--color-border);
}

.qty-out {
  color: #f56c6c;
  font-weight: 600;
  font-size: 16px;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--color-text-tertiary);
  margin-top: 8px;
}

.card-remark {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-top: 6px;
  background: var(--color-bg-tertiary);
  padding: 4px 8px;
  border-radius: 4px;
}

.stock-form {
  padding: 20px 0;
}

@media (prefers-color-scheme: dark) {
  .sparepart-card, .record-card {
    background: var(--color-bg-card);
  }
}
</style>
