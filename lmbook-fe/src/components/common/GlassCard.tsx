import React from 'react'

interface GlassCardProps {
  children: React.ReactNode
  className?: string
  hoverable?: boolean
  gold?: boolean
  onClick?: () => void
  style?: React.CSSProperties
}

/** HOK 风格暗色毛玻璃卡片 */
const GlassCard: React.FC<GlassCardProps> = ({
  children, className = '', hoverable = true, gold = false, onClick, style
}) => {
  const baseStyle: React.CSSProperties = {
    background: gold
      ? 'linear-gradient(135deg, rgba(240,192,96,0.05) 0%, rgba(19,21,32,0.9) 100%)'
      : 'rgba(19, 21, 32, 0.85)',
    backdropFilter: 'blur(24px) saturate(180%)',
    WebkitBackdropFilter: 'blur(24px) saturate(180%)',
    border: gold
      ? '1px solid rgba(240, 192, 96, 0.3)'
      : '1px solid rgba(240, 192, 96, 0.08)',
    borderRadius: 16,
    boxShadow: gold
      ? '0 0 30px rgba(240, 192, 96, 0.15)'
      : '0 4px 24px rgba(0, 0, 0, 0.4)',
    overflow: 'hidden',
    transition: hoverable ? 'all 250ms cubic-bezier(0.4, 0, 0.2, 1)' : 'none',
    cursor: hoverable && onClick ? 'pointer' : 'default',
    ...style,
  }

  const hoverHandlers = hoverable ? {
    onMouseEnter: (e: React.MouseEvent<HTMLDivElement>) => {
      const el = e.currentTarget
      el.style.transform = 'translateY(-4px)'
      el.style.borderColor = gold
        ? 'rgba(240, 192, 96, 0.5)'
        : 'rgba(240, 192, 96, 0.25)'
      el.style.boxShadow = gold
        ? '0 12px 40px rgba(0,0,0,0.5), 0 0 30px rgba(240,192,96,0.25)'
        : '0 12px 40px rgba(0,0,0,0.5), 0 0 20px rgba(240,192,96,0.1)'
    },
    onMouseLeave: (e: React.MouseEvent<HTMLDivElement>) => {
      const el = e.currentTarget
      el.style.transform = 'translateY(0)'
      el.style.borderColor = gold
        ? 'rgba(240, 192, 96, 0.3)'
        : 'rgba(240, 192, 96, 0.08)'
      el.style.boxShadow = gold
        ? '0 0 30px rgba(240, 192, 96, 0.15)'
        : '0 4px 24px rgba(0, 0, 0, 0.4)'
    },
  } : {}

  return (
    <div
      className={className}
      style={baseStyle}
      onClick={onClick}
      {...hoverHandlers}
    >
      {children}
    </div>
  )
}

export default GlassCard
