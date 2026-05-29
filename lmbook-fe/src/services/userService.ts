import request from '@/utils/request'
import type { User } from '@/types/user'

// 解包后端响应: { code: 0, msg: "OK", data: ... }
const unwrap = (response: any) => response.data.data

// ========== 用户认证 ==========

/** 邮箱注册 - POST /users/signup */
export const register = async (data: {
  email: string
  password: string
  confirmPassword: string
}): Promise<void> => {
  const response = await request.post('/users/signup', data)
  return unwrap(response)
}

/** 邮箱登录 - POST /users/login */
export const login = async (email: string, password: string): Promise<void> => {
  const response = await request.post('/users/login', { email, password })
  return unwrap(response)
  // Token 通过响应头 x-jwt-token / x-refresh-token 返回
}

/** 退出登录 - POST /users/logout */
export const logout = async (): Promise<void> => {
  const response = await request.post('/users/logout')
  return unwrap(response)
}

/** 刷新 Token - POST /users/refresh_token */
export const refreshToken = async (): Promise<void> => {
  const response = await request.post('/users/refresh_token')
  return unwrap(response)
}

// ========== 微信登录 ==========

/** 获取微信授权 URL - GET /oauth2/wechat/authurl */
export const getWechatAuthUrl = async (): Promise<string> => {
  const response = await request.get('/oauth2/wechat/authurl')
  return response.data.data
}

// ========== 短信登录 ==========

/** 发送短信验证码 - POST /users/login_sms/code/send */
export const sendSmsCode = async (phone: string): Promise<void> => {
  const response = await request.post('/users/login_sms/code/send', { phone })
  return unwrap(response)
}

/** 手机号登录 - POST /users/login_sms */
export const loginBySms = async (phone: string, code: string): Promise<void> => {
  const response = await request.post('/users/login_sms', { phone, code })
  return unwrap(response)
}

// ========== 用户信息 ==========

/** 获取用户详情 - GET /users/profile */
export const getUserProfile = async (): Promise<User> => {
  const response = await request.get('/users/profile')
  return response.data  // ProfileJWT 直接返回 user 对象
}

/** 编辑资料 - POST /users/edit */
export const editUserProfile = async (data: {
  nickname: string
  birthday?: string
  aboutMe?: string
}): Promise<void> => {
  const response = await request.post('/users/edit', data)
  return unwrap(response)
}
