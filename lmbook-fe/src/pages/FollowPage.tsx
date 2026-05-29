import React, { useState, useEffect } from 'react'
import { Avatar, Button, Spin, Tabs, message } from 'antd'
import { UserOutlined, MinusOutlined } from '@ant-design/icons'
import { useParams, Link } from 'react-router-dom'
import { getFolloweeList, getFollowerList, unfollowUser } from '@/services/socialService'
import { useUserStore } from '@/store/userStore'
import { GlassCard, EmptyState } from '@/components/common'

interface FollowRelation {
  id: number
  follower: number
  followee: number
  nickname?: string
  avatar?: string
}

/** HOK 风格关注/粉丝列表页 */
const FollowPage: React.FC = () => {
  const { uid } = useParams<{ uid: string }>()
  const [activeTab, setActiveTab] = useState('followee')
  const [followees, setFollowees] = useState<FollowRelation[]>([])
  const [followers, setFollowers] = useState<FollowRelation[]>([])
  const [loading, setLoading] = useState(false)
  const { user } = useUserStore()

  const userId = uid ? Number(uid) : user?.id

  useEffect(() => {
    if (userId) fetchData()
  }, [userId, activeTab])

  const fetchData = async () => {
    if (!userId) return
    setLoading(true)
    try {
      if (activeTab === 'followee') {
        const res = await getFolloweeList(userId)
        setFollowees(res?.follow_relations || [])
      } else {
        const res = await getFollowerList(userId)
        setFollowers(res?.follow_relations || [])
      }
    } catch {} finally { setLoading(false) }
  }

  const handleUnfollow = async (targetId: number) => {
    try {
      await unfollowUser(targetId)
      message.success('已取消关注')
    } catch { message.error('操作失败') }
  }

  const renderUserList = (list: FollowRelation[], type: 'followee' | 'follower') => {
    if (loading) return <div className="loading-spinner"><Spin size="large" /></div>
    if (list.length === 0) {
      return <EmptyState
        title={type === 'followee' ? '还没有关注任何人' : '还没有粉丝'}
        description={type === 'followee' ? '去发现感兴趣的内容创作者吧' : '发布优质内容吸引更多粉丝'}
      />
    }

    return (
      <div className="flex flex-col gap-3">
        {list.map((relation) => {
          const targetId = type === 'followee' ? relation.followee : relation.follower
          return (
            <GlassCard key={relation.id} className="p-3 sm:p-4" style={{ cursor: 'default' }}>
              <div className="flex items-center justify-between">
                <Link to={`/profile/${targetId}`} className="flex items-center gap-3 no-underline">
                  <Avatar
                    size={44}
                    src={relation.avatar}
                    icon={<UserOutlined />}
                    style={{ border: '2px solid rgba(240,192,96,0.2)', background: 'rgba(240,192,96,0.1)' }}
                  />
                  <span style={{ color: '#E8E0D0', fontWeight: 600, fontSize: 15 }}>
                    {relation.nickname || `用户 ${targetId}`}
                  </span>
                </Link>
                {targetId !== user?.id && (
                  <Button
                    size="small"
                    icon={<MinusOutlined />}
                    onClick={() => handleUnfollow(targetId)}
                    className="h-9"
                    style={{
                      background: 'rgba(240,192,96,0.1)',
                      border: '1px solid rgba(240,192,96,0.3)',
                      color: '#F0C060',
                      borderRadius: 8,
                    }}
                  >
                    取消关注
                  </Button>
                )}
              </div>
            </GlassCard>
          )
        })}
      </div>
    )
  }

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)' }}>
      <div className="max-w-2xl mx-auto px-3 sm:px-4 py-6 sm:py-8">
        <h1 className="text-xl sm:text-2xl" style={{ fontWeight: 800, color: '#F0C060', marginBottom: 24 }}>
          {uid ? '他的关注' : '我的关注'}
        </h1>

        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={[
            {
              key: 'followee',
              label: `关注 (${followees.length})`,
              children: renderUserList(followees, 'followee'),
            },
            {
              key: 'follower',
              label: `粉丝 (${followers.length})`,
              children: renderUserList(followers, 'follower'),
            },
          ]}
        />
      </div>
    </div>
  )
}

export default FollowPage
