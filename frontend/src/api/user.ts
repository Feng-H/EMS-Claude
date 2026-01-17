import request from './request'

export interface User {
  id: number
  username: string
  name: string
  role: string
  phone: string
  is_active: boolean
  approval_status: string
  must_change_password: boolean
  factory_id?: number
  created_at: string
}

export interface CreateUserRequest {
  username: string
  password: string
  name: string
  role: string
  factory_id?: number
  phone?: string
}

export interface UpdateUserRequest {
  name?: string
  role?: string
  factory_id?: number
  phone?: string
  is_active?: boolean
}

export interface ApproveUserRequest {
  approve: boolean
  reason?: string
}

export const userApi = {
  // 获取用户列表
  getUsers: () => request.get<User[]>('/users'),

  // 创建用户
  createUser: (data: CreateUserRequest) => request.post<{ id: number; message: string }>('/users', data),

  // 更新用户
  updateUser: (id: number, data: UpdateUserRequest) => request.put<{ message: string }>(`/users/${id}`, data),

  // 获取待审核的申请
  getPendingApplications: () => request.get<User[]>('/users/applications'),

  // 审核用户申请
  approveApplication: (id: number, data: ApproveUserRequest) => request.put<{ message: string }>(`/users/${id}/approve`, data),
}
