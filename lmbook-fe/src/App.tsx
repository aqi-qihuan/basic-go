import React, { Suspense, lazy } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import AppLayout from './components/AppLayout'
import ErrorBoundary from './components/ErrorBoundary'

// 懒加载页面组件 - 实现代码分割
const HomePage = lazy(() => import('./pages/HomePage'))
const ArticleDetailPage = lazy(() => import('./pages/ArticleDetailPage'))
const LoginPage = lazy(() => import('./pages/LoginPage'))
const RegisterPage = lazy(() => import('./pages/RegisterPage'))
const WriteArticlePage = lazy(() => import('./pages/WriteArticlePage'))
const UserProfilePage = lazy(() => import('./pages/UserProfilePage'))
const SettingsPage = lazy(() => import('./pages/SettingsPage'))
const LeaderboardPage = lazy(() => import('./pages/LeaderboardPage'))
const FollowPage = lazy(() => import('./pages/FollowPage'))
const SearchPage = lazy(() => import('./pages/SearchPage'))
const MessagesPage = lazy(() => import('./pages/MessagesPage'))
const CollectionsPage = lazy(() => import('./pages/CollectionsPage'))

// 加载中组件
const LoadingFallback = () => (
  <div className="flex items-center justify-center min-h-[60vh]">
    <div className="text-center">
      {/* 金色旋转加载动画 */}
      <div className="relative w-16 h-16 mx-auto mb-4">
        <div
          className="absolute inset-0 rounded-full"
          style={{
            border: '3px solid rgba(240, 192, 96, 0.2)',
            borderTopColor: '#F0C060',
            animation: 'retroSpin 1s linear infinite',
          }}
        />
        <div
          className="absolute inset-2 rounded-full"
          style={{
            border: '2px solid rgba(59, 130, 246, 0.2)',
            borderBottomColor: '#3B82F6',
            animation: 'retroSpin 1.5s linear infinite reverse',
          }}
        />
      </div>
      <p
        className="text-sm"
        style={{ color: '#9C9688' }}
      >
        加载中...
      </p>
    </div>
  </div>
)

const App: React.FC = () => {
  return (
    <div>
      <ErrorBoundary>
        <Suspense fallback={<LoadingFallback />}>
          <Routes>
            {/* 主布局路由 */}
            <Route element={<AppLayout />}>
              {/* 首页 */}
              <Route path="/" element={<HomePage />} />

              {/* 搜索结果 */}
              <Route path="/search" element={<SearchPage />} />

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

              {/* 关注/粉丝列表 */}
              <Route path="/follow" element={<FollowPage />} />
              <Route path="/follow/:uid" element={<FollowPage />} />

              {/* 设置 */}
              <Route path="/settings" element={<SettingsPage />} />

              {/* 消息中心 */}
              <Route path="/messages" element={<MessagesPage />} />

              {/* 收藏夹 */}
              <Route path="/collections" element={<CollectionsPage />} />

              {/* 默认重定向 */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Route>
          </Routes>
        </Suspense>
      </ErrorBoundary>
    </div>
  )
}

export default App
