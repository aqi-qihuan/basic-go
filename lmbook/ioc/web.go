package ioc

import (
	"basic-go/lmbook/internal/web"
	ijwt "basic-go/lmbook/internal/web/jwt"
	"basic-go/lmbook/internal/web/middleware"
	"basic-go/lmbook/pkg/ginx/middleware/ratelimit"
	"basic-go/lmbook/pkg/limiter"
	"basic-go/lmbook/pkg/logger"
	"context"
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

// InitWebServer 初始化Web服务器，并注册中间件和路由处理程序
// 参数 mdls: 中间件列表，类型为 []gin.HandlerFunc
// 参数 userHdl: 用户处理程序，类型为 *web.UserHandler
// 参数 wechatHdl: 微信OAuth2处理程序，类型为 *web.OAuth2WechatHandler
// 返回值: 初始化后的 gin.Engine 实例
func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	// 创建一个默认的gin引擎实例
	server := gin.Default()
	// 使用传入的中间件列表，将中间件应用到服务器上
	server.Use(mdls...)
	// 调用用户处理程序的 RegisterRoutes 方法，将用户相关的路由注册到服务器上
	userHdl.RegisterRoutes(server)
	// 调用微信OAuth2处理程序的 RegisterRoutes 方法，将微信OAuth2相关的路由注册到服务器上
	wechatHdl.RegisterRoutes(server)
	// 返回初始化并注册了路由和中间件的服务器实例
	return server
}

// InitGinMiddlewares 初始化Gin中间件
// 参数：
//   - redisClient: Redis命令客户端，用于实现限流等功能
//   - hdl: JWT处理器，用于处理JWT相关的逻辑
//   - l: 日志记录器，用于记录日志
//
// 返回值：
//   - []gin.HandlerFunc: 返回一个Gin中间件函数的切片
func InitGinMiddlewares(redisClient redis.Cmdable,
	hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 初始化CORS中间件
		cors.New(cors.Config{
			// 是否允许所有来源
			//AllowAllOrigins: true,
			// 允许的来源列表
			//AllowOrigins:     []string{"http://localhost:3000"},
			// 是否允许携带凭证（如Cookies）
			AllowCredentials: true,

			// 允许的请求头列表
			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			//     定义一个允许跨域请求的函数
			AllowOriginFunc: func(origin string) bool {
				// 检查请求的origin是否以"http://localhost"开头
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
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
			l.Debug("", logger.Field{Key: "req", Val: al})
		}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
	}
}
