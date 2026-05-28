import React, { useEffect, useState, useCallback, useRef } from 'react'
import { Input, Spin, Pagination, message } from 'antd'
import { SearchOutlined, FireOutlined, TrophyOutlined } from '@ant-design/icons'
import ArticleCard from '@/components/ArticleCard'
import { TagPill, EmptyState } from '@/components/common'
import { getPublishedArticleList } from '@/services/articleService'
import { searchArticles, getUserTags } from '@/services/searchService'
import type { Article } from '@/types/article'
import dayjs from 'dayjs'

const { Search } = Input

interface TagItem {
  name: string
  count: number
}

/** HOK 营地风格首页 - 暗黑沉浸 + 金色荣耀 */
const HomePage: React.FC = () => {
  const [articles, setArticles] = useState<Article[]>([])
  const [featuredArticles, setFeaturedArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [pagination, setPagination] = useState({ current: 1, pageSize: 12, total: 0 })
  const [keyword, setKeyword] = useState('')
  const [activeTag, setActiveTag] = useState('')
  const [tags, setTags] = useState<TagItem[]>([])
  const debounceRef = useRef<ReturnType<typeof setTimeout>>()

  const fetchArticles = useCallback(async (page = 1, pageSize = 12, kw = '', tag = '') => {
    setLoading(true)
    try {
      const offset = (page - 1) * pageSize
      let result: { articles: any[]; total?: number }

      if (kw.trim()) {
        result = await searchArticles({ keyword: kw.trim(), offset, limit: pageSize })
      } else if (tag) {
        result = await getArticlesByTag(tag, offset, pageSize)
      } else {
        result = await getPublishedArticleList({
          offset, limit: pageSize,
          start_time: dayjs().subtract(1, 'year').toISOString()
        })
      }

      const arts = result.articles || []
      setArticles(arts)

      // 首次加载取前3篇作为热门推荐
      if (page === 1 && !kw && !tag) {
        setFeaturedArticles(arts.slice(0, 3))
      }

      setPagination(prev => ({
        ...prev,
        current: page,
        pageSize,
        total: result.total ?? (arts.length < pageSize ? offset + arts.length : offset + pageSize + 1)
      }))
    } catch {
      message.error('加载文章失败')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    getUserTags().then((res) => setTags(res?.tags || [])).catch(() => {})
  }, [])

  useEffect(() => { fetchArticles() }, [fetchArticles])

  const handleSearch = (value: string) => {
    setKeyword(value)
    setActiveTag('')
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      fetchArticles(1, pagination.pageSize, value, '')
    }, 300)
  }

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setKeyword(value)
    setActiveTag('')
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      fetchArticles(1, pagination.pageSize, value, '')
    }, 400)
  }

  const handleTagClick = (tagName: string) => {
    const newTag = tagName === activeTag ? '' : tagName
    setActiveTag(newTag)
    setKeyword('')
    fetchArticles(1, pagination.pageSize, '', newTag)
  }

  const handlePageChange = (page: number, pageSize?: number) => {
    fetchArticles(page, pageSize || pagination.pageSize, keyword, activeTag)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  return (
    <div style={{ minHeight: '100vh', background: 'var(--bg-deep)' }}>
      {/* 荣耀光晕背景 */}
      <div className="bg-glory" style={{ paddingTop: 32, paddingBottom: 16 }}>
        <div className="max-w-6xl mx-auto px-4">
          {/* 搜索栏 */}
          <div className="flex justify-center mb-8">
            <Search
              placeholder="搜索文章、作者、标签..."
              value={keyword}
              onChange={handleSearchChange}
              onSearch={handleSearch}
              enterButton={<SearchOutlined />}
              size="large"
              allowClear
              style={{ maxWidth: 560 }}
              className="hok-search"
            />
          </div>

          {/* 热门推荐区（仅首页无搜索时显示） */}
          {!keyword && !activeTag && featuredArticles.length > 0 && (
            <div className="mb-8">
              <div className="flex items-center gap-2 mb-4">
                <FireOutlined style={{ color: '#F0C060', fontSize: 20 }} />
                <h2 style={{ color: '#F0C060', fontSize: 20, fontWeight: 700, margin: 0 }}>
                  热门推荐
                </h2>
                <div className="gold-line flex-1 ml-3" />
              </div>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-5">
                {featuredArticles.map((article) => (
                  <ArticleCard key={article.id} article={article} featured />
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="max-w-6xl mx-auto px-4 py-6">
        {/* 标签筛选 */}
        {tags.length > 0 && (
          <div className="flex items-center flex-wrap gap-2 mb-6">
            <TrophyOutlined style={{ color: '#F0C060', marginRight: 4 }} />
            <TagPill
              active={!activeTag}
              onClick={() => { setActiveTag(''); setKeyword(''); fetchArticles(1, pagination.pageSize, '', '') }}
            >
              全部
            </TagPill>
            {tags.map((tag) => (
              <TagPill
                key={tag.name}
                active={activeTag === tag.name}
                onClick={() => handleTagClick(tag.name)}
              >
                {tag.name} ({tag.count})
              </TagPill>
            ))}
          </div>
        )}

        {/* 区块标题 */}
        <div className="flex items-center justify-between mb-6">
          <h2 style={{ color: '#E8E0D0', fontSize: 22, fontWeight: 700, margin: 0 }}>
            {keyword ? `搜索: ${keyword}` : activeTag ? `标签: ${activeTag}` : '最新文章'}
          </h2>
          <span style={{ color: '#6B6558', fontSize: 13 }}>
            共 {pagination.total} 篇
          </span>
        </div>

        {/* 文章网格 */}
        {loading ? (
          <div className="loading-spinner">
            <Spin size="large" />
          </div>
        ) : articles.length > 0 ? (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5 mb-8">
              {articles.map((article) => (
                <ArticleCard key={article.id} article={article} />
              ))}
            </div>

            {/* 分页 */}
            <div className="flex justify-center py-4">
              <Pagination
                current={pagination.current}
                pageSize={pagination.pageSize}
                total={pagination.total}
                onChange={handlePageChange}
                showSizeChanger
                showQuickJumper
                showTotal={(total) => (
                  <span style={{ color: '#6B6558' }}>共 {total} 篇文章</span>
                )}
              />
            </div>
          </>
        ) : (
          <EmptyState
            icon={<SearchOutlined />}
            title={keyword ? `未找到与 "${keyword}" 相关的文章` : activeTag ? `暂无 "${activeTag}" 标签的文章` : '暂无文章'}
            description="试试其他关键词或标签"
          />
        )}
      </div>
    </div>
  )
}

export default HomePage
