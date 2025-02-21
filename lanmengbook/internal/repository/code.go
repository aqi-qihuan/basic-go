package repository

import (
	"basic-go/lanmengbook/internal/repository/cache"
	"context"
)

// ErrCodeVerifyTooMany 表示验证验证码过于频繁
var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

// CodeRepository 是验证码存储库的接口
type CodeRepository interface {
	// Set 设置验证码
	Set(ctx context.Context, biz, phone, code string) error
	// Verify 验证验证码
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

// CachedCodeRepository 是基于缓存的验证码存储库实现
type CachedCodeRepository struct {
	// cache 是一个 CodeCache 接口，用于缓存验证码
	cache cache.CodeCache
}

// NewCodeRepository 创建一个新的 CodeRepository 实例
func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

// Set 方法用于设置验证码
func (c *CachedCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	// 调用缓存的 Set 方法设置验证码
	return c.cache.Set(ctx, biz, phone, code)
}

// Verify 方法用于验证验证码
func (c *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 调用缓存的 Verify 方法验证验证码
	return c.cache.Verify(ctx, biz, phone, code)
}
