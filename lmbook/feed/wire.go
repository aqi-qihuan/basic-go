//go:build wireinject

package main

import (
	"basic-go/lmbook/feed/events"
	"basic-go/lmbook/feed/grpc"
	"basic-go/lmbook/feed/ioc"
	"basic-go/lmbook/feed/repository"
	"basic-go/lmbook/feed/repository/cache"
	"basic-go/lmbook/feed/repository/dao"
	"basic-go/lmbook/feed/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewFeedPushEventDAO,
	dao.NewFeedPullEventDAO,
	cache.NewFeedEventCache,
	repository.NewFeedEventRepo,
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitDB,
	ioc.InitFollowClient,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.RegisterHandler,
		service.NewFeedService,
		grpc.NewFeedEventGrpcSvc,
		events.NewArticleEventConsumer,
		events.NewFeedEventConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
