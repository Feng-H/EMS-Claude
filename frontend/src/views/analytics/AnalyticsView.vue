<template>
  <div class="analytics-view">
    <div class="header">
      <h2>统计分析</h2>
    </div>

    <!-- Overview Cards -->
    <el-row :gutter="16" class="overview-row">
      <el-col :span="6">
        <el-card class="stat-card equipment">
          <div class="stat-icon"><el-icon><Box /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ overview.equipment?.total_equipment || 0 }}</div>
            <div class="stat-label">设备总数</div>
            <div class="stat-detail">
              运行: {{ overview.equipment?.running_equipment || 0 }} |
              停机: {{ overview.equipment?.stopped_equipment || 0 }} |
              维修: {{ overview.equipment?.maintenance_equipment || 0 }}
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card mttr">
          <div class="stat-icon"><el-icon><Timer /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ formatDecimal(overview.mttr_mtbf?.mttr) }}h</div>
            <div class="stat-label">MTTR (平均修复时间)</div>
            <div class="stat-detail">MTBF: {{ formatDecimal(overview.mttr_mtbf?.mtbf) }}h</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card availability">
          <div class="stat-icon"><el-icon><CircleCheck /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ formatDecimal(overview.mttr_mtbf?.availability) }}%</div>
            <div class="stat-label">设备可用率</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card completion">
          <div class="stat-icon"><el-icon><DataAnalysis /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ formatDecimal(overview.tasks?.repair_completion_rate) }}%</div>
            <div class="stat-label">维修完成率</div>
            <div class="stat-detail">
              点检: {{ formatDecimal(overview.tasks?.inspection_completion_rate) }}% |
              保养: {{ formatDecimal(overview.tasks?.maintenance_completion_rate) }}%
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Pending Tasks -->
    <el-row :gutter="16" class="pending-row">
      <el-col :span="8">
        <el-card class="pending-card warning">
          <div class="pending-item">
            <span class="pending-label">待执行点检</span>
            <span class="pending-value">{{ overview.pending_inspections || 0 }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="pending-card info">
          <div class="pending-item">
            <span class="pending-label">待执行保养</span>
            <span class="pending-value">{{ overview.pending_maintenances || 0 }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="pending-card danger">
          <div class="pending-item">
            <span class="pending-label">待处理维修</span>
            <span class="pending-value">{{ overview.pending_repairs || 0 }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Charts -->
    <el-row :gutter="16" class="charts-row">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>任务趋势 (近30天)</span>
          </template>
          <div ref="trendChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>故障分析 (按设备类型)</span>
          </template>
          <div ref="failureChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Top Failure Equipment -->
    <el-row :gutter="16" class="failure-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>故障TOP10设备</span>
          </template>
          <el-table :data="topFailures" stripe v-loading="loading">
            <el-table-column type="index" label="排名" width="60" align="center" />
            <el-table-column prop="equipment_code" label="设备编号" width="120" />
            <el-table-column prop="equipment_name" label="设备名称" min-width="180" />
            <el-table-column prop="failure_count" label="故障次数" width="100" align="center" />
            <el-table-column prop="downtime_hours" label="停机时长(小时)" width="130" align="right">
              <template #default="{ row }">{{ formatDecimal(row.downtime_hours) }}</template>
            </el-table-column>
            <el-table-column prop="mttr" label="MTTR(小时)" width="110" align="right">
              <template #default="{ row }">{{ formatDecimal(row.mttr) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, shallowRef } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
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

const loading = ref(false)
const chartLoading = ref(false)
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

const loadOverview = async () => {
  loading.value = true
  try {
    const res = await getDashboardOverview()
    overview.value = res.data
  } catch (err) {
    ElMessage.error('加载概览数据失败')
  } finally {
    loading.value = false
  }
}

const loadTrendData = async () => {
  chartLoading.value = true
  try {
    const endDate = new Date().toISOString().split('T')[0]
    const startDate = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]

    const res = await getTrendData({ start_date: startDate, end_date: endDate })
    trendData.value = res.data
    renderTrendChart()
  } catch (err) {
    ElMessage.error('加载趋势数据失败')
  } finally {
    chartLoading.value = false
  }
}

const loadFailureAnalysis = async () => {
  chartLoading.value = true
  try {
    const res = await getFailureAnalysis({ limit: 6 })
    failureAnalysis.value = res.data
    renderFailureChart()
  } catch (err) {
    ElMessage.error('加载故障分析失败')
  } finally {
    chartLoading.value = false
  }
}

const loadTopFailures = async () => {
  try {
    const res = await getTopFailureEquipment({ limit: 10 })
    topFailures.value = res.data
  } catch (err) {
    ElMessage.error('加载TOP设备失败')
  }
}

const renderTrendChart = () => {
  if (!trendChartRef.value) return

  if (!trendChart.value) {
    trendChart.value = echarts.init(trendChartRef.value)
  }

  const dates = trendData.value.map(d => d.date.slice(5)) // MM-DD
  const inspectionData = trendData.value.map(d => d.inspection_tasks)
  const maintenanceData = trendData.value.map(d => d.maintenance_tasks)
  const repairData = trendData.value.map(d => d.repair_orders)

  const option: EChartsOption = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['点检', '保养', '维修']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '点检',
        type: 'line',
        data: inspectionData,
        smooth: true,
        itemStyle: { color: '#409eff' }
      },
      {
        name: '保养',
        type: 'line',
        data: maintenanceData,
        smooth: true,
        itemStyle: { color: '#67c23a' }
      },
      {
        name: '维修',
        type: 'line',
        data: repairData,
        smooth: true,
        itemStyle: { color: '#f56c6c' }
      }
    ]
  }

  trendChart.value.setOption(option)
}

const renderFailureChart = () => {
  if (!failureChartRef.value) return

  if (!failureChart.value) {
    failureChart.value = echarts.init(failureChartRef.value)
  }

  const types = failureAnalysis.value.map(f => f.equipment_type_name)
  const counts = failureAnalysis.value.map(f => f.failure_count)

  const option: EChartsOption = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}次'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: failureAnalysis.value.map((f, i) => ({
          name: f.equipment_type_name,
          value: f.failure_count,
          itemStyle: {
            color: ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'][i % 9]
          }
        })),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        label: {
          formatter: '{b}: {c}'
        }
      }
    ]
  }

  failureChart.value.setOption(option)
}

onMounted(() => {
  loadOverview()
  loadTrendData()
  loadFailureAnalysis()
  loadTopFailures()
})
</script>

<style scoped>
.analytics-view {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.overview-row,
.pending-row,
.charts-row,
.failure-row {
  margin-bottom: 20px;
}

.stat-card {
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 4px;
  height: 100%;
}

.stat-card.equipment::before { background: #409eff; }
.stat-card.mttr::before { background: #e6a23c; }
.stat-card.availability::before { background: #67c23a; }
.stat-card.completion::before { background: #909399; }

.stat-card .el-card__body {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  font-size: 48px;
  opacity: 0.2;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin: 4px 0;
}

.stat-detail {
  font-size: 12px;
  color: #606266;
}

.pending-card.warning .pending-value { color: #e6a23c; font-size: 24px; font-weight: bold; }
.pending-card.info .pending-value { color: #409eff; font-size: 24px; font-weight: bold; }
.pending-card.danger .pending-value { color: #f56c6c; font-size: 24px; font-weight: bold; }

.pending-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pending-label {
  font-size: 14px;
  color: #606266;
}

.chart-container {
  height: 300px;
}
</style>
