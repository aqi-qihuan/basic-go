//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(lanmengbook-record-mysql:3308)/lanmengbook",
	},
	Redis: RedisConfig{
		Addr: "lanmengbook-record-redis:6379",
	},
}
