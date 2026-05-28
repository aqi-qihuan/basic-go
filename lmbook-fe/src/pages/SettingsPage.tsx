import React, { useState } from 'react'
import { Form, Input, Button, Avatar, message } from 'antd'
import { UserOutlined, MailOutlined, SaveOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useUserStore } from '@/store/userStore'
import { editUserProfile } from '@/services/userService'
import { GlassCard } from '@/components/common'

/** HOK 营地风格设置页 */
const SettingsPage: React.FC = () => {
  const { user, isLogin, setUser } = useUserStore()
  const navigate = useNavigate()
  const [saving, setSaving] = useState(false)

  if (!isLogin || !user) {
    navigate('/login')
    return null
  }

  const handleSave = async (values: any) => {
    setSaving(true)
    try {
      await editUserProfile({ nickname: values.nickname, aboutMe: values.bio })
      setUser({ ...user, ...values })
      message.success('个人信息更新成功')
    } catch {
      message.error('更新失败，请重试')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)', padding: '32px 16px' }}>
      <div className="max-w-2xl mx-auto">
        <h1 style={{ fontSize: 28, fontWeight: 800, color: '#F0C060', marginBottom: 24 }}>
          账号设置
        </h1>

        {/* 头像卡片 */}
        <GlassCard className="mb-6" style={{ padding: 24, cursor: 'default' }}>
          <div className="flex items-center gap-6">
            <Avatar
              size={80}
              src={user.avatar}
              icon={<UserOutlined />}
              style={{
                border: '3px solid rgba(240, 192, 96, 0.3)',
                background: 'rgba(240, 192, 96, 0.1)',
              }}
            />
            <div>
              <h2 style={{ fontSize: 20, fontWeight: 700, color: '#F5F0E8', margin: '0 0 4px 0' }}>
                {user.nickname || user.email}
              </h2>
              <p style={{ color: '#6B6558', fontSize: 14, margin: 0 }}>个人资料</p>
            </div>
          </div>
        </GlassCard>

        {/* 编辑表单 */}
        <GlassCard style={{ padding: '28px 24px', cursor: 'default' }}>
          <h3 style={{ fontSize: 18, fontWeight: 700, color: '#E8E0D0', marginBottom: 24 }}>
            编辑个人资料
          </h3>

          <Form
            layout="vertical"
            initialValues={{ nickname: user.nickname || '', bio: user.bio || '' }}
            onFinish={handleSave}
          >
            <Form.Item label={<span style={{ color: '#9C9688' }}>昵称</span>} name="nickname"
              rules={[{ required: true, message: '请输入昵称' }]}>
              <Input
                prefix={<UserOutlined style={{ color: '#6B6558' }} />}
                placeholder="输入你的昵称"
                style={{
                  background: 'rgba(26, 29, 43, 0.8)', border: '1px solid rgba(240, 192, 96, 0.1)',
                  borderRadius: 12, height: 44, color: '#E8E0D0',
                }}
              />
            </Form.Item>

            <Form.Item label={<span style={{ color: '#9C9688' }}>邮箱</span>}>
              <Input
                prefix={<MailOutlined style={{ color: '#6B6558' }} />}
                disabled
                value={user.email}
                style={{
                  background: 'rgba(26, 29, 43, 0.5)', border: '1px solid rgba(240, 192, 96, 0.06)',
                  borderRadius: 12, height: 44, color: '#6B6558',
                }}
              />
            </Form.Item>

            <Form.Item label={<span style={{ color: '#9C9688' }}>个人简介</span>} name="bio">
              <Input.TextArea
                rows={4}
                placeholder="介绍一下自己..."
                maxLength={200}
                showCount
                style={{
                  background: 'rgba(26, 29, 43, 0.8)', border: '1px solid rgba(240, 192, 96, 0.1)',
                  borderRadius: 12, color: '#E8E0D0', resize: 'none',
                }}
              />
            </Form.Item>

            <div className="gold-line mb-6" />

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                icon={<SaveOutlined />}
                loading={saving}
                style={{
                  background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                  border: 'none', borderRadius: 12, height: 44,
                  fontWeight: 700, color: '#0B0D17',
                  boxShadow: '0 4px 16px rgba(240, 192, 96, 0.3)',
                }}
              >
                保存修改
              </Button>
            </Form.Item>
          </Form>
        </GlassCard>
      </div>
    </div>
  )
}

export default SettingsPage
