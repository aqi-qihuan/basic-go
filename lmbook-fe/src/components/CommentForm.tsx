import React, { useState } from 'react'
import { Input, Button, message } from 'antd'
import { createComment } from '@/services/commentService'

const { TextArea } = Input

interface CommentFormProps {
  articleId: number
  parentId?: number
  onSuccess: () => void
  placeholder?: string
}

const CommentForm: React.FC<CommentFormProps> = ({ articleId, parentId, onSuccess, placeholder }) => {
  const [content, setContent] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const handleSubmit = async () => {
    if (!content.trim()) {
      message.warning('请输入评论内容')
      return
    }
    setSubmitting(true)
    try {
      await createComment({
        articleId,
        content: content.trim(),
        parentId,
      })
      message.success('评论发表成功')
      setContent('')
      onSuccess()
    } catch {
      message.error('评论发表失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.ctrlKey && e.key === 'Enter') {
      handleSubmit()
    }
  }

  return (
    <div className="comment-form">
      <TextArea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        onKeyDown={handleKeyDown}
        rows={4}
        placeholder={placeholder || '写下你的评论... (Ctrl+Enter 发送)'}
        className="mb-3"
      />
      <div className="flex justify-end">
        <Button
          type="primary"
          onClick={handleSubmit}
          loading={submitting}
          disabled={!content.trim()}
        >
          发表评论
        </Button>
      </div>
    </div>
  )
}

export default CommentForm
