//go:build k8s

package config

// Config 定义了应用程序的配置
var Config = config{
	// DB 配置
	DB: DBConfig{
		// DSN 是数据库的连接字符串
		DSN: "root:root@tcp(lanmengbook-record-mysql:3308)/lanmengbook",
	},
	// Redis 配置
	Redis: RedisConfig{
		// Addr 是 Redis 服务器的地址
		Addr: "lanmengbook-record-redis:6379",
	},
}
