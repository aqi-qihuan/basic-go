import request from '@/utils/request'
import type { Comment } from '@/types/comment'

const unwrap = (response: any) => response.data.data

// ========== 评论系统 ==========
// 注意: 后端 CommentService 可能独立部署，路径需要根据实际 BFF 代理确认
// 暂定路径前缀 /comments

/** 获取评论列表 - 支持分页，按 ID 降序 */
export const getComments = async (params: {
  biz: string
  bizId: number
  minId?: number
  limit?: number
}): Promise<{ comments: Comment[] }> => {
  const response = await request.get('/comments/list', {
    params: { biz: params.biz, bizid: params.bizId, min_id: params.minId, limit: params.limit || 20 }
  })
  return unwrap(response)
}

/** 创建评论 */
export const createComment = async (data: {
  biz: string
  bizId: number
  content: string
  rootId?: number
  parentId?: number
}): Promise<void> => {
  const response = await request.post('/comments/create', {
    comment: {
      biz: data.biz,
      bizid: data.bizId,
      content: data.content,
      root_comment: data.rootId ? { id: data.rootId } : undefined,
      parent_comment: data.parentId ? { id: data.parentId } : undefined,
    }
  })
  return unwrap(response)
}

/** 删除评论 */
export const deleteComment = async (id: number): Promise<void> => {
  const response = await request.post('/comments/delete', { id })
  return unwrap(response)
}

/** 加载更多回复（嵌套评论懒加载） */
export const getMoreReplies = async (rid: number, maxId?: number, limit = 10) => {
  const response = await request.get('/comments/replies', {
    params: { rid, max_id: maxId, limit }
  })
  return unwrap(response)
}
