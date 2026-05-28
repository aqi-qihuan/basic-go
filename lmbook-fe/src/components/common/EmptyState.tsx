import React from 'react'

interface EmptyStateProps {
  icon?: React.ReactNode
  title: string
  description?: string
  action?: React.ReactNode
}

/** HOK 风格空状态 */
const EmptyState: React.FC<EmptyStateProps> = ({ icon, title, description, action }) => {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-4">
      {icon && (
        <div
          className="mb-4"
          style={{
            width: 80, height: 80, borderRadius: '50%',
            background: 'rgba(240, 192, 96, 0.08)',
            border: '1px solid rgba(240, 192, 96, 0.12)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: 36, color: '#F0C060',
          }}
        >
          {icon}
        </div>
      )}
      <h3 style={{ color: '#E8E0D0', fontSize: 18, fontWeight: 600, marginBottom: 8 }}>
        {title}
      </h3>
      {description && (
        <p style={{ color: '#6B6558', fontSize: 14, marginBottom: 16, textAlign: 'center' }}>
          {description}
        </p>
      )}
      {action}
    </div>
  )
}

export default EmptyState
