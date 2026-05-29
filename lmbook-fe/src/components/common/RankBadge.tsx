import React from 'react'

interface RankBadgeProps {
  rank: number
  size?: 'sm' | 'md' | 'lg'
}

/** HOK 金色风格排行榜排名徽章 */
const RankBadge: React.FC<RankBadgeProps> = ({ rank, size = 'md' }) => {
  const sizeMap = {
    sm: { w: 28, h: 28, font: 12 },
    md: { w: 36, h: 36, font: 16 },
    lg: { w: 48, h: 48, font: 22 },
  }
  const s = sizeMap[size]

  if (rank === 1) {
    return (
      <div
        className="animate-neon-flicker"
        style={{
          width: s.w, height: s.h, borderRadius: '50%',
          background: 'linear-gradient(135deg, #F0C060 0%, #C8982A 100%)',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: s.font, fontWeight: 900, color: '#0B0D17',
          boxShadow: '0 0 20px rgba(240, 192, 96, 0.6)',
          border: '2px solid #9C9688',
        }}
      >
        1
      </div>
    )
  }

  if (rank === 2) {
    return (
      <div style={{
        width: s.w, height: s.h, borderRadius: '50%',
        background: 'linear-gradient(135deg, #3B82F6 0%, #2563EB 100%)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        fontSize: s.font, fontWeight: 900, color: '#FFFFFF',
        boxShadow: '0 0 15px rgba(59, 130, 246, 0.4)',
        border: '2px solid #60A5FA',
      }}>
        2
      </div>
    )
  }

  if (rank === 3) {
    return (
      <div style={{
        width: s.w, height: s.h, borderRadius: '50%',
        background: 'linear-gradient(135deg, #3B82F6 0%, #2563EB 100%)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        fontSize: s.font, fontWeight: 900, color: '#FFFFFF',
        boxShadow: '0 0 15px rgba(59, 130, 246, 0.4)',
        border: '2px solid #60A5FA',
      }}>
        3
      </div>
    )
  }

  return (
    <div style={{
      width: s.w, height: s.h, borderRadius: '50%',
      background: 'rgba(240, 192, 96, 0.1)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      fontSize: s.font, fontWeight: 700, color: '#9C9688',
      border: '1px solid rgba(240, 192, 96, 0.3)',
    }}>
      {rank}
    </div>
  )
}

export default RankBadge
