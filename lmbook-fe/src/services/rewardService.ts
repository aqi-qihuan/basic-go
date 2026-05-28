import request from '@/utils/request'

const unwrap = (response: any) => response.data.data

// ========== 打赏系统 ==========

/** 查询打赏详情 - POST /reward/detail */
export const getRewardDetail = async (rid: number): Promise<{ status: string }> => {
  const response = await request.post('/reward/detail', { Rid: rid })
  return unwrap(response)
}
