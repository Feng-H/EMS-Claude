<template>
  <div class="h5-analytics-view">
    <mobile-header title="统计分析" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <div class="content">
        <!-- 核心指标卡片 -->
        <div class="overview-grid">
          <div class="overview-item">
            <div class="val">{{ overview.equipment?.total_equipment || 0 }}</div>
            <div class="lab">设备总数</div>
          </div>
          <div class="overview-item">
            <div class="val">{{ formatDecimal(overview.mttr_mtbf?.availability) }}%</div>
            <div class="lab">可用率</div>
          </div>
          <div class="overview-item">
            <div class="val">{{ formatDecimal(overview.mttr_mtbf?.mttr) }}h</div>
            <div class="lab">MTTR</div>
          </div>
          <div class="overview-item">
            <div class="val">{{ formatDecimal(overview.tasks?.repair_completion_rate) }}%</div>
            <div class="lab">维修完成率</div>
          </div>
        </div>

        <!-- 任务趋势图表 -->
        <van-cell-group inset class="chart-group">
          <van-cell title="任务趋势 (30天)" />
          <div class="chart-wrapper">
            <div ref="trendChartRef" class="chart-container"></div>
          </div>
        </van-cell-group>

        <!-- 故障分析图表 -->
        <van-cell-group inset class="chart-group">
          <van-cell title="故障分布 (按类型)" />
          <div class="chart-wrapper">
            <div ref="failureChartRef" class="chart-container"></div>
          </div>
        </van-cell-group>

        <!-- 故障TOP设备 -->
        <van-cell-group inset title="故障 TOP5 设备">
          <van-cell
            v-for="(item, index) in topFailures.slice(0, 5)"
            :key="item.equipment_id"
            :title="`${index + 1}. ${item.equipment_name}`"
            :label="`故障: ${item.failure_count}次 | MTTR: ${formatDecimal(item.mttr)}h`"
          >
            <template #value>
              <span class="downtime">{{ formatDecimal(item.downtime_hours) }}h</span>
            </template>
          </van-cell>
        </van-cell-group>
      </div>
    </van-pull-refresh>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, shallowRef, nextTick } from 'vue'
import { showToast } from 'vant'
import echarts from '@/utils/echarts'
import type { EChartsOption } from 'echarts'
import MobileHeader from '@/components/mobile/MobileHeader.vue'
import {
  getDashboardOverview,
  getTrendData,
  getFailureAnalysis,
  getTopFailureEquipment,
  type DashboardOverview,
  type TrendData,
  type FailureAnalysis,
  type TopFailureEquipment
} from '@/api/analytics'

const refreshing = ref(false)
const overview = ref<DashboardOverview>({
  equipment: { total_equipment: 0, running_equipment: 0, stopped_equipment: 0, maintenance_equipment: 0, scrapped_equipment: 0 },
  mttr_mtbf: { mttr: 0, mtbf: 0, availability: 0 },
  tasks: { inspection_completion_rate: 0, maintenance_completion_rate: 0, repair_completion_rate: 0 },
  pending_inspections: 0,
  pending_maintenances: 0,
  pending_repairs: 0,
  low_stock_alerts: 0
})

const trendChartRef = ref<HTMLElement>()
const failureChartRef = ref<HTMLElement>()
const topFailures = ref<TopFailureEquipment[]>([])
const trendData = ref<TrendData[]>([])
const failureAnalysis = ref<FailureAnalysis[]>([])

const trendChart = shallowRef<echarts.ECharts>()
const failureChart = shallowRef<echarts.ECharts>()

const formatDecimal = (val: number | undefined) => {
  return val ? val.toFixed(2) : '0.00'
}

const loadData = async () => {
  try {
    const [ovRes, trendRes, failRes, topRes] = await Promise.all([
      getDashboardOverview(),
      getTrendData({ 
        start_date: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        end_date: new Date().toISOString().split('T')[0]
      }),
      getFailureAnalysis({ limit: 6 }),
      getTopFailureEquipment({ limit: 5 })
    ])

    overview.value = ovRes.data
    trendData.value = trendRes.data
    failureAnalysis.value = failRes.data
    topFailures.value = topRes.data

    await nextTick()
    renderCharts()
  } catch (err) {
    showToast('加载数据失败')
  } finally {
    refreshing.value = false
  }
}

const onRefresh = () => {
  loadData()
}

const renderCharts = () => {
  if (trendChartRef.value) {
    if (!trendChart.value) trendChart.value = echarts.init(trendChartRef.value)
    
    const dates = trendData.value.map(d => d.date.slice(5))
    const option: EChartsOption = {
      grid: { top: 30, bottom: 30, left: 40, right: 10 },
      legend: { bottom: 0, icon: 'circle', itemWidth: 8 },
      xAxis: { type: 'category', data: dates },
      yAxis: { type: 'value' },
      series: [
        { name: '点检', type: 'line', data: trendData.value.map(d => d.inspection_tasks), smooth: true, showSymbol: false },
        { name: '维修', type: 'line', data: trendData.value.map(d => d.repair_orders), smooth: true, showSymbol: false }
      ]
    }
    trendChart.value.setOption(option)
  }

  if (failureChartRef.value) {
    if (!failureChart.value) failureChart.value = echarts.init(failureChartRef.value)
    const option: EChartsOption = {
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: false },
        data: failureAnalysis.value.map(f => ({ name: f.equipment_type_name, value: f.failure_count }))
      }]
    }
    failureChart.value.setOption(option)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.h5-analytics-view {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding-top: 46px;
  padding-bottom: 30px;
}

.content {
  padding: var(--space-md);
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.overview-item {
  background: var(--color-bg-card);
  padding: 16px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  text-align: center;
}

.val {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-terracotta);
  margin-bottom: 4px;
}

.lab {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.chart-group {
  margin-bottom: 16px;
}

.chart-wrapper {
  padding: 10px;
  background: var(--color-bg-card);
}

.chart-container {
  height: 220px;
  width: 100%;
}

.downtime {
  color: var(--color-danger);
  font-weight: 600;
}
</style>
