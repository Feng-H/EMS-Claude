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
  }
}

export interface UserInfo {
  id: number
  username: string
  name: string
  role: string
  factory_id?: number
}

export const authApi = {
  login: (data: LoginRequest) => request.post<LoginResponse>('/auth/login', data),
  logout: () => request.post('/auth/logout'),
  getUserInfo: () => request.get<UserInfo>('/auth/me'),
  refreshToken: (token: string) => request.post<{ token: string; expire_at: number }>('/auth/refresh', { token }),
}
