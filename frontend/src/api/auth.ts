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
  lark_openid?: string
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

export interface LarkConfigReq {
  app_id: string;
  app_secret: string;
  verification_token: string;
  encrypt_key: string;
}

export interface LarkConfigResp {
  app_id: string;
  has_app_secret: boolean;
  verification_token: string;
  has_encrypt_key: boolean;
  webhook_url: string;
}

export const authApi = {
  login: (data: LoginRequest) => request.post<LoginResponse>('/auth/login', data),
  logout: () => request.post('/auth/logout'),
  getUserInfo: () => request.get<UserInfo>('/auth/me'),
  refreshToken: (token: string) => request.post<{ token: string; expire_at: number }>('/auth/refresh', { token }),
  changePassword: (data: ChangePasswordRequest) => request.post('/auth/change-password', data),
  applyAccount: (data: ApplyAccountRequest) => request.post('/auth/apply', data),
  bindLark: (openid: string) => request.post('/auth/bind-lark', { openid }),
  getLarkConfig: () => request.get<LarkConfigResp>('/auth/lark-config'),
  updateLarkConfig: (data: LarkConfigReq) => request.put('/auth/lark-config', data),
}
