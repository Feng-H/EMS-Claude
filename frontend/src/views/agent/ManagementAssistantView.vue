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
          <div 
            class="nav-item" 
            :class="{ active: activeMode === 'skill' }"
            @click="activeMode = 'skill'"
          >
            <el-icon><GoldMedal /></el-icon>
            <span>技能管理</span>
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
              <h2>我是您的 EMS 智能专家</h2>
              <p>您可以直接问我任何关于设备管理的问题，例如：</p>
              <div class="guide-chips">
                <el-tag @click="userInput = 'CNC-001 最近 30 天维修费为什么超标？'">CNC-001 维修费分析</el-tag>
                <el-tag @click="userInput = '对比张三和李四负责区域的保养效果'">人效价值对标</el-tag>
                <el-tag @click="userInput = '帮我核查最近是否有级联失效风险'">级联失效审计</el-tag>
              </div>
            </div>
            
            <div v-for="(msg, idx) in messages" :key="idx" :class="['message', msg.role]">
              <div class="avatar">{{ msg.role === 'user' ? 'U' : 'AI' }}</div>
              <div class="content">
                <div class="text">{{ msg.content }}</div>
                <div v-if="msg.trace_id" class="msg-footer">Trace: {{ msg.trace_id }}</div>
              </div>
            </div>
            <div v-if="chatLoading" class="message assistant loading">
              <div class="avatar">AI</div>
              <div class="content"><el-skeleton :rows="2" animated /></div>
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

        <!-- 模式2：专项审计 (原来的 Phase 1 功能) -->
        <div v-if="activeMode === 'audit'" class="audit-mode">
          <div class="audit-config">
            <el-tabs v-model="activeAuditTab">
              <el-tab-pane label="维修合理性审计" name="repair">
                 <el-form label-position="top">
                    <el-form-item label="设备类型">
                      <el-select v-model="maintenanceForm.equipment_type_id" style="width: 100%">
                        <el-option v-for="t in equipmentTypes" :key="t.id" :label="t.name" :value="t.id" />
                      </el-select>
                    </el-form-item>
                    <el-button type="warning" block @click="handleRunRepairAudit">启动 AI 深度审计</el-button>
                 </el-form>
              </el-tab-pane>
              <el-tab-pane label="保养周期优化" name="maintenance">
                 <!-- 保养优化表单 -->
              </el-tab-pane>
            </el-tabs>
          </div>
          <div class="audit-result">
             <div v-if="result" class="result-card">
                <div class="summary-box">{{ result.summary }}</div>
                <!-- 证据展示等 -->
             </div>
          </div>
        </div>

        <!-- 模式3：知识审核 (New!) -->
        <div v-if="activeMode === 'knowledge'" class="knowledge-audit">
          <div class="section-header">
             <h3>知识库待审核 (AI 自主提炼)</h3>
             <p>Agent 在对话过程中识别到的高价值工业经验，请您校准入库。</p>
          </div>
          <el-table :data="knowledgeDrafts" stripe>
            <el-table-column prop="title" label="标题" />
            <el-table-column prop="type" label="类型" width="120" />
            <el-table-column prop="confidence" label="置信度" width="100">
              <template #default="{row}">
                <el-progress :percentage="row.confidence * 100" :format="() => ''" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{row}">
                <el-button type="success" link @click="confirmKnowledge(row.id)">确认入库</el-button>
                <el-button type="danger" link @click="rejectKnowledge(row.id)">拒绝</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </main>

      <!-- 右侧：上下文面板 -->
      <aside class="context-panel">
        <el-card shadow="never">
          <template #header><div class="card-header">设备实时工况</div></template>
          <div class="stat-row">
            <span class="label">当前运行机台</span>
            <span class="value">42 / 50</span>
          </div>
          <div class="stat-row">
            <span class="label">平均负载</span>
            <span class="value success">85%</span>
          </div>
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
import { ChatDotRound, CircleCheck, Reading, GoldMedal, Timer } from '@element-plus/icons-vue'
import { equipmentApi, orgApi, type EquipmentType, type Factory } from '@/api/equipment'
import { agentApi, type ChatResponse, type ConversationResponse, type AgentKnowledge } from '@/api/agent'
import { ElMessage } from 'element-plus'

// 状态管理
const activeMode = ref('chat')
const activeAuditTab = ref('repair')
const loading = ref(false)
const chatLoading = ref(false)
const userInput = ref('')
const currentConvId = ref<number | null>(null)

// 数据
const conversations = ref<ConversationResponse[]>([])
const messages = ref<any[]>([])
const knowledgeDrafts = ref<AgentKnowledge[]>([])
const draftCount = ref(0)
const equipmentTypes = ref<EquipmentType[]>([])
const result = ref<any>(null)

// 保养表单（沿用旧的）
const maintenanceForm = ref({
  equipment_type_id: null as number | null,
  dateRange: [] as string[]
})

// 处理函数
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
    
    // 每次对话后，由于后端是异步学习，我们稍等一下刷新草稿箱
    setTimeout(loadDrafts, 3000)
    await loadConversations()
  } catch (error: any) {
    ElMessage.error('对话失败')
  } finally {
    chatLoading.value = false
    scrollToBottom()
  }
}

async function loadConversation(id: number) {
  currentConvId.value = id
  loading.value = true
  try {
    const res = await agentApi.getConversation(id)
    messages.value = res.data.messages || []
    activeMode.value = 'chat'
  } finally {
    loading.value = false
    scrollToBottom()
  }
}

async function loadConversations() {
  const res = await agentApi.listConversations()
  conversations.value = res.data
}

async function loadDrafts() {
  const res = await agentApi.listKnowledgeDrafts()
  // 模拟过滤 draft 状态
  knowledgeDrafts.value = res.data.filter(k => k.status === 'draft')
  draftCount.value = knowledgeDrafts.value.length
}

function scrollToBottom() {
  nextTick(() => {
    const box = document.querySelector('.chat-messages')
    if (box) box.scrollTop = box.scrollHeight
  })
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString()
}

// 模拟旧功能
async function handleRunRepairAudit() {
  ElMessage.success('正在执行深度审计...')
  // 原有的 API 调用逻辑...
}

onMounted(async () => {
  try {
    const [typesRes, factoriesRes] = await Promise.all([
      equipmentApi.getTypes(),
      orgApi.getFactories(),
      loadConversations(),
      loadDrafts()
    ])
    equipmentTypes.value = Array.isArray(typesRes.data) ? typesRes.data : (typesRes.data as any).items || []
  } catch (error) {
    console.error('Failed to load cockpit data')
  }
})
</script>

<style scoped>
.management-assistant {
  height: calc(100vh - 110px);
  margin: -20px;
  background: var(--color-bg-secondary);
}

.cockpit-layout {
  display: flex;
  height: 100%;
}

/* 侧边栏导航 */
.sidebar-nav {
  width: 240px;
  background: var(--color-bg-primary);
  border-right: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  padding: 16px 0;
}

.nav-group {
  margin-bottom: 24px;
}

.group-title {
  padding: 0 20px;
  font-size: 12px;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 8px;
}

.nav-item {
  padding: 12px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}

.nav-item:hover { background: var(--color-bg-tertiary); }
.nav-item.active {
  background: var(--color-primary-dim);
  color: var(--color-primary);
  border-right: 3px solid var(--color-primary);
}

.history-group {
  flex: 1;
  overflow-y: auto;
}

.history-item {
  padding: 10px 20px;
  cursor: pointer;
}

.history-title {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.history-meta {
  font-size: 11px;
  color: var(--color-text-tertiary);
}

/* 主内容区 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--color-bg-tertiary);
  overflow: hidden;
}

/* 聊天容器 */
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 30px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.welcome-guide {
  text-align: center;
  margin-top: 10vh;
}

.guide-icon { font-size: 48px; margin-bottom: 16px; }

.guide-chips {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.guide-chips .el-tag { cursor: pointer; }

.message {
  display: flex;
  gap: 16px;
  max-width: 85%;
}

.message.user {
  flex-direction: row-reverse;
  align-self: flex-end;
}

.message .avatar {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  flex-shrink: 0;
}

.message.user .avatar { background: #6366f1; }

.message .content {
  background: var(--color-bg-primary);
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.message.user .content {
  background: var(--color-primary);
  color: white;
}

.message.assistant .content {
  border-bottom-left-radius: 2px;
}

.msg-footer {
  margin-top: 8px;
  font-size: 11px;
  opacity: 0.6;
}

/* 输入区 */
.chat-input-area {
  padding: 20px 30px 30px;
  background: var(--color-bg-tertiary);
}

.input-wrapper {
  background: var(--color-bg-primary);
  border: 1px solid var(--color-border);
  border-radius: 12px;
  padding: 12px;
  box-shadow: 0 -4px 12px rgba(0,0,0,0.02);
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

/* 右侧上下文 */
.context-panel {
  width: 280px;
  background: var(--color-bg-primary);
  border-left: 1px solid var(--color-border);
  padding: 20px;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
}

.stat-row .value { font-weight: bold; }
.stat-row .value.success { color: var(--color-success); }

.learning-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  text-align: center;
}

.l-item strong { display: block; font-size: 20px; color: var(--color-primary); }
.l-item span { font-size: 12px; color: var(--color-text-tertiary); }

.mt-16 { margin-top: 16px; }
</style>
