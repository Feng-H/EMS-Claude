import request from './request'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expire_at: number
  user_info: {
    id: number
    username: string
    name: string
    role: string
    factory_id?: number
    approval_status: string
    must_change_password: boolean
  }
  must_change_password: boolean
}

export interface UserInfo {
  id: number
  username: string
  name: string
  role: string
  factory_id?: number
  approval_status: string
  must_change_password: boolean
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface ApplyAccountRequest {
  username: string
  password: string
  name: string
  role: string
  factory_id?: number
  phone?: string
}

export const authApi = {
  login: (data: LoginRequest) => request.post<LoginResponse>('/auth/login', data),
  logout: () => request.post('/auth/logout'),
  getUserInfo: () => request.get<UserInfo>('/auth/me'),
  refreshToken: (token: string) => request.post<{ token: string; expire_at: number }>('/auth/refresh', { token }),
  changePassword: (data: ChangePasswordRequest) => request.post('/auth/change-password', data),
  applyAccount: (data: ApplyAccountRequest) => request.post('/auth/apply', data),
}
