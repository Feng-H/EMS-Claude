import request from './request'

// =====================================================
// Types
// =====================================================

export enum MaintenanceLevel {
  Level1 = 1,  // 一级保养
  Level2 = 2,  // 二级保养
  Level3 = 3   // 精度保养
}

export enum MaintenanceTaskStatus {
  Pending = 'pending',
  InProgress = 'in_progress',
  Completed = 'completed',
  Overdue = 'overdue'
}

export interface MaintenancePlan {
  id: number
  name: string
  equipment_type_id: number
  equipment_type_name?: string
  level: number
  level_name?: string
  cycle_days: number
  flexible_days: number
  work_hours: number
  item_count: number
  created_at: string
}

export interface CreateMaintenancePlanRequest {
  name: string
  equipment_type_id: number
  level: number
  cycle_days: number
  flexible_days?: number
  work_hours?: number
}

export interface MaintenanceItem {
  id: number
  plan_id: number
  name: string
  method?: string
  criteria?: string
  sequence_order: number
}

export interface CreateMaintenanceItemRequest {
  plan_id: number
  name: string
  method?: string
  criteria?: string
  sequence_order: number
}

export interface MaintenanceTask {
  id: number
  plan_id: number
  plan_name?: string
  equipment_id: number
  equipment_code?: string
  equipment_name?: string
  assigned_to?: number
  assignee_name?: string
  scheduled_date: string
  due_date: string
  status: MaintenanceTaskStatus
  started_at?: string
  completed_at?: string
  actual_hours: number
  remark?: string
  item_count: number
  completed_count: number
  created_at: string
}

export interface MaintenanceTaskQuery {
  status?: string
  assigned_to?: number
  date_from?: string
  date_to?: string
  page?: number
  page_size?: number
}

export interface MaintenanceTaskListResponse {
  total: number
  items: MaintenanceTask[]
}

export interface GenerateMaintenanceTasksRequest {
  plan_id: number
  equipment_ids: number[]
  date: string  // YYYY-MM-DD
}

export interface GenerateMaintenanceTasksResponse {
  created_count: number
  task_ids: number[]
  errors: string[]
}

export interface StartMaintenanceRequest {
  task_id: number
  latitude?: number
  longitude?: number
}

export interface MaintenanceItemRecord {
  item_id: number
  result: string  // OK/NG
  remark?: string
  photo_url?: string
}

export interface CompleteMaintenanceRequest {
  task_id: number
  records: MaintenanceItemRecord[]
  latitude?: number
  longitude?: number
  actual_hours?: number
  remark?: string
}

export interface CompleteMaintenanceResponse {
  task_id: number
  completed_at: string
  total_count: number
  ok_count: number
  ng_count: number
  ng_item_ids: number[]
}

export interface MaintenanceStatistics {
  total_plans: number
  total_tasks: number
  pending_tasks: number
  in_progress_tasks: number
  completed_tasks: number
  overdue_tasks: number
  today_completed: number
  completion_rate: number
}

// =====================================================
// API Functions
// =====================================================

// Plan Management
export const getMaintenancePlans = () => {
  return request.get<MaintenancePlan[]>('/maintenance/plans')
}

export const createMaintenancePlan = (data: CreateMaintenancePlanRequest) => {
  return request.post<MaintenancePlan>('/maintenance/plans', data)
}

export const createMaintenanceItem = (data: CreateMaintenanceItemRequest) => {
  return request.post<MaintenanceItem>('/maintenance/items', data)
}

// Task Management
export const generateMaintenanceTasks = (data: GenerateMaintenanceTasksRequest) => {
  return request.post<GenerateMaintenanceTasksResponse>('/maintenance/tasks/generate', data)
}

export const getMaintenanceTasks = (params: MaintenanceTaskQuery) => {
  return request.get<MaintenanceTaskListResponse>('/maintenance/tasks', { params })
}

export const getMaintenanceTask = (id: number) => {
  return request.get<MaintenanceTask>(`/maintenance/tasks/${id}`)
}

export const getMyMaintenanceTasks = () => {
  return request.get<MaintenanceTask[]>('/maintenance/my-tasks')
}

// Execution
export const startMaintenance = (data: StartMaintenanceRequest) => {
  return request.post('/maintenance/start', data)
}

export const completeMaintenance = (data: CompleteMaintenanceRequest) => {
  return request.post<CompleteMaintenanceResponse>('/maintenance/complete', data)
}

// Statistics
export const getMaintenanceStatistics = () => {
  return request.get<MaintenanceStatistics>('/maintenance/statistics')
}

// =====================================================
// Helper Functions
// =====================================================

export const getLevelName = (level: number): string => {
  const names: Record<number, string> = {
    1: '一级保养',
    2: '二级保养',
    3: '精度保养'
  }
  return names[level] || '未知'
}

export const getTaskStatusName = (status: MaintenanceTaskStatus): string => {
  const names: Record<MaintenanceTaskStatus, string> = {
    [MaintenanceTaskStatus.Pending]: '待执行',
    [MaintenanceTaskStatus.InProgress]: '进行中',
    [MaintenanceTaskStatus.Completed]: '已完成',
    [MaintenanceTaskStatus.Overdue]: '已逾期'
  }
  return names[status] || '未知'
}

export const getTaskStatusType = (status: MaintenanceTaskStatus): string => {
  const types: Record<MaintenanceTaskStatus, string> = {
    [MaintenanceTaskStatus.Pending]: 'info',
    [MaintenanceTaskStatus.InProgress]: 'warning',
    [MaintenanceTaskStatus.Completed]: 'success',
    [MaintenanceTaskStatus.Overdue]: 'danger'
  }
  return types[status] || 'info'
}

// Alias exports for backward compatibility
export const getStatusName = getTaskStatusName
export const getStatusType = getTaskStatusType

// API object for unified access
export const maintenanceApi = {
  getPlans: getMaintenancePlans,
  createPlan: createMaintenancePlan,
  createItem: createMaintenanceItem,
  generateTasks: generateMaintenanceTasks,
  getTasks: getMaintenanceTasks,
  getTask: getMaintenanceTask,
  getMyTasks: getMyMaintenanceTasks,
  start: startMaintenance,
  complete: completeMaintenance,
  getStatistics: getMaintenanceStatistics
}
