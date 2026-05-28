// 评论
export interface Comment {
  id: number
  content: string
  author: CommentAuthor
  articleId: number
  parentId?: number
  replyTo?: CommentAuthor
  ctime: string
  likeCount?: number
}

export interface CommentAuthor {
  id: number
  nickname: string
  avatar?: string
}

// 评论列表响应
export interface CommentListResponse {
  comments: Comment[]
  total: number
}
