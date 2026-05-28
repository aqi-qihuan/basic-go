import React from 'react'

interface TagPillProps {
  children: React.ReactNode
  active?: boolean
  tech?: boolean
  onClick?: () => void
  className?: string
}

/** HOK 风格标签胶囊 */
const TagPill: React.FC<TagPillProps> = ({
  children, active = false, tech = false, onClick, className = ''
}) => {
  const baseStyle: React.CSSProperties = {
    display: 'inline-flex',
    alignItems: 'center',
    padding: '4px 14px',
    borderRadius: 20,
    fontSize: 12,
    fontWeight: active ? 700 : 500,
    cursor: onClick ? 'pointer' : 'default',
    transition: 'all 200ms ease',
    userSelect: 'none',
    ...(active ? {
      background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
      color: '#0B0D17',
      border: '1px solid transparent',
    } : tech ? {
      background: 'rgba(59, 130, 246, 0.1)',
      border: '1px solid rgba(59, 130, 246, 0.2)',
      color: '#60A5FA',
    } : {
      background: 'rgba(240, 192, 96, 0.1)',
      border: '1px solid rgba(240, 192, 96, 0.2)',
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
