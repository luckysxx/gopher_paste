import request from '@/utils/request'

export interface User {
  id: number
  username: string
  email: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user_id: number
  username: string
  email: string
}

export interface RegisterResponse {
  user_id: number
  username: string
  email: string
}

// 用户登录
// 对应后端: POST /api/v1/users/login
export const login = (data: LoginRequest) => {
  return request.post<any, LoginResponse>('/api/v1/users/login', data)
}

// 用户注册
// 对应后端: POST /api/v1/users/register
export const register = (data: RegisterRequest) => {
  return request.post<any, RegisterResponse>('/api/v1/users/register', data)
}
