import React, { useEffect, useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { Avatar, Button, message, Spin, Tooltip } from 'antd'
import {
  ArrowLeftOutlined, LikeOutlined, LikeFilled,
  StarOutlined, StarFilled, ShareAltOutlined,
  UserOutlined, ClockCircleOutlined, EyeOutlined,
  MessageOutlined
} from '@ant-design/icons'
import { getPublishedArticle, likeArticle, collectArticle } from '@/services/articleService'
import CommentList from '@/components/CommentList'
import CommentForm from '@/components/CommentForm'
import { TagPill } from '@/components/common'
import { useUserStore } from '@/store/userStore'
import type { Article } from '@/types/article'
import dayjs from 'dayjs'

/** HOK 营地风格文章详情页 */
const ArticleDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const [article, setArticle] = useState<Article | null>(null)
  const [loading, setLoading] = useState(true)
  const [liked, setLiked] = useState(false)
  const [favorited, setFavorited] = useState(false)
  const [likeCount, setLikeCount] = useState(0)
  const [commentRefreshKey, setCommentRefreshKey] = useState(0)
  const navigate = useNavigate()
  const { isLogin } = useUserStore()

  useEffect(() => {
    if (id) fetchArticle(parseInt(id))
  }, [id])

  const fetchArticle = async (articleId: number) => {
    setLoading(true)
    try {
      const data = await getPublishedArticle(articleId)
      setArticle(data)
      setLikeCount(data.likeCnt || 0)
      setLiked(data.liked || false)
      setFavorited(data.collected || false)
    } catch {
      message.error('获取文章详情失败')
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  const checkLogin = (action: () => void) => {
    if (!isLogin) { message.warning('请先登录'); navigate('/login'); return }
    action()
  }

  const handleLike = () => checkLogin(async () => {
    if (!article) return
    try {
      if (liked) {
        await likeArticle(article.id, false)
        setLiked(false); setLikeCount(c => Math.max(0, c - 1))
      } else {
        await likeArticle(article.id, true)
        setLiked(true); setLikeCount(c => c + 1)
      }
    } catch { message.error('操作失败') }
  })

  const handleFavorite = () => checkLogin(async () => {
    if (!article) return
    try {
      if (favorited) {
        await collectArticle(article.id, 0) // cid=0 取消收藏
        setFavorited(false)
      } else {
        await collectArticle(article.id, 0) // cid=0 默认收藏夹
        setFavorited(true)
      }
    } catch { message.error('操作失败') }
  })

  const handleShare = () => {
    navigator.clipboard.writeText(window.location.href)
    message.success('链接已复制到剪贴板')
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center" style={{ minHeight: '60vh', background: 'var(--bg-deep)' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (!article) return null

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)' }}>
      {/* 顶部操作栏 */}
      <div style={{
        position: 'sticky', top: 64, zIndex: 10,
        background: 'rgba(11, 13, 23, 0.95)',
        backdropFilter: 'blur(20px)',
        borderBottom: '1px solid rgba(240, 192, 96, 0.06)',
        padding: '12px 0',
      }}>
        <div className="max-w-4xl mx-auto px-4 flex items-center justify-between">
          <Button
            type="text"
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate(-1)}
            style={{ color: '#9C9688' }}
          >
            返回
          </Button>
          <div className="flex items-center gap-2">
            <Tooltip title={favorited ? '取消收藏' : '收藏'}>
              <Button
                type="text"
                icon={favorited ? <StarFilled style={{ color: '#F0C060' }} /> : <StarOutlined />}
                onClick={handleFavorite}
                style={{ color: favorited ? '#F0C060' : '#9C9688' }}
              />
            </Tooltip>
            <Tooltip title="分享">
              <Button
                type="text"
                icon={<ShareAltOutlined />}
                onClick={handleShare}
                style={{ color: '#9C9688' }}
              />
            </Tooltip>
          </div>
        </div>
      </div>

      <article className="max-w-4xl mx-auto px-4 py-8">
        {/* 文章头部 */}
        <div className="mb-8">
          {/* 标签 */}
          {article.tags && article.tags.length > 0 && (
            <div className="flex flex-wrap gap-2 mb-4">
              {article.tags.map((tag) => (
                <TagPill key={tag} tech>{tag}</TagPill>
              ))}
            </div>
          )}

          {/* 标题 */}
          <h1 style={{
            fontSize: 36, fontWeight: 800, color: '#F5F0E8',
            lineHeight: 1.3, marginBottom: 20,
            fontFamily: "'Inter', 'Noto Sans SC', sans-serif",
          }}>
            {article.title}
          </h1>

          {/* 作者信息 */}
          <div className="flex items-center gap-4 flex-wrap" style={{ marginBottom: 20 }}>
            <div className="flex items-center gap-2">
              <Avatar
                size={36}
                src={(article as any).author?.avatar}
                icon={<UserOutlined />}
                style={{
                  background: 'rgba(240, 192, 96, 0.15)',
                  border: '2px solid rgba(240, 192, 96, 0.3)',
                }}
              />
              <div>
                <Link
                  to={`/profile/${(article as any).author?.id}`}
                  style={{ color: '#E8E0D0', fontWeight: 600, fontSize: 15, textDecoration: 'none' }}
                >
                  {(article as any).author?.name || '匿名'}
                </Link>
              </div>
            </div>
            <span style={{ color: '#6B6558', fontSize: 13 }}>
              <ClockCircleOutlined style={{ marginRight: 4 }} />
              {dayjs(article.ctime).format('YYYY-MM-DD HH:mm')}
            </span>
            {article.viewCount !== undefined && (
              <span style={{ color: '#6B6558', fontSize: 13 }}>
                <EyeOutlined style={{ marginRight: 4 }} />
                {article.viewCount} 次阅读
              </span>
            )}
          </div>

          {/* 金色分隔线 */}
          <div className="gold-divider" />
        </div>

        {/* 文章正文 */}
        <div
          className="article-content mb-12"
          dangerouslySetInnerHTML={{ __html: article.content }}
        />

        {/* 互动栏 */}
        <div
          className="glass-card flex items-center justify-center gap-6 mb-12"
          style={{ padding: '20px 24px', cursor: 'default' }}
        >
          <Button
            type="text"
            size="large"
            icon={liked
              ? <LikeFilled className="animate-like" style={{ color: '#F0C060', fontSize: 22 }} />
              : <LikeOutlined style={{ fontSize: 22 }} />
            }
            onClick={handleLike}
            style={{
              color: liked ? '#F0C060' : '#9C9688',
              fontWeight: liked ? 700 : 400,
              height: 48, padding: '0 24px',
            }}
          >
            {likeCount > 0 && likeCount}
          </Button>

          <Button
            type="text"
            size="large"
            icon={favorited
              ? <StarFilled style={{ color: '#F0C060', fontSize: 22 }} />
              : <StarOutlined style={{ fontSize: 22 }} />
            }
            onClick={handleFavorite}
            style={{
              color: favorited ? '#F0C060' : '#9C9688',
              fontWeight: favorited ? 700 : 400,
              height: 48, padding: '0 24px',
            }}
          />

          <Button
            type="text"
            size="large"
            icon={<ShareAltOutlined style={{ fontSize: 22 }} />}
            onClick={handleShare}
            style={{ color: '#9C9688', height: 48, padding: '0 24px' }}
          />
        </div>

        {/* 评论区 */}
        <div style={{ maxWidth: 768, margin: '0 auto' }}>
          <div className="flex items-center gap-2 mb-6">
            <MessageOutlined style={{ color: '#F0C060', fontSize: 20 }} />
            <h2 style={{ color: '#E8E0D0', fontSize: 22, fontWeight: 700, margin: 0 }}>
              评论
              {(article as any).commentCount !== undefined && (
                <span style={{ color: '#6B6558', fontSize: 14, fontWeight: 400, marginLeft: 8 }}>
                  ({(article as any).commentCount})
                </span>
              )}
            </h2>
            <div className="gold-line flex-1 ml-3" />
          </div>

          {/* 评论输入 */}
          {isLogin ? (
            <div className="glass-card mb-8" style={{ padding: 20, cursor: 'default' }}>
              <CommentForm articleId={article.id} onSuccess={() => setCommentRefreshKey(k => k + 1)} />
            </div>
          ) : (
            <div className="glass-card mb-8 text-center" style={{ padding: 24, cursor: 'default' }}>
              <p style={{ color: '#6B6558', margin: 0 }}>
                请<Link to="/login" style={{ color: '#F0C060' }}>登录</Link>后发表评论
              </p>
            </div>
          )}

          {/* 评论列表 */}
          <CommentList articleId={article.id} refreshKey={commentRefreshKey} />
        </div>
      </article>
    </div>
  )
}

export default ArticleDetailPage
