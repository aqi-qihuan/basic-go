//go:build wireinject

package main

import (
	"basic-go/lmbook/article/events"
	"basic-go/lmbook/article/grpc"
	"basic-go/lmbook/article/ioc"
	"basic-go/lmbook/article/repository"
	"basic-go/lmbook/article/repository/cache"
	"basic-go/lmbook/article/repository/dao"
	"basic-go/lmbook/article/service"
	"basic-go/lmbook/pkg/wego"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitUserRpcClient,
	ioc.InitProducer,
	ioc.InitEtcdClient,
	ioc.InitDB,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		events.NewSaramaSyncProducer,
		cache.NewRedisArticleCache,
		dao.NewGORMArticleDAO,
		repository.NewArticleRepository,
		repository.NewGrpcAuthorRepository,
		service.NewArticleService,
		grpc.NewArticleServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
