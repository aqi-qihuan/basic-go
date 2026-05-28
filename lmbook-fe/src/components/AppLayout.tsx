import React, { useState } from 'react'
import { Layout, Menu, Button, Drawer, Avatar, Dropdown, Space, Input } from 'antd'
import {
  MenuOutlined, UserOutlined, LogoutOutlined, SettingOutlined,
  EditOutlined, SearchOutlined, TrophyOutlined, HomeOutlined
} from '@ant-design/icons'
import { Link, useNavigate, Outlet, useLocation } from 'react-router-dom'
import { useUserStore } from '@/store/userStore'
import type { MenuProps } from 'antd'

const { Header, Content, Footer } = Layout

interface LayoutProps {
  children?: React.ReactNode
}

const AppLayout: React.FC<LayoutProps> = ({ children }) => {
  const [drawerVisible, setDrawerVisible] = useState(false)
  const { user, isLogin, logout } = useUserStore()
  const navigate = useNavigate()
  const location = useLocation()

  const menuItems: MenuProps['items'] = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: <Link to="/">首页</Link>
    },
    {
      key: '/leaderboard',
      icon: <TrophyOutlined />,
      label: <Link to="/leaderboard">排行榜</Link>
    }
  ]

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile')
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
      onClick: () => navigate('/settings')
    },
    { type: 'divider' },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: () => {
        logout()
        navigate('/login')
      }
    }
  ]

  const handleSearch = (value: string) => {
    if (value.trim()) {
      navigate(`/?q=${encodeURIComponent(value.trim())}`)
    }
  }

  return (
    <Layout className="min-h-screen" style={{ background: 'var(--bg-deep)' }}>
      {/* HOK 风格暗色导航栏 */}
      <Header
        className="sticky top-0 z-50 flex items-center justify-between px-4 md:px-6"
        style={{
          background: 'rgba(11, 13, 23, 0.95)',
          backdropFilter: 'blur(20px)',
          borderBottom: '1px solid rgba(240, 192, 96, 0.08)',
          height: 64,
          lineHeight: '64px',
          padding: '0 24px',
        }}
      >
        {/* Logo */}
        <div className="flex items-center gap-6">
          <Link
            to="/"
            className="flex items-center gap-2 no-underline"
            style={{ textShadow: '0 0 20px rgba(240, 192, 96, 0.3)' }}
          >
            <span
              className="text-xl font-extrabold tracking-wide"
              style={{ color: '#F0C060' }}
            >
              蓝梦社区
            </span>
          </Link>

          {/* 桌面端菜单 */}
          <div className="hidden md:block">
            <Menu
              mode="horizontal"
              items={menuItems}
              selectedKeys={[location.pathname]}
              className="border-none"
              style={{
                background: 'transparent',
                color: '#9C9688',
              }}
            />
          </div>
        </div>

        {/* 搜索框 + 用户操作 */}
        <div className="flex items-center gap-4">
          {/* 搜索框 */}
          <div className="hidden md:block">
            <Input.Search
              placeholder="搜索文章..."
              allowClear
              onSearch={handleSearch}
              style={{ width: 240 }}
              className="hok-search"
            />
          </div>

          {/* 用户操作区 */}
          <div className="hidden md:flex items-center gap-3">
            {isLogin ? (
              <>
                <Button
                  type="primary"
                  icon={<EditOutlined />}
                  onClick={() => navigate('/write')}
                  className="btn-gold"
                  style={{
                    background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                    border: 'none',
                    color: '#0B0D17',
                    fontWeight: 700,
                    borderRadius: 10,
                  }}
                >
                  写文章
                </Button>
                <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                  <Avatar
                    src={user?.avatar}
                    icon={!user?.avatar ? <UserOutlined /> : undefined}
                    className="cursor-pointer"
                    style={{
                      border: '2px solid rgba(240, 192, 96, 0.3)',
                      background: 'rgba(240, 192, 96, 0.1)',
                    }}
                  />
                </Dropdown>
              </>
            ) : (
              <Space>
                <Button
                  type="text"
                  onClick={() => navigate('/login')}
                  style={{ color: '#9C9688' }}
                >
                  登录
                </Button>
                <Button
                  type="primary"
                  onClick={() => navigate('/register')}
                  style={{
                    background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                    border: 'none',
                    color: '#0B0D17',
                    fontWeight: 700,
                    borderRadius: 10,
                  }}
                >
                  注册
                </Button>
              </Space>
            )}
          </div>

          {/* 移动端菜单按钮 */}
          <Button
            type="text"
            icon={<MenuOutlined />}
            className="md:hidden"
            style={{ color: '#E8E0D0' }}
            onClick={() => setDrawerVisible(true)}
          />
        </div>
      </Header>

      {/* 移动端抽屉 */}
      <Drawer
        title={
          <span style={{ color: '#F0C060', fontWeight: 700 }}>蓝梦社区</span>
        }
        placement="right"
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        className="md:hidden"
        style={{ background: '#131520' }}
      >
        <Menu
          mode="vertical"
          items={menuItems}
          selectedKeys={[location.pathname]}
          className="border-none"
          style={{ background: 'transparent' }}
          onClick={() => setDrawerVisible(false)}
        />
        <div className="mt-4 pt-4" style={{ borderTop: '1px solid rgba(240, 192, 96, 0.1)' }}>
          {isLogin ? (
            <Space direction="vertical" className="w-full">
              <Button
                block
                type="primary"
                icon={<EditOutlined />}
                onClick={() => { navigate('/write'); setDrawerVisible(false) }}
                style={{
                  background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                  border: 'none',
                  color: '#0B0D17',
                  fontWeight: 700,
                }}
              >
                写文章
              </Button>
              <Button
                block
                icon={<UserOutlined />}
                onClick={() => { navigate('/profile'); setDrawerVisible(false) }}
                style={{
                  background: 'transparent',
                  borderColor: 'rgba(240, 192, 96, 0.3)',
                  color: '#F0C060',
                }}
              >
                个人中心
              </Button>
              <Button
                block
                danger
                onClick={() => { logout(); setDrawerVisible(false) }}
              >
                退出登录
              </Button>
            </Space>
          ) : (
            <Space direction="vertical" className="w-full">
              <Button
                block
                onClick={() => { navigate('/login'); setDrawerVisible(false) }}
                style={{
                  background: 'transparent',
                  borderColor: 'rgba(240, 192, 96, 0.3)',
                  color: '#F0C060',
                }}
              >
                登录
              </Button>
              <Button
                block
                type="primary"
                onClick={() => { navigate('/register'); setDrawerVisible(false) }}
                style={{
                  background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                  border: 'none',
                  color: '#0B0D17',
                  fontWeight: 700,
                }}
              >
                注册
              </Button>
            </Space>
          )}
        </div>
      </Drawer>

      {/* 主内容区 */}
      <Content className="flex-1" style={{ minHeight: 'calc(100vh - 64px - 60px)' }}>
        {children || <Outlet />}
      </Content>

      {/* HOK 风格页脚 */}
      <Footer
        className="text-center"
        style={{
          background: 'rgba(11, 13, 23, 0.95)',
          borderTop: '1px solid rgba(240, 192, 96, 0.06)',
          padding: '16px 24px',
        }}
      >
        <div className="gold-line mb-3" />
        <span style={{ color: '#6B6558', fontSize: 13 }}>
          &copy;2026 蓝梦社区 - 基于 Go 语言的技术社区平台
        </span>
      </Footer>
    </Layout>
  )
}

export default AppLayout
