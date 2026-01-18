import request from './request'

// Types
export interface Base {
  id: number
  code: string
  name: string
  created_at: string
  updated_at: string
}

export interface Factory {
  id: number
  base_id: number
  base_name?: string
  code: string
  name: string
  created_at: string
  updated_at: string
}

export interface Workshop {
  id: number
  factory_id: number
  factory_name?: string
  code: string
  name: string
  created_at: string
  updated_at: string
}

export interface EquipmentType {
  id: number
  name: string
  category?: string
  inspection_template_id?: number
  created_at: string
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
  spec?: string
  purchase_date?: string
  status: 'running' | 'stopped' | 'maintenance' | 'scrapped'
  dedicated_maintenance_id?: number
  dedicated_maintenance_name?: string
  created_at: string
  updated_at: string
}

export interface EquipmentListResponse {
  total: number
  items: Equipment[]
}

export interface EquipmentStatistics {
  total: number
  running: number
  stopped: number
  maintenance: number
  scrapped: number
}

// Organization APIs
export const orgApi = {
  // Bases
  getBases: () => request.get<Base[]>('/organization/bases'),
  createBase: (data: { code: string; name: string }) => request.post<Base>('/organization/bases', data),
  updateBase: (id: number, data: { code: string; name: string }) => request.put(`/organization/bases/${id}`, data),
  deleteBase: (id: number) => request.delete(`/organization/bases/${id}`),

  // Factories
  getFactories: (baseId?: number) =>
    request.get<Factory[]>('/organization/factories', { params: baseId ? { base_id: baseId } : {} }),
  createFactory: (data: { base_id: number; code: string; name: string }) =>
    request.post<Factory>('/organization/factories', data),
  updateFactory: (id: number, data: { base_id: number; code: string; name: string }) =>
    request.put(`/organization/factories/${id}`, data),
  deleteFactory: (id: number) => request.delete(`/organization/factories/${id}`),

  // Workshops
  getWorkshops: (factoryId?: number) =>
    request.get<Workshop[]>('/organization/workshops', {
      params: factoryId ? { factory_id: factoryId } : {},
    }),
  createWorkshop: (data: { factory_id: number; code: string; name: string }) =>
    request.post<Workshop>('/organization/workshops', data),
  updateWorkshop: (id: number, data: { factory_id: number; code: string; name: string }) =>
    request.put(`/organization/workshops/${id}`, data),
  deleteWorkshop: (id: number) => request.delete(`/organization/workshops/${id}`),
}

// Equipment Type APIs
export const equipmentTypeApi = {
  getTypes: () => request.get<EquipmentType[]>('/equipment/types'),
  createType: (data: { name: string; category?: string }) =>
    request.post<EquipmentType>('/equipment/types', data),
  updateType: (id: number, data: { name: string; category?: string }) =>
    request.put(`/equipment/types/${id}`, data),
  deleteType: (id: number) => request.delete(`/equipment/types/${id}`),
}

// Equipment APIs
export const equipmentApi = {
  getList: (params: {
    page?: number
    page_size?: number
    code?: string
    name?: string
    type_id?: number
    factory_id?: number
    workshop_id?: number
    status?: string
  }) => request.get<EquipmentListResponse>('/equipment', { params }),

  getById: (id: number) => request.get<Equipment>(`/equipment/${id}`),

  getByQRCode: (code: string) => request.get<Equipment>(`/equipment/qr/${code}`),

  create: (data: {
    code: string
    name: string
    type_id: number
    workshop_id: number
    spec?: string
    purchase_date?: string
    status?: string
    dedicated_maintenance_id?: number
  }) => request.post<Equipment>('/equipment', data),

  update: (
    id: number,
    data: {
      code: string
      name: string
      type_id: number
      workshop_id: number
      spec?: string
      purchase_date?: string
      status?: string
      dedicated_maintenance_id?: number
    }
  ) => request.put<Equipment>(`/equipment/${id}`, data),

  delete: (id: number) => request.delete(`/equipment/${id}`),

  getStatistics: () => request.get<EquipmentStatistics>('/equipment/statistics'),
  getTypes: equipmentTypeApi.getTypes,
}
