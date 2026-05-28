import React, { useEffect, useState } from 'react'
import { Avatar, Typography, Spin, Empty, message } from 'antd'
import { UserOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { getComments } from '@/services/commentService'
import type { Comment } from '@/types/comment'
import dayjs from 'dayjs'

const { Text } = Typography

interface CommentListProps {
  articleId: number
  refreshKey?: number
}

const CommentList: React.FC<CommentListProps> = ({ articleId, refreshKey }) => {
  const [comments, setComments] = useState<Comment[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)

  useEffect(() => {
    fetchComments()
  }, [articleId, refreshKey])

  const fetchComments = async () => {
    setLoading(true)
    try {
      const res = await getComments(articleId)
      setComments(res.comments || [])
      setTotal(res.total || 0)
    } catch {
      message.error('获取评论失败')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="loading-spinner">
        <Spin size="large" tip="加载评论中..." />
      </div>
    )
  }

  if (comments.length === 0) {
    return (
      <Empty description="暂无评论，快来抢沙发吧" className="my-8" />
    )
  }

  return (
    <div className="comment-list">
      <div className="mb-4 text-gray-500">
        共 {total} 条评论
      </div>
      <div className="space-y-6">
        {comments.map((comment) => (
          <div key={comment.id} className="flex gap-3 p-4 bg-white rounded-lg shadow-sm">
            <Avatar
              size={40}
              src={comment.author.avatar}
              icon={!comment.author.avatar ? <UserOutlined /> : undefined}
            />
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <Text strong className="text-sm">
                  {comment.author.nickname}
                </Text>
                {comment.replyTo && (
                  <Text type="secondary" className="text-xs">
                    回复 @{comment.replyTo.nickname}
                  </Text>
                )}
                <span className="text-xs text-gray-400">
                  <ClockCircleOutlined className="mr-1" />
                  {dayjs(comment.ctime).format('MM-DD HH:mm')}
                </span>
              </div>
              <div className="text-gray-700 text-sm leading-relaxed whitespace-pre-wrap">
                {comment.content}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export default CommentList
