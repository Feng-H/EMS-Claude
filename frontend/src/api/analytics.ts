import request from './request'

// =====================================================
// Types
// =====================================================

export interface EquipmentAnalytics {
  total_equipment: number
  running_equipment: number
  stopped_equipment: number
  maintenance_equipment: number
  scrapped_equipment: number
}

export interface MTTRMTBF {
  mttr: number
  mtbf: number
  availability: number
}

export interface CompletionRate {
  inspection_completion_rate: number
  maintenance_completion_rate: number
  repair_completion_rate: number
}

export interface DashboardOverview {
  equipment: EquipmentAnalytics
  mttr_mtbf: MTTRMTBF
  tasks: CompletionRate
  pending_inspections: number
  pending_maintenances: number
  pending_repairs: number
  low_stock_alerts: number
}

export interface TrendData {
  date: string
  inspection_tasks: number
  maintenance_tasks: number
  repair_orders: number
  downtime_hours: number
}

export interface FailureAnalysis {
  equipment_type_id: number
  equipment_type_name: string
  failure_count: number
  total_downtime: number
}

export interface TopFailureEquipment {
  equipment_id: number
  equipment_code: string
  equipment_name: string
  failure_count: number
  downtime_hours: number
  mttr: number
}

// =====================================================
// API Functions
// =====================================================

export const getDashboardOverview = () => {
  return request.get<DashboardOverview>('/analytics/dashboard')
}

export const getMTTRMTBF = (params?: { factory_id?: number }) => {
  return request.get<MTTRMTBF>('/analytics/mttr-mtbf', { params })
}

export const getTrendData = (params: {
  start_date?: string
  end_date?: string
}) => {
  return request.get<TrendData[]>('/analytics/trends', { params })
}

export const getFailureAnalysis = (params?: { limit?: number }) => {
  return request.get<FailureAnalysis[]>('/analytics/failures', { params })
}

export const getTopFailureEquipment = (params?: { limit?: number }) => {
  return request.get<TopFailureEquipment[]>('/analytics/top-failures', { params })
}
