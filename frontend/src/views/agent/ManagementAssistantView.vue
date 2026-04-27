<template>
  <div class="management-assistant">
    <div class="cockpit-layout">
      <!-- 左侧：导航与历史 -->
      <aside class="sidebar-nav">
        <div class="nav-group">
          <div class="group-title">AI 工作台</div>
          <div 
            class="nav-item" 
            :class="{ active: activeMode === 'chat' }"
            @click="activeMode = 'chat'"
          >
            <el-icon><ChatDotRound /></el-icon>
            <span>专家对话</span>
          </div>
          <div 
            class="nav-item" 
            :class="{ active: activeMode === 'audit' }"
            @click="activeMode = 'audit'"
          >
            <el-icon><CircleCheck /></el-icon>
            <span>专项审计</span>
          </div>
        </div>

        <div class="nav-group">
          <div class="group-title">自主学习中心</div>
          <div 
            class="nav-item" 
            :class="{ active: activeMode === 'knowledge' }"
            @click="activeMode = 'knowledge'"
          >
            <el-icon><Reading /></el-icon>
            <span>知识审核</span>
            <el-badge v-if="draftCount > 0" :value="draftCount" class="badge" />
          </div>
        </div>

        <div class="nav-group history-group">
          <div class="group-title">历史会话</div>
          <div v-for="c in conversations" :key="c.id" 
               class="history-item" 
               :class="{ active: currentConvId === c.id }"
               @click="loadConversation(c.id)">
            <div class="history-title">{{ c.title }}</div>
            <div class="history-meta">{{ formatDate(c.created_at) }}</div>
          </div>
        </div>
      </aside>

      <!-- 中间：核心交互区 -->
      <main class="main-content">
        <!-- 模式1：专家对话 -->
        <div v-if="activeMode === 'chat'" class="chat-container">
          <div class="chat-messages" ref="messageBox">
            <div v-if="messages.length === 0" class="welcome-guide">
              <div class="guide-icon">🤖</div>
              <h2>我是您的 EMS 智能资产专家</h2>
              <p>我已经整合了设备 180 天的运行数据、维修记录与 TCO 财务分析。</p>
              <div class="guide-chips">
                <el-tag @click="userInput = 'CNC-001 最近 30 天维修费分析，建议怎么优化？'">CNC-001 深度诊断</el-tag>
                <el-tag @click="userInput = 'PRESS-05 已经运行 12 年了，从财务角度建议退役吗？'">资产退役评估</el-tag>
                <el-tag @click="userInput = '帮我对比李四和张三负责区域的保养效果'">人效价值对标</el-tag>
              </div>
            </div>
            
            <div v-for="(msg, idx) in messages" :key="idx" :class="['message', msg.role]">
              <div class="avatar">{{ msg.role === 'user' ? 'U' : 'AI' }}</div>
              <div class="content">
                <div class="text" v-html="formatMessage(msg.content)"></div>
                <div v-if="msg.trace_id" class="msg-footer">Trace: {{ msg.trace_id }}</div>
              </div>
            </div>
            <div v-if="chatLoading" class="message assistant loading">
              <div class="avatar">AI</div>
              <div class="content">
                <div class="loading-tip">AI 专家正在深度分析数据并生成决策建议...</div>
                <el-skeleton :rows="2" animated />
              </div>
            </div>
          </div>
          
          <div class="chat-input-area">
            <div class="input-wrapper">
              <el-input
                v-model="userInput"
                type="textarea"
                :rows="3"
                placeholder="在此输入您的分析指令... (Shift + Enter 换行)"
                @keyup.enter.exact.prevent="handleSendChat"
              />
              <div class="input-actions">
                <el-button type="primary" @click="handleSendChat" :loading="chatLoading">发送指令</el-button>
              </div>
            </div>
          </div>
        </div>

        <!-- 模式2：专项审计 -->
        <div v-if="activeMode === 'audit'" class="audit-mode">
          <div class="p-20">
            <h3>自动化管理审计</h3>
            <p class="text-tertiary">基于预定义规则的深度合规性核查</p>
            
            <el-tabs v-model="activeAuditTab" class="mt-20">
              <el-tab-pane label="维修合理性审计" name="repair">
                 <el-form label-position="top" class="max-w-400">
                    <el-form-item label="目标设备类型">
                      <el-select v-model="auditForm.equipment_type_id" placeholder="请选择">
                        <el-option v-for="t in equipmentTypes" :key="t.id" :label="t.name" :value="t.id" />
                      </el-select>
                    </el-form-item>
                    <el-button type="warning" @click="handleRunRepairAudit">启动 AI 级联失效核查</el-button>
                 </el-form>
              </el-tab-pane>
            </el-tabs>

            <div v-if="auditResult" class="mt-24">
               <el-card shadow="never" class="result-card">
                  <pre class="summary-text">{{ auditResult.summary }}</pre>
               </el-card>
            </div>
          </div>
        </div>

        <!-- 模式3：知识审核 -->
        <div v-if="activeMode === 'knowledge'" class="knowledge-audit p-20">
          <div class="section-header">
             <h3>知识库待审核 (AI 自主提炼)</h3>
             <p class="text-tertiary">Agent 从历史对话中自动识别的高价值工业经验。</p>
          </div>
          <el-table :data="knowledgeDrafts" stripe class="mt-20">
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="type" label="类型" width="120" />
            <el-table-column prop="confidence" label="置信度" width="100">
              <template #default="{row}">
                <el-progress :percentage="Math.round(row.confidence * 100)" :format="() => ''" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="{row}">
                <el-button type="success" link @click="confirmKnowledge(row.id)">确认入库</el-button>
                <el-button type="danger" link @click="rejectKnowledge(row.id)">拒绝</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </main>

      <!-- 右侧：预测性洞察面板 -->
      <aside class="context-panel">
        <el-card shadow="never" class="prediction-card">
          <template #header><div class="card-header">设备实时工况 (CNC-001)</div></template>
          <div v-if="prediction" class="prediction-content">
            <div class="stat-main">
              <el-progress 
                type="dashboard" 
                :percentage="Math.round(prediction.rul?.health_score || 0)" 
                :color="customColors"
              />
              <div class="health-label">健康评分</div>
            </div>
            
            <div class="stat-details">
              <div class="stat-row">
                <span class="label">预计 RUL</span>
                <span class="value" :class="(prediction.rul?.estimated_rul_days || 0) < 7 ? 'danger' : 'success'">
                  {{ prediction.rul?.estimated_rul_days || 0 }} 天
                </span>
              </div>
              <div class="stat-row">
                <span class="label">累计 TCO</span>
                <span class="value">¥{{ Math.round(prediction.tco?.total_cost_of_ownership || 0).toLocaleString() }}</span>
              </div>
            </div>

            <div v-if="prediction.symptoms?.length > 0" class="risk-section">
              <div class="risk-title">风险征兆识别:</div>
              <div v-for="(s, i) in prediction.symptoms" :key="i" class="symptom-tag">
                ⚠️ {{ s.title }}
              </div>
            </div>
          </div>
          <el-skeleton v-else :rows="5" animated />
        </el-card>

        <el-card shadow="never" class="mt-16">
          <template #header><div class="card-header">AI 学习记录</div></template>
          <div class="learning-stats">
             <div class="l-item"><strong>{{ draftCount }}</strong><span>新知识</span></div>
             <div class="l-item"><strong>12</strong><span>已获技能</span></div>
          </div>
        </el-card>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { ChatDotRound, CircleCheck, Reading } from '@element-plus/icons-vue'
import { equipmentApi, type EquipmentType } from '@/api/equipment'
import { agentApi, type ConversationResponse, type AgentKnowledge } from '@/api/agent'
import request from '@/api/request'
import { ElMessage } from 'element-plus'

// 状态
const activeMode = ref('chat')
const activeAuditTab = ref('repair')
const chatLoading = ref(false)
const userInput = ref('')
const currentConvId = ref<number | null>(null)
const auditForm = ref({ equipment_type_id: null })
const auditResult = ref<any>(null)

// 数据
const conversations = ref<ConversationResponse[]>([])
const messages = ref<any[]>([])
const knowledgeDrafts = ref<AgentKnowledge[]>([])
const draftCount = ref(0)
const equipmentTypes = ref<EquipmentType[]>([])
const prediction = ref<any>(null)

const customColors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]

// API 调用
async function handleSendChat() {
  if (!userInput.value.trim() || chatLoading.value) return
  const text = userInput.value
  userInput.value = ''
  messages.value.push({ role: 'user', content: text })
  scrollToBottom()
  
  chatLoading.value = true
  try {
    const res = await agentApi.chat({
      conversation_id: currentConvId.value || undefined,
      message: text
    })
    currentConvId.value = res.data.conversation_id
    messages.value.push({ 
      role: 'assistant', 
      content: res.data.reply,
      trace_id: res.data.trace_id 
    })
    setTimeout(loadDrafts, 3000)
    loadConversations()
  } catch (error: any) {
    console.error('Chat failed', error)
    const errorMsg = error.response?.data?.error || '请求失败，请稍后重试'
    messages.value.push({ 
      role: 'assistant', 
      content: `⚠️ 对不起，分析过程中出现错误：${errorMsg}。请检查您的网络连接或稍后再试。` 
    })
    ElMessage.error('对话请求失败')
  } finally {
    chatLoading.value = false
    scrollToBottom()
  }
}

async function handleRunRepairAudit() {
  if (!auditForm.value.equipment_type_id) return
  const loadingMsg = ElMessage.warning({
    message: 'AI 审计任务已启动，正在扫描历史记录...',
    duration: 0
  })
  try {
    const res = await agentApi.auditRepair({
      equipment_type_id: auditForm.value.equipment_type_id,
      time_range: { start_date: '2026-01-01', end_date: '2026-04-26' }
    })
    auditResult.value = res.data
    loadingMsg.close()
    ElMessage.success('审计完成')
  } catch (error: any) {
    loadingMsg.close()
    console.error('Audit failed', error)
    ElMessage.error('审计任务失败：' + (error.response?.data?.error || '服务响应超时'))
  }
}

async function loadConversation(id: number) {
  currentConvId.value = id
  try {
    const res = await agentApi.getConversation(id)
    messages.value = res.data.messages || []
    activeMode.value = 'chat'
  } catch (error) {
    console.error('Failed to load conversation', error)
  } finally {
    scrollToBottom()
  }
}

async function loadConversations() {
  try {
    const res = await agentApi.listConversations()
    conversations.value = res.data || []
  } catch (e) {
    console.error('Failed to load conversations', e)
  }
}

async function loadDrafts() {
  try {
    const res = await agentApi.listKnowledgeDrafts()
    if (res && res.data && Array.isArray(res.data)) {
      knowledgeDrafts.value = res.data.filter((k: any) => k.status === 'draft')
      draftCount.value = knowledgeDrafts.value.length
    }
  } catch (e) {
    console.error('Failed to load drafts', e)
  }
}

async function loadPrediction(id: number = 1) {
  try {
    const res = await agentApi.getEquipmentPrediction(id)
    prediction.value = res.data
  } catch (e) {
    console.error('Failed to load prediction', e)
  }
}

async function confirmKnowledge(id: string) {
  try {
    await request.put(`/agent/knowledge/${id}/status`, { status: 'confirmed' })
    ElMessage.success('知识已正式入库')
    loadDrafts()
  } catch (e) {
    ElMessage.error('确认失败')
  }
}

async function rejectKnowledge(id: string) {
  try {
    await request.put(`/agent/knowledge/${id}/status`, { status: 'rejected' })
    loadDrafts()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

function formatMessage(text: string) {
  return (text || '').replace(/\n/g, '<br>')
}

function scrollToBottom() {
  nextTick(() => {
    const box = document.querySelector('.chat-messages')
    if (box) box.scrollTop = box.scrollHeight
  })
}

function formatDate(d: string) { return d ? new Date(d).toLocaleDateString() : '' }

onMounted(async () => {
  try {
    const typesRes = await equipmentApi.getTypes()
    equipmentTypes.value = Array.isArray(typesRes.data) ? typesRes.data : (typesRes.data as any).items || []
  } catch (e) {
    console.error('Failed to load equipment types')
  }
  loadConversations()
  loadDrafts()
  loadPrediction(1)
})
</script>

<style scoped>
.management-assistant {
  height: calc(100vh - 110px);
  margin: -20px;
  background: var(--color-bg-secondary);
}
.cockpit-layout { display: flex; height: 100%; }
.sidebar-nav {
  width: 240px;
  background: var(--color-bg-primary);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  padding: 16px 0;
}
.nav-group { margin-bottom: 24px; }
.group-title { padding: 0 20px; font-size: 12px; color: var(--color-text-tertiary); text-transform: uppercase; margin-bottom: 8px; }
.nav-item {
  padding: 12px 20px; display: flex; align-items: center; gap: 12px; cursor: pointer; color: var(--color-text-secondary); transition: all 0.2s;
}
.nav-item:hover { background: var(--color-bg-tertiary); }
.nav-item.active { background: var(--color-primary-dim); color: var(--color-primary); border-right: 3px solid var(--color-primary); }
.history-group { flex: 1; overflow-y: auto; }
.history-item { padding: 10px 20px; cursor: pointer; border-bottom: 1px solid rgba(0,0,0,0.02); }
.history-title { font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.history-meta { font-size: 11px; color: var(--color-text-tertiary); }

.main-content { flex: 1; display: flex; flex-direction: column; background: var(--color-bg-tertiary); overflow: hidden; }
.chat-container { display: flex; flex-direction: column; height: 100%; }
.chat-messages { flex: 1; overflow-y: auto; padding: 30px; display: flex; flex-direction: column; gap: 24px; }
.message { display: flex; gap: 16px; max-width: 85%; }
.message.user { flex-direction: row-reverse; align-self: flex-end; }
.message .avatar { width: 36px; height: 36px; border-radius: 8px; background: var(--color-primary); color: white; display: flex; align-items: center; justify-content: center; font-weight: bold; flex-shrink: 0; }
.message.user .avatar { background: #6366f1; }
.message .content { background: var(--color-bg-primary); padding: 12px 16px; border-radius: 12px; box-shadow: 0 2px 8px rgba(0,0,0,0.05); }
.loading-tip { font-size: 13px; color: var(--color-primary); margin-bottom: 8px; font-weight: 500; }
.message.user .content { background: var(--color-primary); color: white; }
.msg-footer { margin-top: 8px; font-size: 11px; opacity: 0.6; }
.chat-input-area { padding: 20px 30px 30px; }
.input-wrapper { background: var(--color-bg-primary); border: 1px solid var(--color-border); border-radius: 12px; padding: 12px; }
.input-actions { display: flex; justify-content: flex-end; margin-top: 8px; }

.context-panel { width: 300px; background: var(--color-bg-primary); border-left: 1px solid var(--color-border); padding: 20px; overflow-y: auto; }
.stat-main { text-align: center; margin-bottom: 20px; }
.health-label { font-size: 14px; color: var(--color-text-secondary); margin-top: -10px; }
.stat-row { display: flex; justify-content: space-between; margin-bottom: 12px; font-size: 13px; }
.stat-row .value { font-weight: bold; }
.stat-row .value.success { color: var(--color-success); }
.stat-row .value.danger { color: var(--color-danger); }
.risk-section { margin-top: 20px; border-top: 1px solid var(--color-border); padding-top: 15px; }
.risk-title { font-size: 12px; color: var(--color-danger); font-weight: bold; margin-bottom: 10px; }
.symptom-tag { font-size: 12px; background: #fff1f0; border: 1px solid #ffa39e; color: #cf1322; padding: 4px 8px; border-radius: 4px; margin-bottom: 6px; }

.welcome-guide { text-align: center; margin-top: 8vh; }
.guide-chips { display: flex; flex-direction: column; align-items: center; gap: 10px; margin-top: 24px; }
.guide-chips .el-tag { cursor: pointer; padding: 10px 20px; font-size: 14px; }

.learning-stats { display: grid; grid-template-columns: 1fr 1fr; text-align: center; }
.l-item strong { display: block; font-size: 20px; color: var(--color-primary); }
.l-item span { font-size: 12px; color: var(--color-text-tertiary); }

.p-20 { padding: 20px; }
.mt-20 { margin-top: 20px; }
.mt-24 { margin-top: 24px; }
.max-w-400 { max-w: 400px; }
.result-card { background: #fdf6ec; border: 1px solid #faecd8; }
.summary-text { white-space: pre-wrap; font-family: inherit; font-size: 14px; line-height: 1.6; }
</style>
