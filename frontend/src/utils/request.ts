import axios from 'axios'
import type { AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

// 1. 创建 axios 实例
const service = axios.create({
  // 这里留空，让它自动匹配当前域名
  // 也就是请求会发给 http://localhost:5173/api/v1/...
  // 然后被 Vite 代理捕获
  baseURL: '',
  timeout: 5000, // 请求超时时间：5秒
})

// 2. 请求拦截器 (Request Interceptor)
service.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 后端统一响应格式
interface ApiResponse<T = unknown> {
  code: number
  msg: string
  data: T
}

// 3. 响应拦截器 (Response Interceptor)
service.interceptors.response.use(
  (response) => {
    // 2xx 范围内的状态码都会触发这里
    const res = response.data as ApiResponse

    // 后端统一包装格式：{ code, msg, data }
    // code 为 200 表示成功，直接返回 data
    if (res.code === 200) {
      return res.data as unknown as AxiosResponse
    }

    // code 不为 200，说明业务逻辑错误
    ElMessage.error(res.msg || '请求失败')
    return Promise.reject(new Error(res.msg || '请求失败'))
  },
  (error) => {
    // 超出 2xx 范围的状态码都会触发这里（网络错误、404、500等）
    console.error('请求出错:', error)

    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      if (window.location.pathname !== '/auth') {
        window.location.href = '/auth'
      }
    }

    // 统一弹出错误提示
    const msg = error.response?.data?.msg || error.message || '网络请求失败，请稍后重试'
    ElMessage.error(msg)

    return Promise.reject(error)
  },
)

export default service
