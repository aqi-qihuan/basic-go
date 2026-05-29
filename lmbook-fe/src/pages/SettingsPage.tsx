import React, { useState } from 'react'
import { Form, Input, Button, Avatar, message, Switch, Divider } from 'antd'
import { UserOutlined, MailOutlined, SaveOutlined, BellOutlined, LockOutlined } from '@ant-design/icons'
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
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)' }} className="px-3 sm:px-4 py-4 sm:py-6">
      <div className="max-w-2xl mx-auto">
        <h1 style={{ fontWeight: 800, color: '#F0C060', marginBottom: 24 }}
          className="text-xl sm:text-2xl">
          账号设置
        </h1>

        {/* 头像卡片 */}
        <GlassCard className="mb-4 sm:mb-6" style={{ cursor: 'default' }}>
          <div className="flex items-center gap-4 sm:gap-6">
            <Avatar
              size={64}
              src={user.avatar}
              icon={<UserOutlined />}
              className="sm:w-20 sm:h-20"
              style={{
                border: '3px solid rgba(240, 192, 96, 0.3)',
                background: 'rgba(240, 192, 96, 0.1)',
              }}
            />
            <div>
              <h2 style={{ fontWeight: 700, color: '#F5F0E8', margin: '0 0 4px 0' }}
                className="text-base sm:text-lg">
                {user.nickname || user.email}
              </h2>
              <p style={{ color: '#6B6558', margin: 0 }} className="text-xs sm:text-sm">个人资料</p>
            </div>
          </div>
        </GlassCard>

        {/* 编辑表单 */}
        <GlassCard className="mb-4 sm:mb-6" style={{ cursor: 'default' }}>
          <div className="p-4 sm:p-6">
            <h3 style={{ fontWeight: 700, color: '#E8E0D0', marginBottom: 20 }}
              className="text-base sm:text-lg">
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
                  rows={3}
                  placeholder="介绍一下自己..."
                  maxLength={200}
                  showCount
                  style={{
                    background: 'rgba(26, 29, 43, 0.8)', border: '1px solid rgba(240, 192, 96, 0.1)',
                    borderRadius: 12, color: '#E8E0D0', resize: 'none',
                  }}
                />
              </Form.Item>

              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  icon={<SaveOutlined />}
                  loading={saving}
                  block
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
          </div>
        </GlassCard>

        {/* 通知设置 */}
        <GlassCard className="mb-4 sm:mb-6" style={{ cursor: 'default' }}>
          <div className="p-4 sm:p-6">
            <h3 style={{ fontWeight: 700, color: '#E8E0D0', marginBottom: 20 }}
              className="text-base sm:text-lg">
              <BellOutlined style={{ marginRight: 8, color: '#F0C060' }} />
              通知设置
            </h3>

            <div className="flex flex-col gap-3 sm:gap-4">
              {[
                { title: '评论通知', desc: '有人评论你的文章时通知' },
                { title: '点赞通知', desc: '有人点赞你的文章时通知' },
                { title: '关注通知', desc: '有人关注你时通知' },
                { title: '系统通知', desc: '接收系统公告和更新通知' },
              ].map((item, idx) => (
                <React.Fragment key={item.title}>
                  {idx > 0 && <Divider style={{ borderColor: 'rgba(240, 192, 96, 0.06)', margin: 0 }} />}
                  <div className="flex items-center justify-between">
                    <div className="flex-1 min-w-0 mr-3">
                      <p style={{ color: '#E8E0D0', margin: 0, fontWeight: 500 }} className="text-sm sm:text-base">{item.title}</p>
                      <p style={{ color: '#6B6558', margin: 0 }} className="text-xs sm:text-sm">{item.desc}</p>
                    </div>
                    <Switch defaultChecked style={{ background: '#F0C060' }} />
                  </div>
                </React.Fragment>
              ))}
            </div>
          </div>
        </GlassCard>

        {/* 安全设置 */}
        <GlassCard style={{ cursor: 'default' }}>
          <div className="p-4 sm:p-6">
            <h3 style={{ fontWeight: 700, color: '#E8E0D0', marginBottom: 20 }}
              className="text-base sm:text-lg">
              <LockOutlined style={{ marginRight: 8, color: '#F0C060' }} />
              安全设置
            </h3>

            <div className="flex flex-col gap-3 sm:gap-4">
              {[
                { title: '修改密码', desc: '定期修改密码以保护账号安全', action: '修改' },
                { title: '绑定手机', desc: user.phone ? `已绑定: ${user.phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')}` : '未绑定', action: user.phone ? '更换' : '绑定' },
                { title: '绑定微信', desc: user.wechat ? '已绑定' : '未绑定', action: user.wechat ? '更换' : '绑定' },
              ].map((item, idx) => (
                <React.Fragment key={item.title}>
                  {idx > 0 && <Divider style={{ borderColor: 'rgba(240, 192, 96, 0.06)', margin: 0 }} />}
                  <div className="flex items-center justify-between">
                    <div className="flex-1 min-w-0 mr-3">
                      <p style={{ color: '#E8E0D0', margin: 0, fontWeight: 500 }} className="text-sm sm:text-base">{item.title}</p>
                      <p style={{ color: '#6B6558', margin: 0 }} className="text-xs sm:text-sm">{item.desc}</p>
                    </div>
                    <Button
                      style={{
                        background: 'transparent',
                        border: '1px solid rgba(240, 192, 96, 0.3)',
                        color: '#F0C060',
                        borderRadius: 8,
                        height: 36,
                      }}
                    >
                      {item.action}
                    </Button>
                  </div>
                </React.Fragment>
              ))}
            </div>
          </div>
        </GlassCard>
      </div>
    </div>
  )
}

export default SettingsPage
