package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

// UserCache 是一个缓存用户信息的接口
type UserCache interface {
	// Get 获取用户信息
	Get(ctx context.Context, uid int64) (domain.User, error)
	// Set 设置用户信息
	Set(ctx context.Context, du domain.User) error
}

// RedisUserCache 是一个基于 Redis 的用户缓存实现
type RedisUserCache struct {
	// cmd 是一个 redis.Cmdable 接口，用于执行 Redis 命令
	cmd redis.Cmdable
	// expiration 是缓存的过期时间
	expiration time.Duration
}

// Get 方法用于获取用户信息
func (c *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	// 生成 Redis 键
	key := c.key(uid)
	// 从 Redis 中获取用户信息
	data, err := c.cmd.Get(ctx, key).Result()
	// 如果获取失败，返回错误
	if err != nil {
		return domain.User{}, err
	}
	// 反序列化用户信息
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	// 如果反序列化失败，返回错误
	return u, err
}

// Set 方法用于设置用户信息
func (c *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	// 生成 Redis 键
	key := c.key(du.Id)
	// 序列化用户信息
	data, err := json.Marshal(du)
	// 如果序列化失败，返回错误
	if err != nil {
		return err
	}
	// 将用户信息保存到 Redis 中
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

// key 方法用于生成 Redis 键
func (c *RedisUserCache) key(uid int64) string {
	// 使用 fmt.Sprintf 格式化字符串，生成 Redis 键
	return fmt.Sprintf("user:info:%d", uid)
}

// UserCacheV1 是一个基于 Redis 的用户缓存实现
type UserCacheV1 struct {
	// client 是一个 redis.Client 实例，用于连接 Redis 服务器
	client *redis.Client
}

// NewUserCache 创建一个新的 UserCache 实例
func NewUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

//func NewUserCacheV1(addr string) *UserCache {
//	cmd := redis.NewClient(&redis.Options{Addr: addr})
//	return &UserCache{
//		cmd:        cmd,
//		expiration: time.Minute * 15,
//	}
//}
