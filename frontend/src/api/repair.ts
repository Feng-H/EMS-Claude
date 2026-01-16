import request from './request'

// Types
export type RepairStatus = 'pending' | 'assigned' | 'in_progress' | 'testing' | 'confirmed' | 'audited' | 'closed'

export interface RepairOrder {
  id: number
  equipment_id: number
  equipment_code?: string
  equipment_name?: string
  fault_description: string
  fault_code?: string
  reporter_id: number
  reporter_name?: string
  assigned_to?: number
  assignee_name?: string
  status: RepairStatus
  priority: number // 1=高,2=中,3=低
  photos?: string[]
  solution?: string
  spare_parts?: string
  actual_hours?: number
  created_at: string
  started_at?: string
  completed_at?: string
  confirmed_at?: string
  audited_at?: string
}

export interface RepairOrderListResponse {
  total: number
  items: RepairOrder[]
}

export interface RepairOrderDetail extends RepairOrder {
  logs?: RepairLog[]
}

export interface RepairLog {
  id: number
  order_id: number
  user_id: number
  user_name?: string
  action: string
  content: string
  created_at: string
}

export interface CreateRepairRequest {
  equipment_id: number
  fault_description: string
  fault_code?: string
  photos?: string[]
  priority?: number
}

export interface AssignRepairRequest {
  assign_to: number
}

export interface StartRepairRequest {
  latitude?: number
  longitude?: number
}

export interface UpdateRepairRequest {
  solution?: string
  spare_parts?: string
  actual_hours?: number
  photos?: string[]
  next_status?: string // testing, confirmed
}

export interface ConfirmRepairRequest {
  accepted: boolean
  comment?: string
  photos?: string[]
}

export interface AuditRepairRequest {
  approved: boolean
  comment?: string
  actual_hours?: number
}

export interface RepairStatistics {
  total_orders: number
  pending_orders: number
  in_progress_orders: number
  completed_orders: number
  today_completed: number
  today_created: number
  avg_repair_time: number
  avg_response_time: number
}

export interface MyRepairStatistics {
  pending_count: number
  in_progress_count: number
  completed_count: number
  today_completed: number
}

// Repair Order APIs
export const repairOrderApi = {
  getOrders: (params: {
    page?: number
    page_size?: number
    status?: string
    priority?: number
    assigned_to?: number
    date_from?: string
    date_to?: string
  }) => request.get<RepairOrderListResponse>('/repair/orders', { params }),

  getOrder: (id: number) => request.get<RepairOrderDetail>(`/repair/orders/${id}`),

  createOrder: (data: CreateRepairRequest) =>
    request.post<RepairOrder>('/repair/orders', data),

  assignOrder: (id: number, data: AssignRepairRequest) =>
    request.post<RepairOrder>(`/repair/orders/${id}/assign`, data),

  startRepair: (id: number, data: StartRepairRequest) =>
    request.post(`/repair/orders/${id}/start`, data),

  updateRepair: (id: number, data: UpdateRepairRequest) =>
    request.post(`/repair/orders/${id}/update`, data),

  confirmRepair: (id: number, data: ConfirmRepairRequest) =>
    request.post(`/repair/orders/${id}/confirm`, data),

  auditRepair: (id: number, data: AuditRepairRequest) =>
    request.post(`/repair/orders/${id}/audit`, data),

  getMyTasks: () => request.get<RepairOrder[]>('/repair/my-tasks'),

  getMyStats: () => request.get<MyRepairStatistics>('/repair/my-stats'),

  getStatistics: () => request.get<RepairStatistics>('/repair/statistics'),
}

// Helper functions
export const getRepairStatusText = (status: RepairStatus): string => {
  const map: Record<RepairStatus, string> = {
    pending: '待派单',
    assigned: '已派单',
    in_progress: '维修中',
    testing: '待测试',
    confirmed: '待审核',
    audited: '已审核',
    closed: '已关闭',
  }
  return map[status] || status
}

export const getRepairStatusType = (status: RepairStatus): string => {
  const map: Record<RepairStatus, string> = {
    pending: 'info',
    assigned: 'primary',
    in_progress: 'warning',
    testing: 'warning',
    confirmed: 'success',
    audited: 'success',
    closed: 'info',
  }
  return map[status] || 'info'
}

export const getPriorityText = (priority: number): string => {
  const map: Record<number, string> = {
    1: '高',
    2: '中',
    3: '低',
  }
  return map[priority] || '中'
}

export const getPriorityType = (priority: number): string => {
  const map: Record<number, string> = {
    1: 'danger',
    2: 'warning',
    3: 'info',
  }
  return map[priority] || 'info'
}

// Alias exports for backward compatibility
export const getStatusType = getRepairStatusType
export const getStatusText = getRepairStatusText
