import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import AppLayout from './components/AppLayout'
import HomePage from './pages/HomePage'
import ArticleDetailPage from './pages/ArticleDetailPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import WriteArticlePage from './pages/WriteArticlePage'
import UserProfilePage from './pages/UserProfilePage'
import SettingsPage from './pages/SettingsPage'
import LeaderboardPage from './pages/LeaderboardPage'

const App: React.FC = () => {
  return (
    <Routes>
      {/* 主布局路由 */}
      <Route element={<AppLayout />}>
        {/* 首页 */}
        <Route path="/" element={<HomePage />} />

        {/* 排行榜 */}
        <Route path="/leaderboard" element={<LeaderboardPage />} />

        {/* 文章详情 */}
        <Route path="/article/:id" element={<ArticleDetailPage />} />

        {/* 写文章/编辑文章 */}
        <Route path="/write" element={<WriteArticlePage />} />
        <Route path="/edit/:id" element={<WriteArticlePage />} />

        {/* 用户认证 */}
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        {/* 用户中心 */}
        <Route path="/profile" element={<UserProfilePage />} />
        <Route path="/profile/:uid" element={<UserProfilePage />} />

        {/* 设置 */}
        <Route path="/settings" element={<SettingsPage />} />

        {/* 默认重定向 */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}

export default App
