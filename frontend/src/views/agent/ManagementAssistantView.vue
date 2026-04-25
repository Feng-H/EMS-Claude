<template>
  <div class="management-assistant">
    <div class="assistant-container">
      <!-- 左侧：控制面板 -->
      <div class="control-panel">
        <el-card shadow="never" class="panel-card">
          <template #header>
            <div class="panel-header">
              <el-icon><MagicStick /></el-icon>
              <span>智能管理配置</span>
            </div>
          </template>

          <el-tabs v-model="activeScenario" class="scenario-tabs">
            <el-tab-pane label="保养优化" name="maintenance">
              <el-form :model="maintenanceForm" label-position="top">
                <el-form-item label="设备类型">
                  <el-select v-model="maintenanceForm.equipment_type_id" placeholder="选择类型" style="width: 100%">
                    <el-option v-for="t in equipmentTypes" :key="t.id" :label="t.name" :value="t.id" />
                  </el-select>
                </el-form-item>
                <el-form-item label="分析范围">
                  <el-date-picker
                    v-model="maintenanceForm.dateRange"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始"
                    end-placeholder="结束"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                </el-form-item>
                <el-button type="primary" class="action-btn" @click="handleRunMaintenance" :loading="loading">
                  生成建议
                </el-button>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="维修审计" name="repair">
              <el-form :model="repairForm" label-position="top">
                <el-form-item label="所属工厂">
                  <el-select v-model="repairForm.factory_id" placeholder="所有工厂" clearable style="width: 100%">
                    <el-option v-for="f in factories" :key="f.id" :label="f.name" :value="f.id" />
                  </el-select>
                </el-form-item>
                <el-form-item label="审计时间段">
                  <el-date-picker
                    v-model="repairForm.dateRange"
                    type="daterange"
                    range-separator="至"
                    start-placeholder="开始"
                    end-placeholder="结束"
                    value-format="YYYY-MM-DD"
                    style="width: 100%"
                  />
                </el-form-item>
                <el-button type="warning" class="action-btn" @click="handleRunRepairAudit" :loading="loading">
                  运行审计
                </el-button>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>

        <!-- 历史会话 -->
        <el-card shadow="never" class="history-card" v-if="sessions.length > 0">
          <template #header>
            <div class="panel-header">
              <el-icon><History /></el-icon>
              <span>最近分析</span>
            </div>
          </template>
          <div v-for="s in sessions" :key="s.id" class="history-item" @click="loadSession(s.id)">
            <div class="history-title">{{ s.scenario === 'maintenance_recommendation' ? '保养建议' : '维修审计' }}</div>
            <div class="history-time">{{ formatDate(s.created_at) }}</div>
          </div>
        </el-card>
      </div>

      <!-- 右侧：结果展示 -->
      <div class="result-panel" v-loading="loading">
        <div v-if="!result && !loading" class="empty-state">
          <el-empty description="在左侧配置参数并点击运行，开启 AI 辅助管理">
            <template #image>
              <el-icon :size="80" color="#e4e7ed"><Opportunity /></el-icon>
            </template>
          </el-empty>
        </div>

        <div v-else-if="result" class="result-content animate-fade-in">
          <!-- 核心结论区 -->
          <div class="conclusion-section">
            <div class="section-badge" :class="result.risk_level">
              {{ result.risk_level === 'high' ? '高风险' : result.risk_level === 'medium' ? '中等风险' : '正常' }}
            </div>
            <h2 class="conclusion-summary">{{ result.summary }}</h2>
            <div class="trace-info">Trace ID: {{ result.trace_id }}</div>
          </div>

          <!-- 具体建议/异常项 -->
          <div class="items-section" v-if="hasItems">
            <h3 class="section-title">详细清单</h3>
            <div v-if="result.scenario === 'maintenance_recommendation'" class="items-list">
              <div v-for="(item, idx) in result.data.recommendations" :key="idx" class="item-card recommend">
                <div class="item-header">
                  <span class="item-title">{{ item.title }}</span>
                  <el-tag size="small">{{ item.type }}</el-tag>
                </div>
                <div class="item-desc">{{ item.description }}</div>
                <div class="item-meta">
                  <strong>理由：</strong>{{ item.reason }}
                </div>
              </div>
            </div>
            <div v-else-if="result.scenario === 'repair_audit'" class="items-list">
              <div v-for="(item, idx) in result.data.anomalies" :key="idx" class="item-card audit">
                <div class="item-header">
                  <span class="item-title">{{ item.title }}</span>
                  <el-tag :type="item.severity === 'high' ? 'danger' : 'warning'" size="small">{{ item.severity }}</el-tag>
                </div>
                <div class="item-desc">{{ item.description }}</div>
                <div class="item-meta">
                  <strong>建议行动：</strong>{{ item.suggested_action }}
                </div>
              </div>
            </div>
          </div>

          <!-- 证据支撑区 -->
          <div class="evidence-section" v-if="result.data.evidence && result.data.evidence.length > 0">
            <h3 class="section-title">参考证据</h3>
            <div class="evidence-grid">
              <div v-for="(ev, idx) in result.data.evidence" :key="idx" class="evidence-pill">
                <el-icon><Document /></el-icon>
                <span class="ev-title" :title="ev.excerpt">{{ ev.title }}</span>
                <span class="ev-score">{{ (ev.score * 100).toFixed(0) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { MagicStick, History, Opportunity, Document } from '@element-plus/icons-vue'
import { equipmentApi, type EquipmentType, type Factory } from '@/api/equipment'
import { agentApi, type AgentResponse } from '@/api/agent'
import { ElMessage } from 'element-plus'

const activeScenario = ref('maintenance')
const loading = ref(false)
const result = ref<AgentResponse<any> | null>(null)
const sessions = ref<any[]>([])

const equipmentTypes = ref<EquipmentType[]>([])
const factories = ref<Factory[]>([])

const maintenanceForm = ref({
  equipment_type_id: null as number | null,
  dateRange: [] as string[]
})

const repairForm = ref({
  factory_id: null as number | null,
  dateRange: [] as string[]
})

const hasItems = computed(() => {
  if (!result.value) return false
  if (result.value.scenario === 'maintenance_recommendation') {
    return result.value.data.recommendations?.length > 0
  }
  if (result.value.scenario === 'repair_audit') {
    return result.value.data.anomalies?.length > 0
  }
  return false
})

async function handleRunMaintenance() {
  if (!maintenanceForm.value.equipment_type_id) {
    ElMessage.warning('请选择设备类型')
    return
  }
  
  loading.value = true
  try {
    const res = await agentApi.recommendMaintenance({
      equipment_type_id: maintenanceForm.value.equipment_type_id!,
      time_range: {
        start_date: maintenanceForm.value.dateRange[0] || '2026-01-01',
        end_date: maintenanceForm.value.dateRange[1] || '2026-04-25'
      }
    })
    result.value = res.data
    await loadSessions()
  } catch (error: any) {
    ElMessage.error('分析失败: ' + (error.response?.data?.error?.message || error.message))
  } finally {
    loading.value = false
  }
}

async function handleRunRepairAudit() {
  loading.value = true
  try {
    const res = await agentApi.auditRepair({
      factory_id: repairForm.value.factory_id || undefined,
      equipment_type_id: 1, // Default for now
      time_range: {
        start_date: repairForm.value.dateRange[0] || '2026-01-01',
        end_date: repairForm.value.dateRange[1] || '2026-04-25'
      }
    })
    result.value = res.data
    await loadSessions()
  } catch (error: any) {
    ElMessage.error('分析失败: ' + (error.response?.data?.error?.message || error.message))
  } finally {
    loading.value = false
  }
}

async function loadSessions() {
  try {
    const res = await agentApi.listSessions()
    sessions.value = res.data
  } catch (error) {
    console.error('Failed to load sessions:', error)
  }
}

async function loadSession(id: number) {
  loading.value = true
  try {
    const sessionRes = await agentApi.getSession(id)
    if (sessionRes.data.artifacts?.length > 0) {
      const artifactId = sessionRes.data.artifacts[0]
      const artifactRes = await agentApi.getArtifact(artifactId)
      
      // Map artifact back to AgentResponse shape for rendering
      const art = artifactRes.data
      result.value = {
        success: true,
        trace_id: sessionRes.data.trace_id,
        language: sessionRes.data.language,
        scenario: sessionRes.data.scenario,
        scope_summary: {},
        summary: art.summary,
        risk_level: art.risk_level,
        artifact_id: art.id,
        evidence_count: art.evidence?.length || 0,
        data: art.result_json
      }
    }
  } catch (error) {
    ElMessage.error('加载历史记录失败')
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleString()
}

onMounted(async () => {
  try {
    const [typesRes, factoriesRes] = await Promise.all([
      equipmentApi.getTypes(),
      equipmentApi.getFactories(),
      loadSessions()
    ])
    equipmentTypes.value = typesRes.data
    factories.value = factoriesRes.data
  } catch (error) {
    console.error('Failed to load metadata:', error)
  }
})
</script>

<style scoped>
.management-assistant {
  height: calc(100vh - 120px);
  overflow: hidden;
}

.assistant-container {
  display: flex;
  gap: 20px;
  height: 100%;
}

.control-panel {
  width: 320px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex-shrink: 0;
}

.panel-card {
  border-radius: 12px;
}

.panel-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.action-btn {
  width: 100%;
  margin-top: 8px;
}

.history-card {
  flex: 1;
  border-radius: 12px;
  overflow-y: auto;
}

.history-item {
  padding: 12px;
  border-bottom: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 0.2s;
}

.history-item:hover {
  background: var(--color-bg-tertiary);
}

.history-title {
  font-size: 14px;
  margin-bottom: 4px;
}

.history-time {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.result-panel {
  flex: 1;
  background: var(--color-bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--color-border);
  overflow-y: auto;
  padding: 24px;
}

.empty-state {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.conclusion-section {
  margin-bottom: 32px;
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 24px;
}

.section-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 12px;
}

.section-badge.high { background: #fee2e2; color: #dc2626; }
.section-badge.medium { background: #fef3c7; color: #d97706; }
.section-badge.normal { background: #dcfce7; color: #16a34a; }

.conclusion-summary {
  font-size: 24px;
  line-height: 1.4;
  color: var(--color-text-primary);
  margin: 0 0 12px 0;
}

.trace-info {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 16px;
  color: var(--color-text-secondary);
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.item-card {
  background: var(--color-bg-primary);
  padding: 16px;
  border-radius: 8px;
  border-left: 4px solid #ccc;
}

.item-card.recommend { border-left-color: var(--color-primary); }
.item-card.audit { border-left-color: #f59e0b; }

.item-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.item-title {
  font-weight: 600;
}

.item-desc {
  font-size: 14px;
  margin-bottom: 8px;
  color: var(--color-text-secondary);
}

.item-meta {
  font-size: 13px;
  color: var(--color-text-tertiary);
}

.evidence-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.evidence-pill {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--color-bg-primary);
  border: 1px solid var(--color-border);
  border-radius: 6px;
  font-size: 13px;
}

.ev-title {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ev-score {
  color: var(--color-primary);
  font-weight: 600;
  font-family: monospace;
}

.animate-fade-in {
  animation: fadeIn 0.5s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
