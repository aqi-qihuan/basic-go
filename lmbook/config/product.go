//go:build product

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/lmbook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
