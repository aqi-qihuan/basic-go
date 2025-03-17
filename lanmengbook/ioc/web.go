package ioc

import (
	"basic-go/lanmengbook/internal/web"
	ijwt "basic-go/lanmengbook/internal/web/jwt"
	"basic-go/lanmengbook/internal/web/middleware"
	"basic-go/lanmengbook/pkg/ginx/middleware/ratelimit"
	"basic-go/lanmengbook/pkg/limiter"
	"basic-go/lanmengbook/pkg/logger"
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

// InitWebServer 初始化Web服务器，并注册中间件和用户处理程序的路由
func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler, wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	// 创建一个默认的Gin引擎
	server := gin.Default()
	// 使用传入的中间件
	server.Use(mdls...)
	// 注册用户处理程序的路由
	userHdl.RegisterRoutes(server)
	// 调用 wechatHdl 的 RegisterRoutes 方法，将路由注册到 server 中
	wechatHdl.RegisterRoutes(server)
	// 返回配置好的Gin引擎
	return server
}

// InitGinMiddlewares 初始化Gin中间件，包括CORS、自定义中间件、限流器和JWT登录检查
func InitGinMiddlewares(redisClient redis.Cmdable, hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
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
			// 定义一个允许跨域请求的函数，根据请求的来源判断是否允许
			AllowOriginFunc: func(origin string) bool {
				// 如果请求的来源以 "http://localhost" 开头，则允许跨域请求
				if strings.HasPrefix(origin, "http://localhost") {
					// 注释掉的代码：如果请求的来源包含 "localhost"，则允许跨域请求
					//if strings.Contains(origin, "localhost") {
					return true
				}
				// 如果请求的来源包含 "your_company.com"，则允许跨域请求
				return strings.Contains(origin, "your_company.com")
			},
			// 设置跨域请求的最大缓存时间，这里设置为12小时
			MaxAge: 12 * time.Hour,
		}),
		// 定义一个中间件函数，用于在请求处理之前打印日志
		func(ctx *gin.Context) {
			// 打印日志信息，表示这是自定义的中间件
			println("这是我的 Middleware")
		},
		// 使用Redis滑动窗口限流器构建一个限流中间件
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		// 构建一个日志中间件，用于记录请求和响应的详细信息
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
			// 使用Debug级别记录请求日志
			l.Debug("", logger.Field{Key: "req", Val: al})
		}).AllowReqBody().AllowRespBody().Build(),
		// 构建一个JWT登录验证中间件，用于检查用户是否已登录
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
	}
}
