import axios, { AxiosInstance, AxiosResponse } from 'axios'
import { message } from 'antd'

// 创建 axios 实例
const request: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器 - 自动携带 JWT Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers!.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// 响应拦截器
// 后端返回格式: { code: number, msg: string, data: any }
// code === 0 表示成功
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { data } = response

    // 后端 ginx.Result: code=0 成功, code>0 业务错误
    if (data.code === 0 || data.code === undefined) {
      return response
    }

    // 业务错误
    message.error(data.msg || '请求失败')
    return Promise.reject(new Error(data.msg || '请求失败'))
  },
  (error) => {
    if (error.response?.status === 401) {
      message.error('请先登录')
      localStorage.removeItem('token')
      window.location.href = '/login'
    } else if (error.response?.status === 404) {
      message.error('接口不存在')
    } else {
      message.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export default request
