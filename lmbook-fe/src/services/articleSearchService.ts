import request from '@/utils/request'

// 搜索文章
export const searchArticles = async (params: {
  keyword: string
  offset?: number
  limit?: number
}): Promise<{ articles: any[]; total: number }> => {
  const response = await request.get('/article/search', { params })
  return response.data.data
}

// 获取热门文章排行榜
export const getHotArticles = async (params: {
  period?: 'day' | 'week' | 'month'
  offset?: number
  limit?: number
}): Promise<{ articles: any[] }> => {
  const response = await request.get('/article/hot', { params })
  return response.data.data
}

// 获取标签列表
export const getTags = async (): Promise<{ tags: { name: string; count: number }[] }> => {
  const response = await request.get('/article/tags')
  return response.data.data
}

// 根据标签获取文章
export const getArticlesByTag = async (tag: string, offset = 0, limit = 20): Promise<{ articles: any[]; total: number }> => {
  const response = await request.get('/article/list_by_tag', { params: { tag, offset, limit } })
  return response.data.data
}
