import request from './request'

export interface TimeRange {
  start_date: string
  end_date: string
}

export interface ChatRequest {
  conversation_id?: number
  message: string
  context?: any
  system_prompt?: string
}

export interface ChatResponse {
  conversation_id: number
  reply: string
  trace_id: string
  artifact_id?: number
  suggested_actions?: string[]
}

export interface ConversationResponse {
  id: number
  title: string
  status: string
  created_at: string
  updated_at: string
}

export interface AgentKnowledge {
  id: string
  title: string
  type: string
  summary: string
  confidence: number
  status: string
  created_at: string
}

export interface MaintenanceRecommendRequest {
  factory_id?: number
  workshop_id?: number
  equipment_type_id: number
  equipment_ids?: number[]
  time_range: TimeRange
  question?: string
  language?: string
  system_prompt?: string
}

export interface RepairAuditRequest {
  factory_id?: number
  workshop_id?: number
  equipment_type_id: number
  time_range: TimeRange
  anomaly_types?: string[]
  language?: string
  system_prompt?: string
}

export interface AgentResponse<T> {
  success: boolean
  trace_id: string
  language: string
  scenario: string
  scope_summary: any
  summary: string
  risk_level: string
  artifact_id?: number
  evidence_count: number
  data: T
}

export const agentApi = {
  // 对话
  chat: (data: ChatRequest) => 
    request.post<ChatResponse>('/agent/chat', data),
  
  listConversations: () => 
    request.get<ConversationResponse[]>('/agent/conversations'),
    
  getConversation: (id: number) => 
    request.get<any>(`/agent/conversations/${id}`),

  // 专项审计
  recommendMaintenance: (data: MaintenanceRecommendRequest) => 
    request.post<AgentResponse<any>>('/agent/maintenance/recommend', data),
    
  auditRepair: (data: RepairAuditRequest) => 
    request.post<AgentResponse<any>>('/agent/audit/repair', data),
    
  // 知识与技能
  listSkills: (status?: string) => 
    request.get<any[]>('/agent/skills', { params: { status } }),

  listKnowledgeDrafts: () => 
    request.get<AgentKnowledge[]>('/knowledge'), // 复用知识库接口但过滤状态
    
  // 历史记录
  listSessions: () => 
    request.get<any[]>('/agent/sessions'),
    
  getSession: (id: number) => 
    request.get<any>(`/agent/sessions/${id}`),
    
  getArtifact: (id: number) => 
    request.get<any>(`/agent/artifacts/${id}`),

  // 预测性分析 (Phase 3 补充)
  getEquipmentPrediction: (equipmentId: number) =>
    request.get<any>(`/agent/equipment/${equipmentId}/prediction`),
}
