package ioc

import (
	"basic-go/lmbook/internal/web"
	ijwt "basic-go/lmbook/internal/web/jwt"
	"basic-go/lmbook/internal/web/middleware"
	"basic-go/lmbook/pkg/ginx"
	"basic-go/lmbook/pkg/ginx/middleware/prometheus"
	"basic-go/lmbook/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	otelgin "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler,
	artHdl *web.ArticleHandler,
	wechatHdl *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	wechatHdl.RegisterRoutes(server)
	artHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable,
	hdl ijwt.Handler, l logger.LoggerV1) []gin.HandlerFunc {
	pb := &prometheus.Builder{
		Namespace: "geektime_daming",
		Subsystem: "lmbook",
		Name:      "gin_http",
		Help:      "统计 GIN 的HTTP接口数据",
	}
	ginx.InitCounter(prometheus2.CounterOpts{
		Namespace: "geektime_daming",
		Subsystem: "lmbook",
		Name:      "biz_code",
		Help:      "统计业务错误码",
	})

	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,

			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
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
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		otelgin.Middleware("lmbook"),
		// ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		//middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
		//	l.Debug("", logger.Field{Key: "req", Val: al})
		//}).AllowReqBody().AllowRespBody().Build(),
		middleware.NewLoginJWTMiddlewareBuilder(hdl).CheckLogin(),
	}
}
