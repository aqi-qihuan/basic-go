package ioc

import (
	"basic-go/lanmengbook/internal/web"
	"basic-go/lanmengbook/internal/web/middleware"
	"basic-go/lanmengbook/pkg/ginx/middleware/ratelimit"
	"basic-go/lanmengbook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

//func InitWebServerV1(mdls []gin.HandlerFunc, hdls []web.Handler) *gin.Engine {
//	server := gin.Default()
//	server.Use(mdls...)
//	for _, hdl := range hdls {
//		hdl.RegisterRoutes(server)
//	}
//	//userHdl.RegisterRoutes(server)
//	return server
//}

// InitWebServer 初始化Web服务器，并注册中间件和用户处理程序的路由
func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	// 创建一个默认的Gin引擎
	server := gin.Default()
	// 使用传入的中间件
	server.Use(mdls...)
	// 注册用户处理程序的路由
	userHdl.RegisterRoutes(server)
	// 注册微信OAuth2处理程序的路由
	oauth2WechatHdl.RegisterRoutes(server)
	// 返回配置好的Gin引擎
	return server
}

// InitGinMiddlewares 初始化Gin中间件，包括CORS、自定义中间件、限流器和JWT登录检查
func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 初始化CORS中间件
		cors.New(cors.Config{
			// 允许所有来源访问（注释掉以启用特定来源）
			//AllowAllOrigins: true,
			// 允许特定来源访问（注释掉以启用所有来源）
			//AllowOrigins:     []string{"http://localhost:3000"},
			// 允许携带凭证（如Cookies）
			AllowCredentials: true,

			// 允许的请求头
			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		func(ctx *gin.Context) {
			println("这是我的 Middleware")
		},
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}
