import request from '@/utils/request'

const unwrap = (response: any) => response.data.data

// ========== 关注系统 ==========
// 注意: 后端 FollowService 可能独立部署，路径需要根据实际 BFF 代理确认
// 暂定路径前缀 /follows

/** 关注用户 */
export const followUser = async (followee: number): Promise<void> => {
  const response = await request.post('/follows/follow', { followee })
  return unwrap(response)
}

/** 取消关注 */
export const unfollowUser = async (followee: number): Promise<void> => {
  const response = await request.post('/follows/cancel', { followee })
  return unwrap(response)
}

/** 获取关注列表 */
export const getFolloweeList = async (follower: number, offset = 0, limit = 20) => {
  const response = await request.get('/follows/followee', {
    params: { follower, offset, limit }
  })
  return unwrap(response)
}

/** 获取粉丝列表 */
export const getFollowerList = async (followee: number, offset = 0, limit = 20) => {
  const response = await request.get('/follows/follower', {
    params: { followee, offset, limit }
  })
  return unwrap(response)
}

/** 查询关注关系 */
export const getFollowInfo = async (follower: number, followee: number) => {
  const response = await request.get('/follows/info', {
    params: { follower, followee }
  })
  return unwrap(response)
}

/** 获取关注/粉丝数统计 */
export const getFollowStats = async (followee: number) => {
  const response = await request.get('/follows/stats', {
    params: { followee }
  })
  return unwrap(response)
}
