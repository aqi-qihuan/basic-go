package ioc

import (
	"basic-go/lmbook/config"
	"github.com/redis/go-redis/v9"
)

// InitRedis 初始化并返回一个 Redis 客户端实例
func InitRedis() redis.Cmdable {
	// 创建一个新的 Redis 客户端，使用配置文件中的地址
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr, // 从配置文件中获取 Redis 服务器的地址
	})
}
