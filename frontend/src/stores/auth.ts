import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type LoginRequest, type LoginResponse, type UserInfo } from '@/api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('ems_token') || '')
  const userInfo = ref<UserInfo | null>(JSON.parse(localStorage.getItem('ems_user') || 'null'))

  const isLoggedIn = computed(() => !!token.value)
  const userRole = computed(() => userInfo.value?.role || '')
  const userName = computed(() => userInfo.value?.name || '')
  const userId = computed(() => userInfo.value?.id || 0)
  const mustChangePassword = computed(() => userInfo.value?.must_change_password || false)

  async function login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await authApi.login(credentials)
    token.value = response.token
    userInfo.value = response.user_info

    localStorage.setItem('ems_token', response.token)
    localStorage.setItem('ems_user', JSON.stringify(response.user_info))

    return response
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('ems_token')
    localStorage.removeItem('ems_user')
  }

  async function getUserInfo() {
    const user = await authApi.getUserInfo()
    userInfo.value = user
    localStorage.setItem('ems_user', JSON.stringify(user))
  }

  function hasRole(...roles: string[]): boolean {
    return roles.includes(userRole.value)
  }

  function hasPermission(requiredRole?: string): boolean {
    if (!requiredRole) return true
    const roleHierarchy: Record<string, number> = {
      admin: 5,
      supervisor: 4,
      engineer: 3,
      maintenance: 2,
      operator: 1,
    }
    return (roleHierarchy[userRole.value] || 0) >= (roleHierarchy[requiredRole] || 0)
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userRole,
    userName,
    userId,
    mustChangePassword,
    login,
    logout,
    getUserInfo,
    hasRole,
    hasPermission,
  }
})
