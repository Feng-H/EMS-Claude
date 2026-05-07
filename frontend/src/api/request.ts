import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 30000,
})

// Request interceptor
request.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
request.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    const authStore = useAuthStore()

    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('未授权，请重新登录')
          authStore.logout()
          break
        case 403:
          ElMessage.error('没有权限访问该资源')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500: {
          const errorData = error.response.data
          const errorMsg = (typeof errorData?.error === 'object' ? errorData.error.message : errorData?.error) || 
                          errorData?.message || '服务器错误，请稍后重试'
          ElMessage.error(errorMsg)
          break
        }
        default: {
          const errorData = error.response.data
          const errorMsg = (typeof errorData?.error === 'object' ? errorData.error.message : errorData?.error) || 
                          errorData?.message || '请求失败'
          ElMessage.error(errorMsg)
        }
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }

    return Promise.reject(error)
  }
)

export default request
