import React, { useState } from 'react'
import { Form, Input, Button, Divider, message } from 'antd'
import { UserOutlined, LockOutlined, WechatOutlined, MobileOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { login, getUserProfile } from '@/services/userService'
import { useUserStore } from '@/store/userStore'

/** HOK 营地风格登录页 - 暗黑沉浸 + 金色荣耀 */
const LoginPage: React.FC = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const { setUser, setToken } = useUserStore()
  const navigate = useNavigate()

  const handleLogin = async (values: { email: string; password: string }) => {
    setLoading(true)
    try {
      const response = await login(values.email, values.password)
      // 后端通过响应头 x-jwt-token 返回 Token
      const token = response.headers?.['x-jwt-token']
      if (token) {
        localStorage.setItem('token', token)
        setToken(token)
      }
      // 获取用户信息
      const profile = await getUserProfile()
      setUser(profile)
      message.success('登录成功！')
      navigate('/')
    } catch (error: any) {
      message.error(error.message || '登录失败，请重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="flex items-center justify-center px-4"
      style={{
        minHeight: '100vh',
        background: 'var(--bg-deep)',
      }}
    >
      {/* 荣耀光晕 */}
      <div style={{
        position: 'fixed', top: 0, left: 0, right: 0, height: '50vh',
        background: 'radial-gradient(circle at 50% 30%, rgba(240,192,96,0.06) 0%, transparent 60%)',
        pointerEvents: 'none',
      }} />

      <div className="w-full max-w-md" style={{ position: 'relative', zIndex: 1 }}>
        {/* Logo */}
        <div className="text-center mb-8">
          <h1 style={{
            fontSize: 40, fontWeight: 800, color: '#F0C060',
            margin: '0 0 8px 0',
            textShadow: '0 0 30px rgba(240, 192, 96, 0.3)',
          }}>
            蓝梦社区
          </h1>
          <p style={{ color: '#9C9688', fontSize: 15 }}>登录后享受完整功能</p>
        </div>

        {/* 登录卡片 */}
        <div className="glass-card" style={{ padding: '36px 32px', cursor: 'default' }}>
          <Form
            form={form}
            name="login"
            onFinish={handleLogin}
            layout="vertical"
            size="large"
          >
            <Form.Item
              name="email"
              rules={[
                { required: true, message: '请输入邮箱' },
                { type: 'email', message: '请输入有效的邮箱地址' },
              ]}
            >
              <Input
                prefix={<UserOutlined style={{ color: '#6B6558' }} />}
                placeholder="邮箱"
                style={{
                  background: 'rgba(26, 29, 43, 0.8)',
                  border: '1px solid rgba(240, 192, 96, 0.1)',
                  borderRadius: 12,
                  height: 48,
                  color: '#E8E0D0',
                }}
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[{ required: true, message: '请输入密码' }]}
            >
              <Input.Password
                prefix={<LockOutlined style={{ color: '#6B6558' }} />}
                placeholder="密码"
                style={{
                  background: 'rgba(26, 29, 43, 0.8)',
                  border: '1px solid rgba(240, 192, 96, 0.1)',
                  borderRadius: 12,
                  height: 48,
                }}
              />
            </Form.Item>

            <Form.Item>
              <div className="flex items-center justify-between">
                <Form.Item name="remember" valuePropName="checked" noStyle>
                  <label style={{ color: '#9C9688', fontSize: 13, cursor: 'pointer' }}>
                    <input type="checkbox" style={{ marginRight: 6 }} /> 记住我
                  </label>
                </Form.Item>
                <a href="#" style={{ color: '#F0C060', fontSize: 13 }}>忘记密码？</a>
              </div>
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                block
                style={{
                  background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                  border: 'none',
                  borderRadius: 12,
                  height: 48,
                  fontSize: 16,
                  fontWeight: 700,
                  color: '#0B0D17',
                  boxShadow: '0 4px 16px rgba(240, 192, 96, 0.3)',
                }}
              >
                登录
              </Button>
            </Form.Item>
          </Form>

          <Divider style={{ borderColor: 'rgba(240, 192, 96, 0.1)' }}>
            <span style={{ color: '#6B6558', fontSize: 13 }}>其他登录方式</span>
          </Divider>

          <div className="flex justify-center gap-4 mb-6">
            <Button
              shape="circle"
              size="large"
              icon={<WechatOutlined />}
              onClick={() => message.info('微信登录开发中...')}
              style={{
                background: 'rgba(34, 197, 94, 0.1)',
                border: '1px solid rgba(34, 197, 94, 0.3)',
                color: '#22C55E',
              }}
            />
            <Button
              shape="circle"
              size="large"
              icon={<MobileOutlined />}
              onClick={() => message.info('手机号登录开发中...')}
              style={{
                background: 'rgba(59, 130, 246, 0.1)',
                border: '1px solid rgba(59, 130, 246, 0.3)',
                color: '#3B82F6',
              }}
            />
          </div>

          <div className="text-center">
            <span style={{ color: '#6B6558' }}>还没有账号？</span>
            <Link to="/register" style={{ color: '#F0C060', fontWeight: 600, marginLeft: 6 }}>
              立即注册
            </Link>
          </div>
        </div>

        {/* 底部 */}
        <div className="text-center mt-6">
          <p style={{ color: '#6B6558', fontSize: 12 }}>
            登录即表示同意我们的
            <a href="#" style={{ color: '#F0C060', margin: '0 4px' }}>服务条款</a>
            和
            <a href="#" style={{ color: '#F0C060', marginLeft: 4 }}>隐私政策</a>
          </p>
        </div>
      </div>
    </div>
  )
}

export default LoginPage
