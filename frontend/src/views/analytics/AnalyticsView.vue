<template>
  <div class="analytics-view">
    <div class="header">
      <div class="header-left">
        <h2>管理驾驶舱</h2>
        <span class="header-subtitle">设备运行、维保效率及资产健康度实时分析</span>
      </div>
      <div class="header-actions">
        <el-select v-model="filterFactoryId" placeholder="全工厂" clearable @change="handleFactoryChange" style="width: 200px">
          <el-option v-for="f in factories" :key="f.id" :label="f.name" :value="f.id" />
        </el-select>
        <el-button @click="refreshAll" :icon="Refresh">刷新</el-button>
      </div>
    </div>

    <!-- Overview Cards -->
    <el-row :gutter="16" class="overview-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card equipment">
          <div class="stat-icon"><el-icon><Box /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ overview.equipment?.total_equipment || 0 }}</div>
            <div class="stat-label">设备总数</div>
            <div class="stat-detail">
              运行: {{ overview.equipment?.running_equipment || 0 }} |
              停机: {{ overview.equipment?.stopped_equipment || 0 }}
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card mttr">
          <div class="stat-icon"><el-icon><Timer /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ formatDecimal(overview.mttr_mtbf?.mttr) }}h</div>
            <div class="stat-label">MTTR (平均修复时间)</div>
            <div class="stat-detail">MTBF: {{ formatDecimal(overview.mttr_mtbf?.mtbf) }}h</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card availability">
          <div class="stat-icon"><el-icon><CircleCheck /></el-icon></div>
          <div class="stat-content">
            <div class="stat-value">{{ formatDecimal(overview.mttr_mtbf?.availability) }}%</div>
            <div class="stat-label">设备可用率</div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
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
      <el-col :xs="24" :sm="8">
        <el-card class="pending-card warning">
          <div class="pending-item">
            <span class="pending-label">待执行点检</span>
            <span class="pending-value">{{ overview.pending_inspections || 0 }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
        <el-card class="pending-card info">
          <div class="pending-item">
            <span class="pending-label">待执行保养</span>
            <span class="pending-value">{{ overview.pending_maintenances || 0 }}</span>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="8">
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
      <el-col :xs="24" :lg="16">
        <el-card>
          <template #header>
            <span>任务趋势 (近30天)</span>
          </template>
          <div ref="trendChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card>
          <template #header>
            <span>故障分析 (按设备类型)</span>
          </template>
          <div ref="failureChartRef" class="chart-container" v-loading="chartLoading"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Rankings -->
    <el-row :gutter="16" class="ranking-row">
      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>MTBF 排名 (高 -> 低)</span>
              <el-tooltip content="平均无故障时间，反映设备可靠性"><el-icon><QuestionFilled /></el-icon></el-tooltip>
            </div>
          </template>
          <el-table :data="mtbfRanking" size="small" stripe @row-click="handleRankingClick" class="clickable-table">
            <el-table-column type="index" label="#" width="40" />
            <el-table-column prop="equipment_name" label="设备" show-overflow-tooltip />
            <el-table-column label="MTBF(h)" width="90" align="right">
              <template #default="{ row }">{{ formatDecimal(row.value) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>停机时间排名 (高 -> 低)</span>
              <el-tooltip content="累计故障停机时长，反映停机损失"><el-icon><QuestionFilled /></el-icon></el-tooltip>
            </div>
          </template>
          <el-table :data="downtimeRanking" size="small" stripe @row-click="handleRankingClick" class="clickable-table">
            <el-table-column type="index" label="#" width="40" />
            <el-table-column prop="equipment_name" label="设备" show-overflow-tooltip />
            <el-table-column label="停机(h)" width="90" align="right">
              <template #default="{ row }">{{ formatDecimal(row.value) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>可用性评分排名</span>
              <el-tooltip content="综合可用性得分，100为最优"><el-icon><QuestionFilled /></el-icon></el-tooltip>
            </div>
          </template>
          <el-table :data="performanceRanking" size="small" stripe @row-click="handleRankingClick" class="clickable-table">
            <el-table-column type="index" label="#" width="40" />
            <el-table-column prop="equipment_name" label="设备" show-overflow-tooltip />
            <el-table-column label="评分" width="80" align="right">
              <template #default="{ row }">
                <el-text :type="row.value > 90 ? 'success' : row.value > 70 ? 'warning' : 'danger'">
                  {{ formatDecimal(row.value) }}
                </el-text>
              </template>
            </el-table-column>
          </el-table>
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
          <el-table :data="topFailures" stripe v-loading="loading" @row-click="handleRankingClick" class="clickable-table">
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
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Box, Timer, CircleCheck, DataAnalysis, Refresh, QuestionFilled } from '@element-plus/icons-vue'
import echarts from '@/utils/echarts'
import type { EChartsOption } from 'echarts'
import {
  getDashboardOverview,
  getTrendData,
  getFailureAnalysis,
  getTopFailureEquipment,
  getMTBFRanking,
  getDowntimeRanking,
  getPerformanceRanking,
  type DashboardOverview,
  type TrendData,
  type FailureAnalysis,
  type TopFailureEquipment,
  type EquipmentRanking
} from '@/api/analytics'
import { equipmentApi, type Factory } from '@/api/equipment'

const loading = ref(false)
const chartLoading = ref(false)
const filterFactoryId = ref<number>()
const factories = ref<Factory[]>([])
const router = useRouter()

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
const mtbfRanking = ref<EquipmentRanking[]>([])
const downtimeRanking = ref<EquipmentRanking[]>([])
const performanceRanking = ref<EquipmentRanking[]>([])

const trendData = ref<TrendData[]>([])
const failureAnalysis = ref<FailureAnalysis[]>([])

const trendChart = shallowRef<echarts.ECharts>()
const failureChart = shallowRef<echarts.ECharts>()

const formatDecimal = (val: number | undefined) => {
  return val ? val.toFixed(2) : '0.00'
}

const handleRankingClick = (row: any) => {
  const id = row.equipment_id || row.id
  if (id) router.push(`/equipment/detail/${id}`)
}

const loadFactories = async () => {
  try {
    const res = await equipmentApi.getFactories()
    factories.value = res.data
  } catch (err) {
    console.error('Failed to load factories')
  }
}

const handleFactoryChange = () => {
  refreshAll()
}

const refreshAll = () => {
  loadOverview()
  loadTrendData()
  loadFailureAnalysis()
  loadTopFailures()
  loadRankings()
}

const loadOverview = async () => {
  loading.value = true
  try {
    const res = await getDashboardOverview({ factory_id: filterFactoryId.value })
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

    const res = await getTrendData({ 
      start_date: startDate, 
      end_date: endDate,
      factory_id: filterFactoryId.value
    })
    trendData.value = res.data
    renderTrendChart()
  } catch (err) {
    ElMessage.error('加载趋势数据失败')
  } finally {
    chartLoading.value = false
  }
}

const loadRankings = async () => {
  try {
    const fId = filterFactoryId.value
    const [mtbf, down, perf] = await Promise.all([
      getMTBFRanking({ limit: 5, factory_id: fId }),
      getDowntimeRanking({ limit: 5, factory_id: fId }),
      getPerformanceRanking({ limit: 5, factory_id: fId })
    ])
    mtbfRanking.value = mtbf.data
    downtimeRanking.value = down.data
    performanceRanking.value = perf.data
  } catch (err) {
    console.error('Failed to load rankings')
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
  loadFactories()
  refreshAll()
})
</script>

<style scoped>
.analytics-view {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
}

.header-subtitle {
  font-size: 14px;
  color: #8c8c8c;
  margin-top: 4px;
  display: block;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.overview-row,
.pending-row,
.charts-row,
.ranking-row,
.failure-row {
  margin-bottom: 20px;
}

.stat-card {
  position: relative;
  overflow: hidden;
  border-radius: 8px;
  border: none;
  box-shadow: 0 4px 12px rgba(0,0,0,0.05);
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

.stat-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
}

.stat-icon {
  font-size: 40px;
  opacity: 0.8;
  color: #f0f2f5;
  background: #fafafa;
  padding: 12px;
  border-radius: 50%;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1f1f1f;
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: #8c8c8c;
  margin: 4px 0;
}

.stat-detail {
  font-size: 12px;
  color: #bfbfbf;
}

.pending-card {
  border-radius: 8px;
}

.pending-card.warning { background-color: #fffbe6; border: 1px solid #ffe58f; }
.pending-card.info { background-color: #e6f7ff; border: 1px solid #91d5ff; }
.pending-card.danger { background-color: #fff1f0; border: 1px solid #ffa39e; }

.pending-card.warning .pending-value { color: #faad14; font-size: 24px; font-weight: bold; }
.pending-card.info .pending-value { color: #1890ff; font-size: 24px; font-weight: bold; }
.pending-card.danger .pending-value { color: #f5222d; font-size: 24px; font-weight: bold; }

.pending-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pending-label {
  font-size: 14px;
  font-weight: 500;
  color: #595959;
}

.chart-container {
  height: 320px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header span {
  font-weight: 600;
  font-size: 15px;
}

.clickable-table :deep(.el-table__row) {
  cursor: pointer;
}
</style>
