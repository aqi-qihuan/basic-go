package ratelimit

import (
	"basic-go/lmbook/pkg/limiter"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// 定义一个Builder结构体，用于构建限流中间件
type Builder struct {
	prefix  string          // 前缀，用于标识限流的键
	limiter limiter.Limiter // 限流器接口，用于执行限流逻辑
}

// NewBuilder 创建一个新的Builder实例
func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{
		prefix:  "ip-limiter", // 默认前缀为"ip-limiter"
		limiter: l,            // 传入的限流器实例
	}
}

// Prefix 设置Builder的前缀
func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix // 更新前缀
	return b          // 返回Builder实例，支持链式调用
}

// Build 构建并返回一个gin.HandlerFunc中间件
func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 使用限流器进行限流判断，键由前缀和客户端IP组成
		limited, err := b.limiter.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
		if err != nil {
			log.Println(err) // 打印错误信息
			// 这一步很有意思，就是如果这边出错了
			// 要怎么办？
			// 保守做法：因为借助于 Redis 来做限流，那么 Redis 崩溃了，为了防止系统崩溃，直接限流
			ctx.AbortWithStatus(http.StatusInternalServerError) // 返回500内部服务器错误
			// 激进做法：虽然 Redis 崩溃了，但是这个时候还是要尽量服务正常的用户，所以不限流
			// ctx.Next()  // 继续处理请求
			return
		}
		if limited {
			log.Println(err)                                // 打印错误信息
			ctx.AbortWithStatus(http.StatusTooManyRequests) // 返回429请求过多
			return
		}
		ctx.Next() // 继续处理请求
	}
}
