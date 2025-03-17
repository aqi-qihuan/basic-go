package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	// luaSetCode 定义了一个 Lua 脚本，用于设置验证码
	luaSetCode string
	// luaVerifyCode 定义了一个 Lua 脚本，用于验证验证码
	luaVerifyCode string
	// ErrCodeSendTooMany 表示发送验证码过于频繁
	ErrCodeSendTooMany = errors.New("发送太频繁")
	// ErrCodeVerifyTooMany 表示验证验证码过于频繁
	ErrCodeVerifyTooMany = errors.New("验证太频繁")
)

// CodeCache 是一个缓存验证码的接口
type CodeCache interface {
	// Set 设置验证码
	Set(ctx context.Context, biz, phone, code string) error
	// Verify 验证验证码
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

// RedisCodeCache 是一个基于 Redis 的验证码缓存实现
type RedisCodeCache struct {
	// cmd 是一个 redis.Cmdable 接口，用于执行 Redis 命令
	cmd redis.Cmdable
}

// NewCodeCache 创建一个新的 RedisCodeCache 实例
func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

// Set 方法用于设置验证码
func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	// 调用 Redis 的 Eval 方法执行 Lua 脚本，设置验证码
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return err
	}
	// 根据 Lua 脚本的返回值判断设置结果
	switch res {
	case -2:
		// 验证码存在，但是没有过期时间
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		// 发送验证码过于频繁
		return ErrCodeSendTooMany
	default:
		// 设置成功
		return nil
	}
}

// Verify 方法用于验证验证码
func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	// 调用 Redis 的 Eval 方法执行 Lua 脚本，验证验证码
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}
	// 根据 Lua 脚本的返回值判断验证结果
	switch res {
	case -2:
		// 验证码不存在
		return false, nil
	case -1:
		// 验证验证码过于频繁
		return false, ErrCodeVerifyTooMany
	default:
		// 验证成功
		return false, nil
	}
}

// key 方法用于生成 Redis 键
func (c *RedisCodeCache) key(biz, phone string) string {
	// 使用 fmt.Sprintf 格式化字符串，生成 Redis 键
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
