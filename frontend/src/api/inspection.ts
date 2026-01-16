import request from './request'

// Types
export interface InspectionTemplate {
  id: number
  name: string
  equipment_type_id: number
  equipment_type_name?: string
  item_count: number
  created_at: string
}

export interface InspectionTemplateDetail extends InspectionTemplate {
  items: InspectionItem[]
}

export interface InspectionItem {
  id: number
  template_id: number
  template_name?: string
  name: string
  method?: string
  criteria?: string
  sequence_order: number
}

export interface InspectionTask {
  id: number
  equipment_id: number
  equipment_code?: string
  equipment_name?: string
  template_id: number
  template_name?: string
  assigned_to: number
  assignee_name?: string
  scheduled_date: string // YYYY-MM-DD format
  status: 'pending' | 'in_progress' | 'completed' | 'overdue'
  started_at?: string
  completed_at?: string
  latitude?: number
  longitude?: number
  item_count: number
  completed_count: number
}

export interface InspectionTaskListResponse {
  total: number
  items: InspectionTask[]
}

export interface InspectionRecord {
  id: number
  task_id: number
  item_id: number
  item_name?: string
  result: 'OK' | 'NG'
  remark?: string
  photo_url?: string
  created_at: string
}

export interface StartInspectionRequest {
  equipment_id: number
  qr_code: string
  latitude?: number
  longitude?: number
  timestamp: number
}

export interface StartInspectionResponse {
  task_id: number
  equipment_id: number
  equipment?: Equipment
  items?: InspectionItem[]
  started_at: string
}

export interface Equipment {
  id: number
  code: string
  name: string
  type_id: number
  type_name?: string
  workshop_id: number
  workshop_name?: string
  factory_id?: number
  factory_name?: string
  qr_code: string
  status: string
}

export interface CompleteInspectionRequest {
  task_id: number
  records: {
    item_id: number
    result: 'OK' | 'NG'
    remark?: string
    photo_url?: string
  }[]
  latitude?: number
  longitude?: number
}

export interface CompleteInspectionResponse {
  task_id: number
  completed_at: string
  total_count: number
  ok_count: number
  ng_count: number
  ng_items: number[]
}

export interface InspectionStatistics {
  total_tasks: number
  pending_tasks: number
  in_progress_tasks: number
  completed_tasks: number
  overdue_tasks: number
  today_completed: number
  completion_rate: number
}

export interface MyTasksStatistics {
  pending_count: number
  in_progress_count: number
  today_tasks: number
}

// Template APIs
export const inspectionTemplateApi = {
  getTemplates: () => request.get<InspectionTemplate[]>('/inspection/templates'),
  getTemplate: (id: number) => request.get<InspectionTemplateDetail>(`/inspection/templates/${id}`),
  createTemplate: (data: { name: string; equipment_type_id: number }) =>
    request.post<InspectionTemplate>('/inspection/templates', data),
}

// Item APIs
export const inspectionItemApi = {
  createItem: (data: {
    template_id: number
    name: string
    method?: string
    criteria?: string
    sequence_order: number
  }) => request.post<InspectionItem>('/inspection/items', data),
}

// Task APIs
export const inspectionTaskApi = {
  getTasks: (params: {
    page?: number
    page_size?: number
    assigned_to?: number
    status?: string
    date_from?: string
    date_to?: string
  }) => request.get<InspectionTaskListResponse>('/inspection/tasks', { params }),

  getTask: (id: number) => request.get<InspectionTask>(`/inspection/tasks/${id}`),

  getMyTasks: () => request.get<InspectionTask[]>('/inspection/my-tasks'),

  getMyStats: () => request.get<MyTasksStatistics>('/inspection/my-stats'),

  start: (data: StartInspectionRequest) =>
    request.post<StartInspectionResponse>('/inspection/start', data),

  complete: (data: CompleteInspectionRequest) =>
    request.post<CompleteInspectionResponse>('/inspection/complete', data),

  getStatistics: () => request.get<InspectionStatistics>('/inspection/statistics'),
}
