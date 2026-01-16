import request from './request'

// =====================================================
// Types
// =====================================================

export interface SparePart {
  id: number
  code: string
  name: string
  specification?: string
  unit?: string
  factory_id?: number
  factory_name?: string
  safety_stock: number
  current_stock?: number
  created_at: string
}

export interface CreateSparePartRequest {
  code: string
  name: string
  specification?: string
  unit?: string
  factory_id?: number
  safety_stock?: number
}

export interface UpdateSparePartRequest {
  code: string
  name: string
  specification?: string
  unit?: string
  factory_id?: number
  safety_stock?: number
}

export interface SparePartListResponse {
  total: number
  items: SparePart[]
}

export interface Inventory {
  id: number
  spare_part_id: number
  spare_part_code: string
  spare_part_name: string
  factory_id: number
  factory_name: string
  quantity: number
  is_low_stock: boolean
  updated_at: string
}

export interface InventoryListResponse {
  total: number
  items: Inventory[]
}

export interface StockInRequest {
  spare_part_id: number
  factory_id: number
  quantity: number
  remark?: string
}

export interface StockOutRequest {
  spare_part_id: number
  factory_id: number
  quantity: number
  order_id?: number
  task_id?: number
  remark?: string
}

export interface LowStockAlert {
  spare_part_id: number
  spare_part_code: string
  spare_part_name: string
  factory_id: number
  factory_name: string
  current_stock: number
  safety_stock: number
  shortage: number
}

export interface Consumption {
  id: number
  spare_part_id: number
  spare_part_code: string
  spare_part_name: string
  order_id?: number
  task_id?: number
  quantity: number
  user_id: number
  user_name: string
  remark?: string
  created_at: string
}

export interface ConsumptionListResponse {
  total: number
  items: Consumption[]
}

export interface CreateConsumptionRequest {
  spare_part_id: number
  quantity: number
  order_id?: number
  task_id?: number
  remark?: string
}

export interface SparePartStatistics {
  total_parts: number
  low_stock_count: number
  total_stock_value: number
  monthly_consumption: number
}

// =====================================================
// API Functions
// =====================================================

// Part Management
export const getSpareParts = (params: {
  code?: string
  name?: string
  page?: number
  page_size?: number
}) => {
  return request.get<SparePartListResponse>('/spareparts', { params })
}

export const createSparePart = (data: CreateSparePartRequest) => {
  return request.post<SparePart>('/spareparts', data)
}

export const updateSparePart = (id: number, data: UpdateSparePartRequest) => {
  return request.put(`/spareparts/${id}`, data)
}

export const deleteSparePart = (id: number) => {
  return request.delete(`/spareparts/${id}`)
}

// Inventory
export const getInventory = (params: {
  spare_part_id?: number
  factory_id?: number
  low_stock?: boolean
  page?: number
  page_size?: number
}) => {
  return request.get<InventoryListResponse>('/spareparts/inventory', { params })
}

export const stockIn = (data: StockInRequest) => {
  return request.post('/spareparts/stock-in', data)
}

export const stockOut = (data: StockOutRequest) => {
  return request.post('/spareparts/stock-out', data)
}

export const getLowStockAlerts = () => {
  return request.get<LowStockAlert[]>('/spareparts/alerts')
}

// Consumption
export const getConsumptions = (params: {
  spare_part_id?: number
  order_id?: number
  task_id?: number
  date_from?: string
  date_to?: string
  page?: number
  page_size?: number
}) => {
  return request.get<ConsumptionListResponse>('/spareparts/consumptions', { params })
}

export const createConsumption = (data: CreateConsumptionRequest) => {
  return request.post('/spareparts/consumptions', data)
}

// Statistics
export const getSparePartStatistics = () => {
  return request.get<SparePartStatistics>('/spareparts/statistics')
}
