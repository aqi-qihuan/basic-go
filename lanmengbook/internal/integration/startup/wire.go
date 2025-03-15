//go:build wireinject

package startup

import (
	"basic-go/lanmengbook/internal/repository"
	"basic-go/lanmengbook/internal/repository/cache"
	"basic-go/lanmengbook/internal/repository/dao"
	"basic-go/lanmengbook/internal/service"
	"basic-go/lanmengbook/internal/web"
	"basic-go/lanmengbook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		InitRedis, ioc.InitDB,
		// DAO 部分
		dao.NewUserDAO,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
