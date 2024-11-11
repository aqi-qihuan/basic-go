package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 定义一个LoginMiddlewareBuilder结构体
type LoginMiddlewareBuilder struct {
}

// CheckLogin方法用于检查用户是否已经登录
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 返回一个gin.HandlerFunc类型的函数
	return func(ctx *gin.Context) {
		// 获取请求的URL路径
		path := ctx.Request.URL.Path
		// 如果路径是/users/signup或/users/login，则不需要登录校验
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		// 获取session
		sess := sessions.Default(ctx)
		// 如果session中没有userId，则表示用户未登录
		if sess.Get("userId") == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
