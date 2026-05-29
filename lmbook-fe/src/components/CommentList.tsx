import React, { useEffect, useState } from 'react'
import { Avatar, Button, Spin, message } from 'antd'
import { UserOutlined, LikeOutlined, MessageOutlined, DownOutlined } from '@ant-design/icons'
import { getComments, getMoreReplies } from '@/services/commentService'
import CommentForm from '@/components/CommentForm'
import { EmptyState } from '@/components/common'
import type { Comment } from '@/types/comment'
import dayjs from 'dayjs'

interface CommentListProps {
  articleId: number
  refreshKey?: number
}

/** HOK 金色风格评论列表 - 支持嵌套回复 + 懒加载 */
const CommentList: React.FC<CommentListProps> = ({ articleId, refreshKey }) => {
  const [comments, setComments] = useState<Comment[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [replyTo, setReplyTo] = useState<number | null>(null)
  const [replyRoot, setReplyRoot] = useState<number | null>(null)
  const [repliesMap, setRepliesMap] = useState<Record<number, Comment[]>>({})
  const [moreRepliesMap, setMoreRepliesMap] = useState<Record<number, { maxId: number; hasMore: boolean }>>({})
  const [loadingReplies, setLoadingReplies] = useState<Record<number, boolean>>({})

  useEffect(() => { fetchComments() }, [articleId, refreshKey])

  const fetchComments = async () => {
    setLoading(true)
    try {
      const res = await getComments({ biz: 'article', bizId: articleId })
      setComments(res?.comments || [])
      setTotal(res?.comments?.length || 0)
    } catch {
      message.error('获取评论失败')
    } finally {
      setLoading(false)
    }
  }

  // 加载某条评论的回复
  const loadReplies = async (commentId: number) => {
    setLoadingReplies(prev => ({ ...prev, [commentId]: true }))
    try {
      const res = await getMoreReplies(commentId)
      const replies = res?.replies || []
      setRepliesMap(prev => ({ ...prev, [commentId]: replies }))
      if (replies.length > 0) {
        setMoreRepliesMap(prev => ({
          ...prev,
          [commentId]: { maxId: replies[replies.length - 1].id, hasMore: replies.length >= 10 }
        }))
      }
    } catch {} finally {
      setLoadingReplies(prev => ({ ...prev, [commentId]: false }))
    }
  }

  // 加载更多回复
  const loadMoreReplies = async (commentId: number) => {
    const info = moreRepliesMap[commentId]
    if (!info) return
    setLoadingReplies(prev => ({ ...prev, [commentId]: true }))
    try {
      const res = await getMoreReplies(commentId, info.maxId)
      const newReplies = res?.replies || []
      setRepliesMap(prev => ({
        ...prev,
        [commentId]: [...(prev[commentId] || []), ...newReplies]
      }))
      if (newReplies.length > 0) {
        setMoreRepliesMap(prev => ({
          ...prev,
          [commentId]: { maxId: newReplies[newReplies.length - 1].id, hasMore: newReplies.length >= 10 }
        }))
      } else {
        setMoreRepliesMap(prev => ({ ...prev, [commentId]: { ...info, hasMore: false } }))
      }
    } catch {} finally {
      setLoadingReplies(prev => ({ ...prev, [commentId]: false }))
    }
  }

  const handleReplySuccess = () => {
    setReplyTo(null)
    setReplyRoot(null)
    fetchComments()
    if (replyRoot) loadReplies(replyRoot)
  }

  const commentStyle: React.CSSProperties = {
    background: 'rgba(19, 21, 32, 0.6)',
    borderRadius: 10,
    padding: '16px 20px',
    border: '1px solid rgba(240, 192, 96, 0.2)',
  }

  const renderComment = (comment: Comment, isReply = false) => (
    <div key={comment.id} style={{ ...commentStyle, marginLeft: isReply ? 40 : 0 }}>
      <div className="flex gap-3">
        <Avatar
          size={isReply ? 28 : 36}
          icon={<UserOutlined />}
          style={{ background: 'rgba(240,192,96,0.15)', border: '1px solid rgba(240,192,96,0.3)', flexShrink: 0 }}
        />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span style={{ color: '#F5F0E8', fontWeight: 600, fontSize: isReply ? 13 : 14 }}>
              {comment.author?.nickname || '匿名用户'}
            </span>
            {comment.replyTo && (
              <span style={{ color: '#6B6558', fontSize: 12 }}>
                回复 #{comment.replyTo.id}
              </span>
            )}
            <span style={{ color: '#6B6558', fontSize: 12 }}>
              {dayjs(comment.ctime).format('MM-DD HH:mm')}
            </span>
          </div>

          <p style={{ color: '#F5F0E8', fontSize: isReply ? 13 : 14, lineHeight: 1.6, margin: '0 0 8px 0', whiteSpace: 'pre-wrap' }}>
            {comment.content}
          </p>

          <div className="flex items-center gap-4">
            <Button
              type="text"
              size="small"
              icon={<LikeOutlined />}
              style={{ color: '#6B6558', fontSize: 12, padding: '0 4px', height: 24 }}
            />
            {!isReply && (
              <Button
                type="text"
                size="small"
                icon={<MessageOutlined />}
                onClick={() => { setReplyTo(comment.id); setReplyRoot(comment.id) }}
                style={{ color: '#6B6558', fontSize: 12, padding: '0 4px', height: 24 }}
              >
                回复
              </Button>
            )}
          </div>

          {/* 回复输入框 */}
          {replyTo === comment.id && (
            <div className="mt-3">
              <CommentForm
                articleId={articleId}
                parentId={comment.id}
                onSuccess={handleReplySuccess}
                placeholder={`回复 #${comment.id}...`}
              />
            </div>
          )}

          {/* 嵌套回复 */}
          {!isReply && (
            <div className="mt-3">
              {repliesMap[comment.id] ? (
                <div className="flex flex-col gap-2">
                  {repliesMap[comment.id].map(reply => renderComment(reply, true))}
                  {moreRepliesMap[comment.id]?.hasMore && (
                    <Button
                      type="text"
                      size="small"
                      loading={loadingReplies[comment.id]}
                      onClick={() => loadMoreReplies(comment.id)}
                      icon={<DownOutlined />}
                      style={{ color: '#F0C060', fontSize: 12, alignSelf: 'flex-start' }}
                    >
                      加载更多回复
                    </Button>
                  )}
                </div>
              ) : (
                <Button
                  type="text"
                  size="small"
                  loading={loadingReplies[comment.id]}
                  onClick={() => loadReplies(comment.id)}
                  icon={<DownOutlined />}
                  style={{ color: '#6B6558', fontSize: 12, padding: '0 4px', height: 24 }}
                >
                  查看回复
                </Button>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )

  if (loading) {
    return <div className="loading-spinner"><Spin size="large" tip="加载评论中..." /></div>
  }

  if (comments.length === 0) {
    return <EmptyState title="暂无评论" description="快来抢沙发吧" />
  }

  return (
    <div>
      <div style={{ color: '#6B6558', fontSize: 13, marginBottom: 16 }}>
        共 {total} 条评论
      </div>
      <div className="flex flex-col gap-3">
        {comments.map(comment => renderComment(comment))}
      </div>
    </div>
  )
}

export default CommentList
