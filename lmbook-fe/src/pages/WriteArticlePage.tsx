import React, { useState, useEffect, useRef } from 'react'
import { Button, Input, Tag, message, Spin } from 'antd'
import { SaveOutlined, SendOutlined, ArrowLeftOutlined } from '@ant-design/icons'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { Editor, Toolbar } from '@wangeditor/editor-for-react'
import { IDomEditor, IEditorConfig, IToolbarConfig } from '@wangeditor/editor'
import { saveArticle, publishArticle, getArticleById } from '@/services/articleService'
import { useUserStore } from '@/store/userStore'
import '@wangeditor/editor/dist/css/style.css'
import type { Author } from '@/types/article'

const WriteArticlePage: React.FC = () => {
  const [title, setTitle] = useState('')
  const [content, setContent] = useState('')
  const [tags, setTags] = useState<string[]>([])
  const [inputTag, setInputTag] = useState('')
  const [saving, setSaving] = useState(false)
  const [publishing, setPublishing] = useState(false)
  const [loading, setLoading] = useState(false)
  const editorRef = useRef<IDomEditor | null>(null)
  const { user } = useUserStore()
  const navigate = useNavigate()
  const { id } = useParams<{ id: string }>()

  // 编辑模式：加载文章数据
  useEffect(() => {
    if (id) {
      fetchArticle(id)
    }
  }, [id]) // eslint-disable-line react-hooks/exhaustive-deps

  const fetchArticle = async (articleId: string) => {
    setLoading(true)
    try {
      const article = await getArticleById(Number(articleId))
      setTitle(article.title)
      setContent(article.content)
      setTags(article.tags || [])
      if (editorRef.current) {
        editorRef.current.setHtml(article.content)
      }
    } catch {
      message.error('加载文章失败')
    } finally {
      setLoading(false)
    }
  }

  // 编辑器配置
  const editorConfig: Partial<IEditorConfig> = {
    placeholder: '请输入文章内容...',
    onChange(editor: IDomEditor) {
      setContent(editor.getHtml())
    },
    MENU_CONF: {
      uploadImage: {
        server: '/api/upload/image',
        fieldName: 'file',
      },
    },
  }

  const toolbarConfig: Partial<IToolbarConfig> = {
    excludeKeys: [
      'group-video',
    ],
  }

  const handleCreateEditor = (editor: IDomEditor) => {
    editorRef.current = editor
    if (content) {
      editor.setHtml(content)
    }
  }

  // 构建 author 对象
  const getAuthor = (): Author | undefined => {
    if (!user?.id) return undefined
    return {
      id: user.id,
      name: user.nickname || user.email || '',
      avatar: user.avatar,
    }
  }

  // 添加标签
  const handleAddTag = () => {
    if (inputTag && !tags.includes(inputTag)) {
      setTags([...tags, inputTag])
      setInputTag('')
    }
  }

  // 删除标签
  const handleRemoveTag = (removedTag: string) => {
    setTags(tags.filter(tag => tag !== removedTag))
  }

  // 保存草稿
  const handleSaveDraft = async () => {
    if (!title.trim()) {
      message.warning('请输入文章标题')
      return
    }

    setSaving(true)
    try {
      const articleId = await saveArticle({
        title: title.trim(),
        content,
        tags,
        author: getAuthor(),
      })
      message.success('草稿保存成功！')
      if (articleId) {
        navigate(`/edit/${articleId}`)
      }
    } catch {
      message.error('保存失败，请重试')
    } finally {
      setSaving(false)
    }
  }

  // 发布文章
  const handlePublish = async () => {
    if (!title.trim()) {
      message.warning('请输入文章标题')
      return
    }
    if (!content.trim() || content === '<p><br></p>') {
      message.warning('请输入文章内容')
      return
    }

    setPublishing(true)
    try {
      await publishArticle({
        title: title.trim(),
        content,
        tags,
        author: getAuthor(),
      })
      message.success('文章发布成功！')
      navigate('/')
    } catch {
      message.error('发布失败，请重试')
    } finally {
      setPublishing(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Spin size="large" tip="加载中..." />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 顶部工具栏 */}
      <div className="sticky top-0 z-10 bg-white shadow-sm px-4 py-3 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Link to="/" className="text-gray-600 hover:text-primary-600 transition-colors duration-200">
            <ArrowLeftOutlined className="text-xl" />
          </Link>
          <h1 className="text-xl font-semibold text-gray-800">
            {id ? '编辑文章' : '写文章'}
          </h1>
        </div>

        <div className="flex items-center space-x-3">
          <Button
            icon={<SaveOutlined />}
            loading={saving}
            onClick={handleSaveDraft}
            className="border-primary-600 text-primary-600 hover:bg-primary-50 transition-all duration-200"
          >
            保存草稿
          </Button>
          <Button
            type="primary"
            icon={<SendOutlined />}
            loading={publishing}
            onClick={handlePublish}
            className="bg-primary-600 hover:bg-primary-700 transition-all duration-200"
          >
            发布文章
          </Button>
        </div>
      </div>

      {/* 编辑区域 */}
      <div className="max-w-5xl mx-auto px-4 py-8">
        {/* 标题输入 */}
        <div className="mb-6">
          <Input
            placeholder="请输入文章标题..."
            value={title}
            onChange={e => setTitle(e.target.value)}
            className="text-3xl font-bold border-none shadow-none focus:shadow-none text-gray-800"
            style={{ fontSize: '2rem', fontWeight: 700 }}
          />
        </div>

        {/* 标签输入 */}
        <div className="mb-6">
          <div className="flex items-center space-x-2 mb-3">
            <Input
              placeholder="添加标签"
              value={inputTag}
              onChange={e => setInputTag(e.target.value)}
              onPressEnter={handleAddTag}
              className="w-48"
            />
            <Button type="dashed" onClick={handleAddTag}>
              添加标签
            </Button>
          </div>
          <div className="flex flex-wrap gap-2">
            {tags.map(tag => (
              <Tag
                key={tag}
                closable
                onClose={() => handleRemoveTag(tag)}
                className="bg-primary-50 text-primary-600 border-primary-200"
              >
                {tag}
              </Tag>
            ))}
          </div>
        </div>

        {/* 富文本编辑器 */}
        <div className="bg-white rounded-xl shadow-card hover:shadow-card-hover transition-all duration-300">
          <Toolbar
            editor={editorRef.current}
            defaultConfig={toolbarConfig}
            className="border-b"
          />
          <Editor
            defaultConfig={editorConfig}
            value={content}
            onCreated={handleCreateEditor}
            className="min-h-screen-50 p-6"
          />
        </div>
      </div>
    </div>
  )
}

export default WriteArticlePage
