// 作者信息
export interface Author {
  id: number
  name: string
  avatar?: string
}

// 文章
export interface Article {
  id: number
  title: string
  status: number
  content: string
  abstract: string
  tags?: string[]
  coverImage?: string
  author: Author
  ctime: string
  utime: string
  viewCount?: number
  likeCount?: number
  favoriteCount?: number
  commentCount?: number
}

// 文章列表响应
export interface ArticleListResponse {
  articles: Article[]
  total?: number
}
