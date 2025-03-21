//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(lmbook-record-mysql:3308)/lmbook",
	},
	Redis: RedisConfig{
		Addr: "lmbook-record-redis:6379",
	},
}
