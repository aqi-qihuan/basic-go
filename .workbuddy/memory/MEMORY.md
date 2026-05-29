# 蓝梦社区(lmbook)项目长期记忆

## 项目概述
- **类型**: 技术社区论坛（类掘金/CSDN）
- **技术栈**: React + TypeScript + Tailwind CSS + Vite
- **后端**: Go + Gin + gRPC + GORM
- **富文本编辑器**: wangEditor

## 设计风格演变
1. **Phase 1**: HOK 暗黑金色风格（#F0C060 金色主题）
2. **Phase 2**: 复古未来主义风格（#7C3AED 紫色主题 + 霓虹灯效果）
3. **最终决定**: 保留 HOK 金色主题，只保留性能优化部分

## 关键技术决策
- 采用 React.lazy 实现路由级别代码分割
- Vite manualChunks 分包策略：react-vendor / antd-vendor / utils-vendor / editor-vendor
- ErrorBoundary 错误边界组件捕获异常

## 用户偏好
- 喜欢 HOK 金色风格，不喜欢紫色
- 重视性能优化（代码分割、懒加载）
- 使用简体中文交流，混用英文技术术语
- 先规划方案，审批后用"开始编码"触发执行

## 项目结构
- 11 个页面（全部懒加载）
- 6 个业务组件 + 4 个公共组件 + ErrorBoundary
- 6 个服务层（article, user, comment, social, search, reward）

## 构建状态
- `npx vite build` 成功（约 10s）
- 所有页面组件 < 16 kB
- 第三方库独立分包（react/antd/utils/editor）

## 待优化项
- antd-vendor (719KB) 和 editor-vendor (811KB) 体积较大
- 可考虑 Tree Shaking 进一步优化
- 可添加 PWA 支持

## 重要工作记录

### 2026-05-19: 前端开发完成
- TypeScript 类型修复（30+ 个错误）
- 运行时错误修复（CSS 500、vite.svg 404）
- 环境变量问题处理

### 2026-05-21: AqiCloud-Web UI/UX 修复
- BreadCrumb 路径显示修复
- MyShareView 提取码交互优化
- FileTable 颜色统一（粉色→金色）
- RecycleView 空状态 i18n

### 2026-05-28: lmbook-fe HOK 营地风格重构
- 设计系统：tailwind.config.js + index.css 全面改为 HOK 暗黑金色风格
- 全局布局：AppLayout.tsx HOK 暗色导航栏 + 金色 Logo
- 公共组件：GlassCard、RankBadge、TagPill、EmptyState
- API 路径全面修复
- P2 功能全部完成

### 2026-05-29: 风格升级 + 性能优化
- 微信 OAuth2 登录实现
- SettingsPage 增强（通知设置、安全设置）
- 收藏夹全链路实现（后端+前端）
- 移动端全量适配（12个页面）
- 复古未来主义设计系统（最终未采用）

### 2026-05-30: 系统架构分析
- 创建 3 个系统架构图
- 整理项目文档，删除重复文件
