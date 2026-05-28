import React from 'react'

interface RankBadgeProps {
  rank: number
  size?: 'sm' | 'md' | 'lg'
}

/** HOK 风格排行榜排名徽章 */
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
        className="animate-gold-pulse"
        style={{
          width: s.w, height: s.h, borderRadius: '50%',
          background: 'linear-gradient(135deg, #FFD700 0%, #FFA500 100%)',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          fontSize: s.font, fontWeight: 900, color: '#0B0D17',
          boxShadow: '0 0 20px rgba(255, 215, 0, 0.4)',
          border: '2px solid #FFD700',
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
        background: 'linear-gradient(135deg, #C0C0C0 0%, #A8A8A8 100%)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        fontSize: s.font, fontWeight: 900, color: '#0B0D17',
        boxShadow: '0 0 15px rgba(192, 192, 192, 0.3)',
        border: '2px solid #C0C0C0',
      }}>
        2
      </div>
    )
  }

  if (rank === 3) {
    return (
      <div style={{
        width: s.w, height: s.h, borderRadius: '50%',
        background: 'linear-gradient(135deg, #CD7F32 0%, #B8860B 100%)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        fontSize: s.font, fontWeight: 900, color: '#0B0D17',
        boxShadow: '0 0 15px rgba(205, 127, 50, 0.3)',
        border: '2px solid #CD7F32',
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
      border: '1px solid rgba(240, 192, 96, 0.15)',
    }}>
      {rank}
    </div>
  )
}

export default RankBadge
