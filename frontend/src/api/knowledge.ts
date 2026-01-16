import request from './request'

// =====================================================
// Types
// =====================================================

export interface KnowledgeArticle {
  id: number
  title: string
  equipment_type_id?: number
  equipment_type_name?: string
  fault_phenomenon?: string
  cause_analysis?: string
  solution: string
  source_type: string
  source_id?: number
  tags: string[]
  created_by: number
  creator_name?: string
  created_at: string
  updated_at: string
}

export interface CreateKnowledgeArticleRequest {
  title: string
  equipment_type_id?: number
  fault_phenomenon?: string
  cause_analysis?: string
  solution: string
  source_type?: string
  source_id?: number
  tags?: string[]
}

export interface UpdateKnowledgeArticleRequest {
  title: string
  equipment_type_id?: number
  fault_phenomenon?: string
  cause_analysis?: string
  solution: string
  tags?: string[]
}

export interface KnowledgeArticleListResponse {
  total: number
  items: KnowledgeArticle[]
}

export interface ConvertFromRepairRequest {
  order_id: number
  title: string
  fault_phenomenon: string
  cause_analysis: string
  tags: string[]
}

// =====================================================
// API Functions
// =====================================================

export const getKnowledgeArticles = (params: {
  keyword?: string
  equipment_type_id?: number
  tag?: string
  source_type?: string
  page?: number
  page_size?: number
}) => {
  return request.get<KnowledgeArticleListResponse>('/knowledge', { params })
}

export const getKnowledgeArticle = (id: number) => {
  return request.get<KnowledgeArticle>(`/knowledge/${id}`)
}

export const createKnowledgeArticle = (data: CreateKnowledgeArticleRequest) => {
  return request.post<KnowledgeArticle>('/knowledge', data)
}

export const updateKnowledgeArticle = (id: number, data: UpdateKnowledgeArticleRequest) => {
  return request.put(`/knowledge/${id}`, data)
}

export const deleteKnowledgeArticle = (id: number) => {
  return request.delete(`/knowledge/${id}`)
}

export const searchKnowledgeArticles = (keyword: string) => {
  return request.get<KnowledgeArticle[]>('/knowledge/search', { params: { keyword } })
}

export const convertFromRepair = (data: ConvertFromRepairRequest) => {
  return request.post<KnowledgeArticle>('/knowledge/convert-repair', data)
}
