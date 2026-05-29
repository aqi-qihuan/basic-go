import React, { useState } from 'react'
import { Form, Input, Button, Divider, message, Segmented } from 'antd'
import { UserOutlined, LockOutlined, WechatOutlined, MobileOutlined, SendOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { login, loginBySms, sendSmsCode, getUserProfile, getWechatAuthUrl } from '@/services/userService'
import { useUserStore } from '@/store/userStore'

type LoginMethod = 'email' | 'sms'

/** HOK 营地风格登录页 - 邮箱/手机号双模式 */
const LoginPage: React.FC = () => {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [method, setMethod] = useState<LoginMethod>('email')
  const [countdown, setCountdown] = useState(0)
  const { setUser, setToken } = useUserStore()
  const navigate = useNavigate()

  // 邮箱登录
  const handleEmailLogin = async (values: { email: string; password: string }) => {
    setLoading(true)
    try {
      await login(values.email, values.password)
      const token = localStorage.getItem('token')
      if (token) setToken(token)
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

  // 手机号登录
  const handleSmsLogin = async (values: { phone: string; code: string }) => {
    setLoading(true)
    try {
      await loginBySms(values.phone, values.code)
      const token = localStorage.getItem('token')
      if (token) setToken(token)
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

  // 发送验证码
  const handleSendSms = async () => {
    const phone = form.getFieldValue('phone')
    if (!phone) { message.warning('请先输入手机号'); return }
    try {
      await sendSmsCode(phone)
      message.success('验证码已发送')
      setCountdown(60)
      const timer = setInterval(() => {
        setCountdown(prev => {
          if (prev <= 1) { clearInterval(timer); return 0 }
          return prev - 1
        })
      }, 1000)
    } catch {
      message.error('发送失败，请重试')
    }
  }

  const inputStyle = {
    background: 'rgba(26, 29, 43, 0.8)',
    border: '1px solid rgba(240, 192, 96, 0.1)',
    borderRadius: 12,
    height: 48,
    color: '#E8E0D0',
  }

  return (
    <div className="flex items-center justify-center px-3 sm:px-4 py-8" style={{ minHeight: '100vh', background: 'var(--bg-deep)' }}>
      <div style={{
        position: 'fixed', top: 0, left: 0, right: 0, height: '50vh',
        background: 'radial-gradient(circle at 50% 30%, rgba(240,192,96,0.06) 0%, transparent 60%)',
        pointerEvents: 'none',
      }} />

      <div className="w-full max-w-md" style={{ position: 'relative', zIndex: 1 }}>
        <div className="text-center mb-8">
          <h1 style={{ fontSize: 'clamp(28px, 8vw, 40px)', fontWeight: 800, color: '#F0C060', margin: '0 0 8px 0', textShadow: '0 0 30px rgba(240,192,96,0.3)' }}>
            蓝梦社区
          </h1>
          <p style={{ color: '#9C9688', fontSize: 15 }}>登录后享受完整功能</p>
        </div>

        <div className="glass-card" style={{ padding: '24px 20px', cursor: 'default' }}>
          {/* 登录方式切换 */}
          <div className="flex justify-center mb-6">
            <Segmented
              value={method}
              onChange={(v) => { setMethod(v as LoginMethod); form.resetFields() }}
              options={[
                { value: 'email', label: '邮箱登录' },
                { value: 'sms', label: '手机号登录' },
              ]}
              style={{
                background: 'rgba(26, 29, 43, 0.8)',
                border: '1px solid rgba(240, 192, 96, 0.1)',
              }}
            />
          </div>

          <Form form={form} name="login" onFinish={method === 'email' ? handleEmailLogin : handleSmsLogin as any} layout="vertical" size="large">
            {method === 'email' ? (
              <>
                <Form.Item name="email" rules={[{ required: true, message: '请输入邮箱' }, { type: 'email', message: '请输入有效的邮箱' }]}>
                  <Input prefix={<UserOutlined style={{ color: '#6B6558' }} />} placeholder="邮箱" style={inputStyle} />
                </Form.Item>
                <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                  <Input.Password prefix={<LockOutlined style={{ color: '#6B6558' }} />} placeholder="密码" style={inputStyle} />
                </Form.Item>
              </>
            ) : (
              <>
                <Form.Item name="phone" rules={[{ required: true, message: '请输入手机号' }]}>
                  <Input prefix={<MobileOutlined style={{ color: '#6B6558' }} />} placeholder="手机号" style={inputStyle} />
                </Form.Item>
                <Form.Item>
                  <div className="flex gap-3">
                    <Form.Item name="code" rules={[{ required: true, message: '请输入验证码' }]} noStyle>
                      <Input placeholder="验证码" style={{ ...inputStyle, flex: 1 }} />
                    </Form.Item>
                    <Button
                      icon={<SendOutlined />}
                      disabled={countdown > 0}
                      onClick={handleSendSms}
                      style={{ height: 48, borderRadius: 12, background: 'transparent', border: '1px solid rgba(240,192,96,0.3)', color: '#F0C060' }}
                    >
                      {countdown > 0 ? `${countdown}s` : '发送'}
                    </Button>
                  </div>
                </Form.Item>
              </>
            )}

            {method === 'email' && (
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
            )}

            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading} block style={{
                background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                border: 'none', borderRadius: 12, height: 48, fontSize: 16, fontWeight: 700, color: '#0B0D17',
                boxShadow: '0 4px 16px rgba(240, 192, 96, 0.3)',
              }}>
                登录
              </Button>
            </Form.Item>
          </Form>

          <Divider style={{ borderColor: 'rgba(240, 192, 96, 0.1)' }}>
            <span style={{ color: '#6B6558', fontSize: 13 }}>其他登录方式</span>
          </Divider>

          <div className="flex justify-center gap-4 mb-6">
            <Button shape="circle" size="large" icon={<WechatOutlined />}
              onClick={async () => {
                try {
                  const url = await getWechatAuthUrl()
                  window.location.href = url
                } catch {
                  message.error('获取微信授权链接失败')
                }
              }}
              style={{ background: 'rgba(34,197,94,0.1)', border: '1px solid rgba(34,197,94,0.3)', color: '#22C55E' }}
            />
          </div>

          <div className="text-center">
            <span style={{ color: '#6B6558' }}>还没有账号？</span>
            <Link to="/register" style={{ color: '#F0C060', fontWeight: 600, marginLeft: 6 }}>立即注册</Link>
          </div>
        </div>

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
