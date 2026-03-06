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

type RawPaste = Partial<{
  short_link: string
  ShortLink: string
  content: string
  Content: string
  language: string
  Language: string
  created_at: string
  CreatedAt: string
  expires_at: string
  ExpiresAt: string
  valid: boolean
  Valid: boolean
  id: string | number
  ID: string | number
}>

const normalizePaste = (raw: RawPaste): Paste => {
  const shortLink = raw.short_link ?? raw.ShortLink
  const content = raw.content ?? raw.Content
  const language = raw.language ?? raw.Language
  const createdAt = raw.created_at ?? raw.CreatedAt

  if (!shortLink || !content || !language || !createdAt) {
    throw new Error('Paste 响应字段缺失')
  }

  return {
    short_link: shortLink,
    content,
    language,
    created_at: createdAt,
    expires_at: raw.expires_at ?? raw.ExpiresAt,
    valid: raw.valid ?? raw.Valid,
    id: raw.id ?? raw.ID,
  }
}

// 1. 获取帖子详情的接口
// 对应后端: GET /api/v1/pastes/:id
export const getPaste = (id: string) => {
  return request.get<unknown, RawPaste>(`/api/v1/pastes/${id}`).then(normalizePaste)
}

// 2. 创建帖子的接口
// 对应后端: POST /api/v1/pastes
export const createPaste = (data: { content: string; language: string }) => {
  return request.post<unknown, RawPaste>('/api/v1/pastes', data).then(normalizePaste)
}
