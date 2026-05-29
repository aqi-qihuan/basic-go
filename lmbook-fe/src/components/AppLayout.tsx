import React, { useState } from 'react'
import { Layout, Menu, Button, Drawer, Avatar, Dropdown, Space, Input } from 'antd'
import {
  MenuOutlined, UserOutlined, LogoutOutlined, SettingOutlined,
  EditOutlined, TrophyOutlined, HomeOutlined,
  BellOutlined, StarOutlined
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
    { key: '/', icon: <HomeOutlined />, label: <Link to="/">首页</Link> },
    { key: '/leaderboard', icon: <TrophyOutlined />, label: <Link to="/leaderboard">排行榜</Link> },
  ]

  const userMenuItems: MenuProps['items'] = [
    { key: 'profile', icon: <UserOutlined />, label: '个人中心', onClick: () => navigate('/profile') },
    { key: 'collections', icon: <StarOutlined />, label: '我的收藏', onClick: () => navigate('/collections') },
    { key: 'settings', icon: <SettingOutlined />, label: '设置', onClick: () => navigate('/settings') },
    { type: 'divider' },
    { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: () => { logout(); navigate('/login') } },
  ]

  const handleSearch = (value: string) => {
    if (value.trim()) {
      navigate(`/search?q=${encodeURIComponent(value.trim())}`)
      setDrawerVisible(false)
    }
  }

  const closeDrawer = () => setDrawerVisible(false)

  return (
    <Layout className="min-h-screen" style={{ background: 'var(--bg-deep)' }}>
      {/* HOK 风格暗色导航栏 */}
      <Header
        className="sticky top-0 z-50 flex items-center justify-between"
        style={{
          background: 'rgba(11, 13, 23, 0.95)',
          backdropFilter: 'blur(20px)',
          borderBottom: '1px solid rgba(240, 192, 96, 0.08)',
          height: 56,
          lineHeight: '56px',
          padding: '0 12px',
          boxShadow: 'none',
        }}
      >
        {/* Logo */}
        <div className="flex items-center gap-3 sm:gap-6">
          <Link
            to="/"
            className="flex items-center gap-2 no-underline"
            style={{ textShadow: '0 0 20px rgba(240, 192, 96, 0.3)' }}
          >
            <span
              className="text-lg sm:text-xl font-extrabold tracking-wide font-heading"
              style={{ color: '#F0C060' }}
            >
              蓝梦社区
            </span>
          </Link>

          {/* 桌面端菜单 */}
          <div className="hidden lg:block">
            <Menu
              mode="horizontal"
              items={menuItems}
              selectedKeys={[location.pathname]}
              className="border-none"
              style={{ background: 'transparent', color: '#9C9688' }}
            />
          </div>
        </div>

        {/* 右侧操作区 */}
        <div className="flex items-center gap-2 sm:gap-4">
          {/* 桌面端搜索框 */}
          <div className="hidden lg:block">
            <Input.Search
              placeholder="搜索文章..."
              allowClear
              onSearch={handleSearch}
              style={{ width: 220 }}
              className=""
            />
          </div>

          {/* 桌面端用户操作 */}
          <div className="hidden md:flex items-center gap-2">
            {isLogin ? (
              <>
                <Button
                  type="primary"
                  icon={<EditOutlined />}
                  onClick={() => navigate('/write')}
                  className=""
                  style={{
                    background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                    border: 'none',
                    color: '#0B0D17',
                    fontWeight: 700,
                    borderRadius: 10,
                    height: 36,
                    padding: '0 14px',
                  }}
                >
                  写文章
                </Button>
                <Button
                  type="text"
                  icon={<BellOutlined />}
                  onClick={() => navigate('/messages')}
                  style={{ color: '#9C9688' }}
                />
                <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                  <Avatar
                    src={user?.avatar}
                    icon={!user?.avatar ? <UserOutlined /> : undefined}
                    className="cursor-pointer"
                    size={32}
                    style={{ border: '2px solid rgba(240, 192, 96, 0.3)', background: 'rgba(240, 192, 96, 0.1)' }}
                  />
                </Dropdown>
              </>
            ) : (
              <Space size={4}>
                <Button type="text" onClick={() => navigate('/login')} style={{ color: '#9C9688', height: 36 }}>
                  登录
                </Button>
                <Button
                  type="primary"
                  onClick={() => navigate('/register')}
                  className=""
                  style={{
                    background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                    border: 'none',
                    color: '#0B0D17',
                    fontWeight: 700,
                    borderRadius: 10,
                    height: 36,
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
            style={{ color: '#9C9688', width: 40, height: 40 }}
            onClick={() => setDrawerVisible(true)}
          />
        </div>
      </Header>

      {/* 移动端抽屉 - 包含搜索、导航、用户操作 */}
      <Drawer
        title={<span style={{ color: '#F0C060', fontWeight: 700 }}>蓝梦社区</span>}
        placement="right"
        onClose={closeDrawer}
        open={drawerVisible}
        width={280}
        styles={{
          body: { background: '#0B0D17', padding: '16px' },
          header: { background: '#0B0D17', borderBottom: '1px solid rgba(240, 192, 96, 0.08)' },
        }}
      >
        {/* 移动端搜索框 */}
        <div className="mb-4">
          <Input.Search
            placeholder="搜索文章..."
            allowClear
            onSearch={handleSearch}
            style={{ width: '100%' }}
            className=""
          />
        </div>

        {/* 导航菜单 */}
        <Menu
          mode="vertical"
          items={menuItems}
          selectedKeys={[location.pathname]}
          className="border-none mb-4"
          style={{ background: 'transparent' }}
          onClick={closeDrawer}
        />

        <div style={{ borderTop: '1px solid rgba(240, 192, 96, 0.1)', paddingTop: 16 }}>
          {isLogin ? (
            <Space direction="vertical" className="w-full" size={8}>
              <Button
                block
                type="primary"
                icon={<EditOutlined />}
                onClick={() => { navigate('/write'); closeDrawer() }}
                className=""
                  style={{
                    background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                    border: 'none',
                    color: '#0B0D17',
                    fontWeight: 700,
                    height: 44,
                    borderRadius: 10,
                  }}
                >
                  写文章
              </Button>
              <Button
                block
                icon={<BellOutlined />}
                onClick={() => { navigate('/messages'); closeDrawer() }}
                style={{ background: 'transparent', borderColor: 'rgba(240, 192, 96, 0.1)', color: '#F0C060', height: 44 }}
              >
                消息
              </Button>
              <Button
                block
                icon={<UserOutlined />}
                onClick={() => { navigate('/profile'); closeDrawer() }}
                style={{ background: 'transparent', borderColor: 'rgba(240, 192, 96, 0.1)', color: '#F0C060', height: 44 }}
              >
                个人中心
              </Button>
              <Button
                block
                icon={<TrophyOutlined />}
                onClick={() => { navigate('/leaderboard'); closeDrawer() }}
                style={{ background: 'transparent', borderColor: 'rgba(240, 192, 96, 0.1)', color: '#F0C060', height: 44 }}
              >
                排行榜
              </Button>
              <Button block danger onClick={() => { logout(); closeDrawer() }} style={{ height: 44 }}>
                退出登录
              </Button>
            </Space>
          ) : (
            <Space direction="vertical" className="w-full" size={8}>
              <Button
                block
                onClick={() => { navigate('/login'); closeDrawer() }}
                style={{ background: 'transparent', borderColor: 'rgba(240, 192, 96, 0.1)', color: '#F0C060', height: 44 }}
              >
                登录
              </Button>
              <Button
                block
                type="primary"
                onClick={() => { navigate('/register'); closeDrawer() }}
                className=""
                  style={{
                    background: 'linear-gradient(135deg, #F0C060 0%, #FF8C00 100%)',
                    border: 'none',
                    color: '#0B0D17',
                    fontWeight: 700,
                    height: 44,
                    borderRadius: 10,
                  }}
                >
                  注册
              </Button>
            </Space>
          )}
        </div>
      </Drawer>

      {/* 主内容区 */}
      <Content className="flex-1" style={{ minHeight: 'calc(100vh - 56px - 52px)' }}>
        {children || <Outlet />}
      </Content>

      {/* HOK 风格页脚 */}
      <Footer
        className="text-center"
        style={{
          background: 'rgba(11, 13, 23, 0.95)',
          borderTop: '1px solid rgba(240, 192, 96, 0.06)',
          padding: '12px 16px',
        }}
      >
        <div className="gold-line mb-2" />
        <span style={{ color: '#6B6558', fontSize: 12,  }}>
          &copy;2026 蓝梦社区 - 基于 Go 语言的技术社区平台
        </span>
      </Footer>
    </Layout>
  )
}

export default AppLayout
