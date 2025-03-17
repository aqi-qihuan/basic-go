package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaScript string

// 嵌入Lua脚本文件，用于在Redis中执行滑动窗口限流逻辑
type RedisSlidingWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

// 定义一个结构体RedisSlidingWindowLimiter，用于实现Redis滑动窗口限流器
// cmd: Redis命令接口
// interval: 时间窗口间隔
// rate: 在时间窗口内的最大请求次数
func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

// 构造函数，用于创建一个新的Redis滑动窗口限流器实例
// cmd: Redis命令接口
// interval: 时间窗口间隔
// rate: 在时间窗口内的最大请求次数
// 返回一个新的RedisSlidingWindowLimiter实例
func (b *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return b.cmd.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
