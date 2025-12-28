import request from '@/utils/request'

export interface Paste {
  short_link: string
  content: string
  language: string
  created_at: string
  expires_at?: string
  valid?: boolean
  id?: string | number
}

// 1. 获取帖子详情的接口
// 对应后端: GET /api/v1/pastes/:id
export const getPaste = (id: string) => {
  return request.get<any, Paste>(`/api/v1/pastes/${id}`)
}

// 2. 创建帖子的接口
// 对应后端: POST /api/v1/pastes
export const createPaste = (data: { content: string; language: string }) => {
  return request.post<any, Paste>('/api/v1/pastes', data)
}
