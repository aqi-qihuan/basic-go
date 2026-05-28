// 用户信息
export interface User {
  id: number
  email: string
  nickname: string
  avatar?: string
  bio?: string
  articleCount?: number
  followerCount?: number
  followingCount?: number
  ctime?: string
}

// 登录请求
export interface LoginRequest {
  email: string
  password: string
}

// 注册请求
export interface RegisterRequest {
  email: string
  password: string
  nickname: string
  captcha?: string
}

// 通用 API 响应
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
