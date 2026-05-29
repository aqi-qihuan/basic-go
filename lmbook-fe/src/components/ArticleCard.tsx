import React from 'react'
import { Avatar } from 'antd'
import { EyeOutlined, LikeOutlined, UserOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { TagPill } from '@/components/common'
import type { Article } from '@/types/article'
import dayjs from 'dayjs'

interface ArticleCardProps {
  article: Article
  featured?: boolean
}

/** HOK 金色风格文章卡片 */
const ArticleCard: React.FC<ArticleCardProps> = ({ article, featured = false }) => {
  const navigate = useNavigate()

  const cardStyle: React.CSSProperties = featured ? {
    background: 'linear-gradient(135deg, rgba(240,192,96,0.1) 0%, rgba(11,13,23,0.9) 100%)',
    backdropFilter: 'blur(24px) saturate(180%)',
    border: '1px solid rgba(240, 192, 96, 0.4)',
    borderRadius: 16,
    overflow: 'hidden',
    cursor: 'pointer',
    transition: 'all 250ms cubic-bezier(0.4, 0, 0.2, 1)',
    boxShadow: '0 0 20px rgba(240, 192, 96, 0.2)',
  } : {
    background: 'rgba(19, 21, 32, 0.85)',
    backdropFilter: 'blur(24px) saturate(180%)',
    border: '1px solid rgba(240, 192, 96, 0.2)',
    borderRadius: 16,
    overflow: 'hidden',
    cursor: 'pointer',
    transition: 'all 250ms cubic-bezier(0.4, 0, 0.2, 1)',
    boxShadow: '0 4px 24px rgba(0, 0, 0, 0.4)',
  }

  const handleHover = (e: React.MouseEvent<HTMLDivElement>, enter: boolean) => {
    const el = e.currentTarget
    if (enter) {
      el.style.transform = 'translateY(-4px)'
      el.style.borderColor = 'rgba(240, 192, 96, 0.6)'
      el.style.boxShadow = '0 12px 40px rgba(0,0,0,0.5), 0 0 20px rgba(240,192,96,0.3)'
    } else {
      el.style.transform = 'translateY(0)'
      el.style.borderColor = featured
        ? 'rgba(240, 192, 96, 0.4)'
        : 'rgba(240, 192, 96, 0.2)'
      el.style.boxShadow = featured
        ? '0 0 20px rgba(240, 192, 96, 0.2)'
        : '0 4px 24px rgba(0, 0, 0, 0.4)'
    }
  }

  return (
    <div
      style={cardStyle}
      onClick={() => navigate(`/article/${article.id}`)}
      onMouseEnter={(e) => handleHover(e, true)}
      onMouseLeave={(e) => handleHover(e, false)}
    >
      {/* 封面图 */}
      {article.coverImage && (
        <div style={{ position: 'relative', overflow: 'hidden' }}>
          <img
            src={article.coverImage}
            alt={article.title}
            style={{
              width: '100%',
              height: featured ? 220 : 180,
              objectFit: 'cover',
              display: 'block',
            }}
          />
          {/* 金色渐变遮罩 */}
          <div style={{
            position: 'absolute', bottom: 0, left: 0, right: 0, height: 80,
            background: 'linear-gradient(transparent, rgba(11,13,23,0.9))',
          }} />
        </div>
      )}

      {/* 内容区 */}
      <div style={{ padding: featured ? '20px 24px' : '16px 20px' }}>
        {/* 标签 */}
        {article.tags && article.tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mb-3">
            {article.tags.slice(0, 3).map((tag) => (
              <TagPill key={tag} tech>{tag}</TagPill>
            ))}
          </div>
        )}

        {/* 标题 */}
        <h3
          style={{
            fontSize: featured ? 22 : 16,
            fontWeight: 700,
            color: '#F5F0E8',
            marginBottom: 8,
            lineHeight: 1.4,
            display: '-webkit-box',
            WebkitLineClamp: featured ? 2 : 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
          }}
        >
          <Link
            to={`/article/${article.id}`}
            onClick={(e) => e.stopPropagation()}
            style={{
              color: 'inherit',
              textDecoration: 'none',
              transition: 'color 200ms ease',
            }}
            onMouseEnter={(e) => { (e.target as HTMLElement).style.color = '#F0C060' }}
            onMouseLeave={(e) => { (e.target as HTMLElement).style.color = '#F5F0E8' }}
          >
            {article.title}
          </Link>
        </h3>

        {/* 摘要 */}
        {article.abstract && (
          <p style={{
            fontSize: 14, color: '#9C9688', lineHeight: 1.6,
            marginBottom: 16,
            display: '-webkit-box',
            WebkitLineClamp: featured ? 3 : 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
          }}>
            {article.abstract}
          </p>
        )}

        {/* 底部信息 */}
        <div className="flex items-center justify-between" style={{ paddingTop: 12, borderTop: '1px solid rgba(240, 192, 96, 0.2)' }}>
          <div className="flex items-center gap-2">
            <Avatar
              size={22}
              src={(article as any).author?.avatar}
              icon={<UserOutlined />}
              style={{ background: 'rgba(240, 192, 96, 0.15)', border: '1px solid rgba(240, 192, 96, 0.3)' }}
            />
            <span style={{ fontSize: 13, color: '#9C9688' }}>
              {(article as any).author?.name || '匿名'}
            </span>
            <span style={{ fontSize: 12, color: '#6B6558' }}>
              &middot; {dayjs(article.ctime).format('MM-DD')}
            </span>
          </div>

          <div className="flex items-center gap-4" style={{ fontSize: 12, color: '#6B6558' }}>
            <span className="flex items-center gap-1">
              <EyeOutlined /> {article.viewCount || 0}
            </span>
            <span className="flex items-center gap-1">
              <LikeOutlined /> {article.likeCount || 0}
            </span>
          </div>
        </div>
      </div>
    </div>
  )
}

export default ArticleCard
