// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"basic-go/lmbook/internal/events/article"
	"basic-go/lmbook/internal/repository"
	"basic-go/lmbook/internal/repository/cache"
	"basic-go/lmbook/internal/repository/dao"
	"basic-go/lmbook/internal/service"
	"basic-go/lmbook/internal/web"
	"basic-go/lmbook/internal/web/jwt"
	"basic-go/lmbook/ioc"
	"github.com/google/wire"
)

import (
	_ "github.com/spf13/viper/remote"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	loggerV1 := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, loggerV1)
	db := ioc.InitDB(loggerV1)
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
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, loggerV1, interactiveCache)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(loggerV1, articleService, interactiveService)
	wechatService := ioc.InitWechatService(loggerV1)
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, handler, userService)
	engine := ioc.InitWebServer(v, userHandler, articleHandler, oAuth2WechatHandler)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(interactiveRepository, client, loggerV1)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	rankingService := service.NewBatchRankingService(interactiveService, articleService)
	rlockClient := ioc.InitRlockClient(cmdable)
	rankingJob := ioc.InitRankingJob(rankingService, rlockClient, loggerV1)
	cron := ioc.InitJobs(loggerV1, rankingJob)
	app := &App{
		server:    engine,
		consumers: v2,
		cron:      cron,
	}
	return app
}

// wire.go:

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)

var rankingSvcSet = wire.NewSet(cache.NewRankingRedisCache, repository.NewCachedRankingRepository, service.NewBatchRankingService)
