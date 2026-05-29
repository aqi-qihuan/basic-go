# 蓝梦社区 (lmbook) 系统架构设计文档

> 最后更新: 2026-05-30 | 版本: v1.0

---

## 1. 项目概述

**蓝梦社区**是一个面向开发者的技术社区平台（类掘金/CSDN），采用 Go 微服务架构后端 + React 前端的全栈技术方案。

### 技术栈

| 层级 | 技术选型 |
|------|---------|
| 前端 | React 18 + TypeScript + Tailwind CSS + Vite |
| 后端网关 | Gin (BFF) + JWT 鉴权 |
| 微服务 | Go + gRPC + GORM |
| 服务发现 | etcd |
| 消息队列 | Kafka |
| 缓存 | Redis |
| 数据库 | MySQL |
| 搜索引擎 | Elasticsearch |
| 容器化 | Docker + Kubernetes |

---

## 2. 系统架构

### 2.1 整体架构（四层）

```
┌─────────────────────────────────────────────────┐
│  前端层 (React + TypeScript + Tailwind)          │
│  12 个页面 · 11 个组件 · 6 个 Service            │
├─────────────────────────────────────────────────┤
│  BFF 网关层 (Gin HTTP)                           │
│  JWT 鉴权 · 路由分发 · gRPC 调用                 │
├─────────────────────────────────────────────────┤
│  微服务层 (16 个独立服务)                         │
│  user · article · comment · interactive · follow │
│  search · tag · feed · ranking · reward          │
│  code · sms · oauth2 · payment · cronjob · im    │
├─────────────────────────────────────────────────┤
│  存储层                                          │
│  MySQL · Redis · Kafka · Elasticsearch           │
└─────────────────────────────────────────────────┘
```

### 2.2 微服务分类

#### 核心业务服务（5 个）

| 服务 | 职责 | 主要接口 |
|------|------|---------|
| **user** | 用户注册/登录/资料管理 | Signup, Login, Profile, Edit |
| **article** | 文章 CRUD/发布/撤回 | Publish, Edit, Detail, List, Withdraw |
| **comment** | 评论/回复/嵌套评论 | CreateComment, GetCommentList, DeleteComment |
| **interactive** | 互动数据（点赞/收藏/打赏/阅读计数） | Like, Collect, Get, GetByIds |
| **follow** | 关注/粉丝关系 | Follow, CancelFollow, GetFollowee, GetFollower |

#### 辅助业务服务（5 个）

| 服务 | 职责 |
|------|------|
| **search** | 全文搜索（ES 索引） |
| **tag** | 标签管理 |
| **feed** | 关注动态时间线 |
| **ranking** | 排行榜计算 |
| **reward** | 打赏/支付详情 |

#### 基础设施服务（6 个）

| 服务 | 职责 |
|------|------|
| **code** | 验证码生成/校验 |
| **sms** | 短信发送 |
| **oauth2** | 第三方登录（微信） |
| **payment** | 支付处理 |
| **cronjob** | 定时任务 |
| **im** | 实时消息（开发中） |

### 2.3 服务间通信

- **同步调用**: gRPC（通过 etcd 服务发现）
- **异步消息**: Kafka 事件驱动
- **服务注册**: etcd

---

## 3. 前端架构

### 3.1 目录结构

```
src/
├── App.tsx              # 路由入口（React.lazy 代码分割）
├── components/          # 组件层
│   ├── AppLayout.tsx    # 全局布局（Header + Content + Footer）
│   ├── ArticleCard.tsx  # 文章卡片
│   ├── CommentForm.tsx  # 评论表单
│   ├── CommentList.tsx  # 评论列表（嵌套回复）
│   ├── ErrorBoundary.tsx # 错误边界
│   ├── editor.tsx       # 富文本编辑器（wangEditor）
│   └── common/          # 公共组件
│       ├── GlassCard.tsx
│       ├── TagPill.tsx
│       ├── RankBadge.tsx
│       ├── RewardModal.tsx
│       └── EmptyState.tsx
├── pages/               # 页面层（12 个，全部懒加载）
│   ├── HomePage.tsx
│   ├── ArticleDetailPage.tsx
│   ├── WriteArticlePage.tsx
│   ├── LoginPage.tsx
│   ├── RegisterPage.tsx
│   ├── UserProfilePage.tsx
│   ├── SettingsPage.tsx
│   ├── LeaderboardPage.tsx
│   ├── FollowPage.tsx
│   ├── SearchPage.tsx
│   ├── MessagesPage.tsx
│   └── CollectionsPage.tsx
├── services/            # API 服务层
│   ├── articleService.ts
│   ├── userService.ts
│   ├── commentService.ts
│   ├── socialService.ts
│   ├── searchService.ts
│   └── rewardService.ts
├── store/               # 状态管理（Zustand）
│   └── userStore.ts
├── types/               # TypeScript 类型定义
└── utils/               # 工具函数
    └── request.ts       # axios 封装
```

### 3.2 路由表

| 路由 | 页面 | 功能 |
|------|------|------|
| `/` | HomePage | 首页文章列表 + Feed |
| `/article/:id` | ArticleDetailPage | 文章详情 + 评论 |
| `/write` | WriteArticlePage | 写文章 |
| `/edit/:id` | WriteArticlePage | 编辑文章 |
| `/login` | LoginPage | 登录（邮箱/手机/微信） |
| `/register` | RegisterPage | 注册 |
| `/profile` | UserProfilePage | 个人中心 |
| `/profile/:uid` | UserProfilePage | 他人主页 |
| `/search` | SearchPage | 搜索结果 |
| `/leaderboard` | LeaderboardPage | 排行榜 |
| `/follow` | FollowPage | 关注/粉丝列表 |
| `/settings` | SettingsPage | 账号设置 |
| `/messages` | MessagesPage | 消息中心 |
| `/collections` | CollectionsPage | 收藏夹 |

### 3.3 性能优化

- **代码分割**: React.lazy + Suspense 实现路由级别懒加载
- **分包策略**: react-vendor / antd-vendor / utils-vendor / editor-vendor
- **页面组件**: 全部 < 16 kB
- **ErrorBoundary**: 全局错误捕获

---

## 4. 数据流

### 4.1 文章发布流程

```
用户填写标题/内容/标签
    ↓
POST /articles/publish → BFF (JWT 鉴权)
    ↓ gRPC
article 服务 → MySQL (articles 表)
    ↓ Kafka 事件
├── search 服务 → ES 索引
├── tag 服务 → 标签关联
└── feed 服务 → 粉丝时间线
```

### 4.2 文章阅读流程

```
用户访问 /article/:id
    ↓
BFF → article 服务 (gRPC) → MySQL 查询文章
    ↓
BFF → interactive 服务 (gRPC) → Redis 缓存查询互动数据
    ↓
BFF → comment 服务 (gRPC) → MySQL 查询评论
    ↓
返回聚合数据给前端
```

### 4.3 用户认证流程

```
登录请求 → BFF
    ↓
user 服务 → MySQL 验证
    ↓
JWT Token 写入 Redis
    ↓
响应头 x-jwt-token 返回前端
    ↓
前端存储 Token，后续请求携带
```

---

## 5. 数据库设计

### 核心表结构

| 表名 | 所属服务 | 主要字段 |
|------|---------|---------|
| users | user | id, email, password, nickname, avatar |
| articles | article | id, title, content, author_id, status |
| comments | comment | id, article_id, user_id, content, parent_id |
| interactives | interactive | biz, biz_id, read_cnt, like_cnt, collect_cnt |
| user_like_bizs | interactive | uid, biz, biz_id, status |
| user_collection_bizs | interactive | uid, biz, biz_id, cid |
| follow_relations | follow | follower, followee, status |
| tags | tag | id, name |
| biz_tags | tag | biz, biz_id, tag_id |

---

## 6. API 路由

### BFF HTTP 路由

| 路径 | 方法 | Handler | 功能 |
|------|------|---------|------|
| `/users/signup` | POST | UserHandler | 注册 |
| `/users/login` | POST | UserHandler | 登录 |
| `/users/logout` | POST | UserHandler | 退出 |
| `/users/profile` | GET | UserHandler | 获取资料 |
| `/users/edit` | POST | UserHandler | 编辑资料 |
| `/users/login_sms/code/send` | POST | UserHandler | 发送验证码 |
| `/users/login_sms` | POST | UserHandler | 手机号登录 |
| `/articles/detail/:id` | GET | ArticleHandler | 文章详情 |
| `/articles/list` | POST | ArticleHandler | 文章列表 |
| `/articles/edit` | POST | ArticleHandler | 编辑文章 |
| `/articles/publish` | POST | ArticleHandler | 发布文章 |
| `/articles/pub/:id` | GET | ArticleHandler | 已发布文章详情 |
| `/articles/pub/like` | POST | ArticleHandler | 点赞 |
| `/articles/pub/collect` | POST | ArticleHandler | 收藏 |
| `/articles/pub/reward` | POST | ArticleHandler | 打赏 |
| `/collections/list` | POST | CollectionHandler | 收藏列表 |
| `/collections/cancel` | POST | CollectionHandler | 取消收藏 |
| `/oauth2/wechat/authurl` | GET | OAuth2Handler | 微信授权 URL |
| `/oauth2/wechat/callback` | GET | OAuth2Handler | 微信回调 |

---

## 7. 部署架构

### Kubernetes 部署

```yaml
# 每个微服务独立 Deployment + Service
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lmbook-article
spec:
  replicas: 2
  template:
    spec:
      containers:
        - name: article
          image: lmbook/article:latest
          ports:
            - containerPort: 8080
```

### 基础设施

- **MySQL**: 主从复制，读写分离
- **Redis**: 集群模式
- **Kafka**: 3 节点集群
- **etcd**: 3 节点集群
- **Elasticsearch**: 3 节点集群

---

## 8. 监控与可观测性

- **日志**: 结构化日志（zap）
- **指标**: Prometheus + Grafana
- **链路追踪**: OpenTelemetry
- **健康检查**: 各服务提供 health endpoint

---

## 9. 设计决策记录

| 决策 | 原因 |
|------|------|
| BFF 模式 | 前端只需调用一个网关，BFF 负责聚合多个微服务 |
| gRPC 同步调用 | 核心路径（读文章、发评论）需要低延迟 |
| Kafka 异步消息 | 非核心路径（搜索索引、Feed 推送）可接受延迟 |
| etcd 服务发现 | 轻量级、高可用、支持 Watch |
| React.lazy 代码分割 | 首屏加载优化，按需加载页面组件 |
| Zustand 状态管理 | 轻量级、TypeScript 友好、无 Provider 包裹 |
| wangEditor 富文本 | 开箱即用、中文友好、体积小 |
