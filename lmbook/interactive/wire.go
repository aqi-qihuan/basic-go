//go:build wireinject

package main

import (
	"basic-go/lmbook/interactive/events"
	"basic-go/lmbook/interactive/grpc"
	"basic-go/lmbook/interactive/ioc"
	repository2 "basic-go/lmbook/interactive/repository"
	cache2 "basic-go/lmbook/interactive/repository/cache"
	dao2 "basic-go/lmbook/interactive/repository/dao"
	service2 "basic-go/lmbook/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitSrcDB,
	ioc.InitDstDB,
	ioc.InitDoubleWritePool,
	ioc.InitBizDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitSaramaSyncProducer,
	ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitInteractiveProducer,
		ioc.InitFixerConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		ioc.InitGinxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
