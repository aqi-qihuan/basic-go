// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
	"basic-go/lmbook/internal/web/jwt"
	"basic-go/lmbook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	loggerV1 := InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, loggerV1)
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, handler, codeService)
	articleDAO := dao.NewArticleGORMDAO(db)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(articleDAO, userRepository, articleCache)
	client := InitSaramaClient()
	syncProducer := InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDAO := dao2.NewGORMInteractiveDAO(db)
	interactiveCache := cache2.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository2.NewCachedInteractiveRepository(interactiveDAO, loggerV1, interactiveCache)
	interactiveService := service2.NewInteractiveService(interactiveRepository)
	interactiveServiceClient := ioc.InitIntrClient(interactiveService)
	articleHandler := web.NewArticleHandler(loggerV1, articleService, interactiveServiceClient)
	wechatService := InitWechatService(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService)
	engine := ioc.InitWebServer(v, userHandler, articleHandler, oAuth2WechatHandler)
	return engine
}

func InitAsyncSmsService(svc sms.Service) *async.Service {
	db := InitDB()
	asyncSmsDAO := dao.NewGORMAsyncSmsDAO(db)
	asyncSmsRepository := repository.NewAsyncSMSRepository(asyncSmsDAO)
	loggerV1 := InitLogger()
	asyncService := async.NewService(svc, asyncSmsRepository, loggerV1)
	return asyncService
}

func InitArticleHandler(dao3 dao.ArticleDAO) *web.ArticleHandler {
	loggerV1 := InitLogger()
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	cmdable := InitRedis()
	userCache := cache.NewUserCache(cmdable)
	userRepository := repository.NewCachedUserRepository(userDAO, userCache)
	articleCache := cache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewCachedArticleRepository(dao3, userRepository, articleCache)
	client := InitSaramaClient()
	syncProducer := InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDAO := dao2.NewGORMInteractiveDAO(db)
	interactiveCache := cache2.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository2.NewCachedInteractiveRepository(interactiveDAO, loggerV1, interactiveCache)
	interactiveService := service2.NewInteractiveService(interactiveRepository)
	interactiveServiceClient := ioc.InitIntrClient(interactiveService)
	articleHandler := web.NewArticleHandler(loggerV1, articleService, interactiveServiceClient)
	return articleHandler
}

func InitJobScheduler() *job.Scheduler {
	db := InitDB()
	jobDAO := dao.NewGORMJobDAO(db)
	cronJobRepository := repository.NewPreemptJobRepository(jobDAO)
	loggerV1 := InitLogger()
	cronJobService := service.NewCronJobService(cronJobRepository, loggerV1)
	scheduler := job.NewScheduler(cronJobService, loggerV1)
	return scheduler
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitRedis, InitDB,
	InitSaramaClient,
	InitSyncProducer,
	InitLogger)

var jobProviderSet = wire.NewSet(service.NewCronJobService, repository.NewPreemptJobRepository, dao.NewGORMJobDAO)

var userSvcProvider = wire.NewSet(dao.NewUserDAO, cache.NewUserCache, repository.NewCachedUserRepository, service.NewUserService)

var articlSvcProvider = wire.NewSet(repository.NewCachedArticleRepository, cache.NewArticleRedisCache, dao.NewArticleGORMDAO, service.NewArticleService)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO, cache2.NewInteractiveRedisCache, repository2.NewCachedInteractiveRepository, service2.NewInteractiveService, ioc.InitIntrClient)
