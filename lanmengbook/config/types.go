package config

// config 结构体定义了应用程序的配置
type config struct {
	// DB 配置
	DB DBConfig
	// Redis 配置
	Redis RedisConfig
}

// DBConfig 结构体定义了数据库的配置
type DBConfig struct {
	// DSN 是数据库的连接字符串
	DSN string
}

// RedisConfig 结构体定义了 Redis 的配置
type RedisConfig struct {
	// Addr 是 Redis 服务器的地址
	Addr string
}
