//go:build wireinject

package startup

import (
	"basic-go/lmbook/internal/web"
	ijwt "basic-go/lmbook/internal/web/jwt"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB,
					InitLogger)
var userSvcProvider = wire.NewSet( // Service 部分
	dao.NewUserDAO,
	cache.NewCodeCache,
	repository.NewCachedUserRepository,
	service.NewUserService)
var articlSvcProvider = wire.NewSet(
	repository.NewCachedArticleRepository,
	cache.NewArticleRedisCache,
	dao.NewArticleGORMDAO,
	service.NewArticleService)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		articlSvcProvider,

		// cache 部分
		cache.NewCodeCache,

		// repository 部分
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		InitWechatService,

		// handler 部分
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		articlSvcProvider,
		repository.NewCachedArticleRepository,
		cache.NewArticleRedisCache,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
