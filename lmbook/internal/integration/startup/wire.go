//go:build wireinject

package startup

import (
	repository2 "basic-go/lmbook/interactive/repository"
	cache2 "basic-go/lmbook/interactive/repository/cache"
	dao2 "basic-go/lmbook/interactive/repository/dao"
	service2 "basic-go/lmbook/interactive/service"
	"basic-go/lmbook/internal/events/article"
	"basic-go/lmbook/internal/job"
	"basic-go/lmbook/internal/repository"
	"basic-go/lmbook/internal/repository/cache"
	"basic-go/lmbook/internal/repository/dao"
	"basic-go/lmbook/internal/service"
	"basic-go/lmbook/internal/service/sms"
	"basic-go/lmbook/internal/service/sms/async"
	"basic-go/lmbook/internal/web"
	ijwt "basic-go/lmbook/internal/web/jwt"
	"basic-go/lmbook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB,
	InitSaramaClient,
	InitSyncProducer,
	InitLogger)

var jobProviderSet = wire.NewSet(
	service.NewCronJobService,
	repository.NewPreemptJobRepository,
	dao.NewGORMJobDAO)

var userSvcProvider = wire.NewSet(
	dao.NewUserDAO,
	cache.NewUserCache,
	repository.NewCachedUserRepository,
	service.NewUserService)

var articlSvcProvider = wire.NewSet(
	repository.NewCachedArticleRepository,
	cache.NewArticleRedisCache,
	dao.NewArticleGORMDAO,
	service.NewArticleService)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
	ioc.InitIntrClient,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		articlSvcProvider,
		interactiveSvcSet,
		// cache 部分
		cache.NewCodeCache,

		// repository 部分
		repository.NewCodeRepository,

		article.NewSaramaSyncProducer,

		// Service 部分
		ioc.InitSMSService,
		service.NewCodeService,
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

func InitAsyncSmsService(svc sms.Service) *async.Service {
	wire.Build(thirdPartySet, repository.NewAsyncSMSRepository,
		dao.NewGORMAsyncSmsDAO,
		async.NewService,
	)
	return &async.Service{}
}

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		userSvcProvider,
		interactiveSvcSet,
		repository.NewCachedArticleRepository,
		cache.NewArticleRedisCache,
		service.NewArticleService,
		article.NewSaramaSyncProducer,
		web.NewArticleHandler)
	return &web.ArticleHandler{}
}

func InitJobScheduler() *job.Scheduler {
	wire.Build(jobProviderSet, thirdPartySet, job.NewScheduler)
	return &job.Scheduler{}
}
