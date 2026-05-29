import request from '@/utils/request'

const unwrap = (response: any) => response.data.data

// ========== 搜索系统 ==========
// 注意: 后端 SearchService 可能独立部署，路径需要根据实际 BFF 代理确认
// 暂定路径前缀 /search

/** 全文搜索（同时返回用户+文章） */
export const search = async (expression: string) => {
  const response = await request.get('/search', { params: { expression } })
  return unwrap(response)
}

/** 搜索文章（封装） */
export const searchArticles = async (params: {
  keyword: string
  offset?: number
  limit?: number
}) => {
  const result = await search(params.keyword)
  return { articles: result.article?.articles || [], total: result.article?.articles?.length || 0 }
}

// ========== 排行榜 ==========
// 注意: 后端 RankingService 可能独立部署

/** 获取热门文章 Top N */
export const getHotArticles = async (params?: {
  limit?: number
}): Promise<{ articles: any[] }> => {
  const response = await request.get('/ranking/top', { params })
  return unwrap(response)
}

// ========== 标签系统 ==========
// 注意: 后端 TagService 可能独立部署

/** 获取用户的标签列表 */
export const getUserTags = async () => {
  const response = await request.get('/tags/list')
  return unwrap(response)
}

/** 获取文章的标签 */
export const getBizTags = async (biz: string, bizId: number) => {
  const response = await request.get('/tags/biz', {
    params: { biz, biz_id: bizId }
  })
  return unwrap(response)
}

/** 创建标签 */
export const createTag = async (name: string) => {
  const response = await request.post('/tags/create', { name })
  return unwrap(response)
}

/** 给文章附加标签（覆盖式） */
export const attachTags = async (biz: string, bizId: number, tids: number[]) => {
  const response = await request.post('/tags/attach', { biz, biz_id: bizId, tids })
  return unwrap(response)
}

// ========== 互动系统 ==========

/** 批量获取互动数据（点赞/收藏状态） */
export const getInteractiveByIds = async (biz: string, ids: number[]) => {
  const response = request.get('/interactive/by_ids', {
    params: { biz, ids: ids.join(',') }
  })
  return unwrap(await response)
}

// ========== Feed 流 ==========

/** 获取用户 Feed 流 */
export const getFeedEvents = async (limit?: number, timestamp?: number) => {
  const response = await request.get('/feed/events', {
    params: { limit, timestamp }
  })
  return unwrap(response)
}
