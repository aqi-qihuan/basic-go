import React, { useState } from 'react'
import { Tabs, Avatar, Button } from 'antd'
import {
  BellOutlined, LikeOutlined, MessageOutlined, UserOutlined,
  CheckOutlined
} from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { GlassCard, EmptyState } from '@/components/common'

interface MessageItem {
  id: number
  type: 'like' | 'comment' | 'follow' | 'system'
  title: string
  content: string
  from?: { id: number; nickname: string; avatar?: string }
  targetId?: number
  targetType?: 'article' | 'comment'
  read: boolean
  time: string
}

/** HOK 风格消息中心 */
const MessagesPage: React.FC = () => {
  const [activeTab, setActiveTab] = useState('all')

  // 模拟数据（实际应从 API 获取）
  const mockMessages: MessageItem[] = []

  const getIcon = (type: string) => {
    switch (type) {
      case 'like': return <LikeOutlined style={{ color: '#F0C060' }} />
      case 'comment': return <MessageOutlined style={{ color: '#3B82F6' }} />
      case 'follow': return <UserOutlined style={{ color: '#22C55E' }} />
      default: return <BellOutlined style={{ color: '#9C9688' }} />
    }
  }

  const renderMessageList = (list: MessageItem[]) => {
    if (list.length === 0) {
      return <EmptyState icon={<BellOutlined />} title="暂无消息" description="暂时没有新的通知" />
    }

    return (
      <div className="flex flex-col gap-3">
        {list.map(msg => (
          <GlassCard key={msg.id} className="p-3 sm:p-4" style={{ cursor: 'default' }}>
            <div className="flex items-start gap-3">
              <div style={{
                width: 36, height: 36, borderRadius: '50%',
                background: 'rgba(240,192,96,0.1)',
                display: 'flex', alignItems: 'center', justifyContent: 'center',
                flexShrink: 0,
              }}>
                {msg.from?.avatar ? (
                  <Avatar size={36} src={msg.from.avatar} />
                ) : getIcon(msg.type)}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span style={{ color: '#E8E0D0', fontWeight: msg.read ? 400 : 600, fontSize: 14 }}>
                    {msg.title}
                  </span>
                  {!msg.read && (
                    <span style={{
                      width: 6, height: 6, borderRadius: '50%',
                      background: '#F0C060', display: 'inline-block',
                    }} />
                  )}
                </div>
                <p style={{ color: '#9C9688', fontSize: 13, margin: '0 0 4px 0' }}>
                  {msg.content}
                </p>
                <span style={{ color: '#6B6558', fontSize: 12 }}>{msg.time}</span>
              </div>
              {msg.targetId && msg.targetType === 'article' && (
                <Link
                  to={`/article/${msg.targetId}`}
                  style={{ color: '#F0C060', fontSize: 13, flexShrink: 0 }}
                >
                  查看
                </Link>
              )}
            </div>
          </GlassCard>
        ))}
      </div>
    )
  }

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)' }}>
      <div className="max-w-3xl mx-auto px-3 sm:px-4 py-6 sm:py-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-xl sm:text-2xl" style={{ fontWeight: 800, color: '#F0C060', margin: 0 }}>
            消息中心
          </h1>
          <Button
            type="text"
            icon={<CheckOutlined />}
            className="h-9"
            style={{ color: '#9C9688', fontSize: 13 }}
          >
            全部已读
          </Button>
        </div>

        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            { key: 'all', label: '全部', children: renderMessageList(mockMessages) },
            { key: 'like', label: '点赞', children: renderMessageList(mockMessages.filter(m => m.type === 'like')) },
            { key: 'comment', label: '评论', children: renderMessageList(mockMessages.filter(m => m.type === 'comment')) },
            { key: 'follow', label: '关注', children: renderMessageList(mockMessages.filter(m => m.type === 'follow')) },
            { key: 'system', label: '系统', children: renderMessageList(mockMessages.filter(m => m.type === 'system')) },
          ]}
        />
      </div>
    </div>
  )
}

export default MessagesPage
