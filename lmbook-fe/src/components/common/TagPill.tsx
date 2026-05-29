import React from 'react'

interface TagPillProps {
  children: React.ReactNode
  active?: boolean
  tech?: boolean
  onClick?: () => void
  className?: string
}

/** HOK 金色风格标签胶囊 */
const TagPill: React.FC<TagPillProps> = ({
  children, active = false, tech = false, onClick, className = ''
}) => {
  const baseStyle: React.CSSProperties = {
    display: 'inline-flex',
    alignItems: 'center',
    padding: '4px 14px',
    borderRadius: 10,
    fontSize: 12,
    fontWeight: active ? 700 : 500,
    cursor: onClick ? 'pointer' : 'default',
    transition: 'all 200ms ease',
    userSelect: 'none',
    ...(active ? {
      background: 'linear-gradient(135deg, #F0C060 0%, #C8982A 100%)',
      color: '#0B0D17',
      border: '1px solid #9C9688',
      boxShadow: '0 0 5px rgba(240, 192, 96, 0.3)',
    } : tech ? {
      background: 'rgba(59, 130, 246, 0.1)',
      border: '1px solid rgba(59, 130, 246, 0.3)',
      color: '#3B82F6',
    } : {
      background: 'rgba(240, 192, 96, 0.1)',
      border: '1px solid rgba(240, 192, 96, 0.3)',
      color: '#F0C060',
    }),
  }

  return (
    <span
      className={className}
      style={baseStyle}
      onClick={onClick}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
    >
      {children}
    </span>
  )
}

export default TagPill
