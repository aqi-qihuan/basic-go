package jwt

import "github.com/gin-gonic/gin"

// 定义一个名为 Handler 的接口，用于处理与令牌（Token）相关的操作
type Handler interface {
	// ClearToken 方法用于清除上下文中的令牌
	// 参数 ctx 是一个指向 gin.Context 的指针，表示当前的 HTTP 请求上下文
	// 返回值是一个 error，用于表示操作过程中是否发生错误
	ClearToken(ctx *gin.Context) error
	// ExtractToken 方法用于从上下文中提取令牌
	// 参数 ctx 是一个指向 gin.Context 的指针，表示当前的 HTTP 请求上下文
	// 返回值是一个字符串，表示提取到的令牌
	ExtractToken(ctx *gin.Context) string
	// SetLoginToken 方法用于设置登录令牌
	// 参数 ctx 是一个指向 gin.Context 的指针，表示当前的 HTTP 请求上下文
	// 参数 uid 是一个 int64 类型的用户 ID
	// 返回值是一个 error，用于表示操作过程中是否发生错误
	SetLoginToken(ctx *gin.Context, uid int64) error
	// SetJWTToken 方法用于设置 JWT 令牌
	// 参数 ctx 是一个指向 gin.Context 的指针，表示当前的 HTTP 请求上下文
	// 参数 uid 是一个 int64 类型的用户 ID
	// 参数 ssid 是一个字符串，表示会话 ID
	// 返回值是一个 error，用于表示操作过程中是否发生错误
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	// CheckSession 方法用于检查会话的有效性
	// 参数 ctx 是一个指向 gin.Context 的指针，表示当前的 HTTP 请求上下文
	// 参数 ssid 是一个字符串，表示会话 ID
	// 返回值是一个 error，用于表示操作过程中是否发生错误
	CheckSession(ctx *gin.Context, ssid string) error
}
