import React, { useEffect, useState } from 'react'
import { Button, Spin, message, Pagination } from 'antd'
import { StarFilled, DeleteOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { getCollections, cancelCollection } from '@/services/articleService'
import { useUserStore } from '@/store/userStore'
import { GlassCard, EmptyState } from '@/components/common'
import dayjs from 'dayjs'

interface CollectionItem {
  bizId: number
  title?: string
  abstract?: string
  ctime: number
}

const CollectionsPage: React.FC = () => {
  const [items, setItems] = useState<CollectionItem[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const pageSize = 10
  const { isLogin } = useUserStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (!isLogin) {
      navigate('/login')
      return
    }
    fetchCollections()
  }, [page])

  const fetchCollections = async () => {
    setLoading(true)
    try {
      const res = await getCollections((page - 1) * pageSize, pageSize)
      setItems(res?.items || [])
      setTotal(res?.total || 0)
    } catch {
      message.error('获取收藏列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = async (bizId: number) => {
    try {
      await cancelCollection(bizId)
      message.success('已取消收藏')
      fetchCollections()
    } catch {
      message.error('取消收藏失败')
    }
  }

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)' }} className="px-3 sm:px-4 py-4 sm:py-6">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center gap-2 sm:gap-3 mb-4 sm:mb-6">
          <StarFilled style={{ color: '#F0C060' }} className="text-xl sm:text-2xl" />
          <h1 style={{ fontWeight: 800, color: '#F0C060', margin: 0 }}
            className="text-xl sm:text-2xl">
            我的收藏
          </h1>
          <span style={{ color: '#6B6558' }} className="text-xs sm:text-sm">共 {total} 篇</span>
        </div>

        {loading ? (
          <div className="flex justify-center py-16 sm:py-20">
            <Spin size="large" tip="加载中..." />
          </div>
        ) : items.length === 0 ? (
          <GlassCard style={{ cursor: 'default' }}>
            <div className="p-6 sm:p-10">
              <EmptyState
                icon={<StarFilled />}
                title="暂无收藏"
                description="去发现好文章并收藏吧"
                action={
                  <Button
                    type="primary"
                    onClick={() => navigate('/')}
                    style={{
                      background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                      border: 'none', borderRadius: 10, fontWeight: 700, color: '#0B0D17',
                      height: 44,
                    }}
                  >
                    去首页看看
                  </Button>
                }
              />
            </div>
          </GlassCard>
        ) : (
          <>
            <div className="flex flex-col gap-2 sm:gap-3">
              {items.map((item) => (
                <GlassCard
                  key={item.bizId}
                  className="hover-lift cursor-pointer"
                  onClick={() => navigate(`/article/${item.bizId}`)}
                >
                  <div className="p-3 sm:p-4">
                    <div className="flex items-start justify-between gap-3 sm:gap-4">
                      <div className="flex-1 min-w-0">
                        <h3 style={{
                          fontWeight: 700, color: '#F5F0E8',
                          margin: '0 0 6px 0', lineHeight: 1.4,
                        }} className="text-sm sm:text-base">
                          {item.title || `文章 #${item.bizId}`}
                        </h3>
                        {item.abstract && (
                          <p style={{
                            color: '#9C9688', lineHeight: 1.6,
                            margin: '0 0 8px 0',
                            display: '-webkit-box',
                            WebkitLineClamp: 2,
                            WebkitBoxOrient: 'vertical',
                            overflow: 'hidden',
                          }} className="text-xs sm:text-sm">
                            {item.abstract}
                          </p>
                        )}
                        <div className="flex items-center gap-3 sm:gap-4" style={{ color: '#6B6558' }}>
                          <span className="text-xs">收藏于 {dayjs(item.ctime).format('YYYY-MM-DD')}</span>
                        </div>
                      </div>
                      <div className="flex items-center gap-1 sm:gap-2 flex-shrink-0" onClick={(e) => e.stopPropagation()}>
                        <Button
                          type="text"
                          icon={<EyeOutlined />}
                          onClick={() => navigate(`/article/${item.bizId}`)}
                          style={{ color: '#9C9688', height: 36, width: 36 }}
                        />
                        <Button
                          type="text"
                          danger
                          icon={<DeleteOutlined />}
                          onClick={() => handleCancel(item.bizId)}
                          style={{ height: 36, width: 36 }}
                        />
                      </div>
                    </div>
                  </div>
                </GlassCard>
              ))}
            </div>

            {total > pageSize && (
              <div className="flex justify-center mt-4 sm:mt-6">
                <Pagination
                  current={page}
                  total={total}
                  pageSize={pageSize}
                  onChange={setPage}
                  showSizeChanger={false}
                  size="small"
                />
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}

export default CollectionsPage
