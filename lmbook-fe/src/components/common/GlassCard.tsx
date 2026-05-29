import React from 'react'

interface GlassCardProps {
  children: React.ReactNode
  className?: string
  hoverable?: boolean
  gold?: boolean
  neon?: boolean
  onClick?: () => void
  style?: React.CSSProperties
}

/** HOK 金色风格玻璃卡片 */
const GlassCard: React.FC<GlassCardProps> = ({
  children, className = '', hoverable = true, gold = false, neon, onClick, style
}) => {
  // Backward compat: neon maps to gold
  const isGold = gold || neon || false

  const baseStyle: React.CSSProperties = {
    background: isGold
      ? 'linear-gradient(135deg, rgba(240,192,96,0.1) 0%, rgba(11,13,23,0.9) 100%)'
      : 'rgba(19, 21, 32, 0.85)',
    backdropFilter: 'blur(24px) saturate(180%)',
    WebkitBackdropFilter: 'blur(24px) saturate(180%)',
    border: isGold
      ? '1px solid rgba(240, 192, 96, 0.5)'
      : '1px solid rgba(240, 192, 96, 0.2)',
    borderRadius: 16,
    boxShadow: isGold
      ? '0 0 20px rgba(240, 192, 96, 0.3)'
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
      el.style.borderColor = isGold
        ? 'rgba(240, 192, 96, 0.8)'
        : 'rgba(240, 192, 96, 0.4)'
      el.style.boxShadow = isGold
        ? '0 12px 40px rgba(0,0,0,0.5), 0 0 30px rgba(240,192,96,0.4)'
        : '0 12px 40px rgba(0,0,0,0.5), 0 0 20px rgba(240,192,96,0.2)'
    },
    onMouseLeave: (e: React.MouseEvent<HTMLDivElement>) => {
      const el = e.currentTarget
      el.style.transform = 'translateY(0)'
      el.style.borderColor = isGold
        ? 'rgba(240, 192, 96, 0.5)'
        : 'rgba(240, 192, 96, 0.2)'
      el.style.boxShadow = isGold
        ? '0 0 20px rgba(240, 192, 96, 0.3)'
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
