package middleware

import (
	ijwt "basic-go/lanmengbook/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// LoginJWTMiddlewareBuilder 是一个用于构建登录JWT中间件的构建器结构体
type LoginJWTMiddlewareBuilder struct {
	// Handler 是一个嵌入字段，表示实现了ijwt.Handler接口的类型
	ijwt.Handler
}

// NewLoginJWTMiddlewareBuilder 是一个工厂函数，用于创建一个新的LoginJWTMiddlewareBuilder实例
// 参数 hdl 是一个实现了ijwt.Handler接口的处理器
// 返回一个指向新创建的LoginJWTMiddlewareBuilder实例的指针
func NewLoginJWTMiddlewareBuilder(hdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	// 返回一个新创建的LoginJWTMiddlewareBuilder实例，并将传入的处理器hdl赋值给其Handler字段
	return &LoginJWTMiddlewareBuilder{
		Handler: hdl,
	}
}

// LoginJWTMiddlewareBuilder 是一个用于构建登录 JWT 中间件的构建器
func (m *LoginJWTMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取当前请求的路径
		path := ctx.Request.URL.Path
		// 不需要登录校验的路径
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/users/login_sms/code/send" ||
			path == "/users/login_sms" ||
			path == "/oauth2/wechat/authurl" ||
			path == "/oauth2/wechat/callback" {
			// 不需要登录校验
			return
		}
		tokenStr := m.ExtractToken(ctx)
		var uc ijwt.UserClaims

		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return ijwt.JWTKey, nil
		})
		if err != nil {
			// token 不对，token 是伪造的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			// 在这里发现 access_token 过期了，生成一个新的 access_token

			// token 解析出来了，但是 token 可能是非法的，或者过期了的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 这里看
		err = m.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// token 无效或者 redis 有问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 可以兼容 Redis 异常的情况
		// 做好监控，监控有没有 error
		//if cnt > 0 {
		//	// token 无效或者 redis 有问题
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		ctx.Set("user", uc)
	}
}
