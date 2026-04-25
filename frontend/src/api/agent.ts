import request from './request'

export interface TimeRange {
  start_date: string
  end_date: string
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

export interface EvidenceItem {
  evidence_type: string
  source_table: string
  source_id: number
  title: string
  excerpt: string
  score: number
}

export interface RecommendationItem {
  type: string
  target: string
  target_id: number
  title: string
  description: string
  reason: string
  impact: string
}

export interface AnomalyItem {
  anomaly_type: string
  severity: string
  target_type: string
  target_id: number
  title: string
  description: string
  suggested_action: string
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
  recommendMaintenance: (data: MaintenanceRecommendRequest) => 
    request.post<AgentResponse<any>>('/agent/maintenance/recommend', data),
    
  auditRepair: (data: RepairAuditRequest) => 
    request.post<AgentResponse<any>>('/agent/audit/repair', data),
    
  listSessions: () => 
    request.get<any[]>('/agent/sessions'),
    
  getSession: (id: number) => 
    request.get<any>(`/agent/sessions/${id}`),
    
  getArtifact: (id: number) => 
    request.get<any>(`/agent/artifacts/${id}`),
}
