package main

import (
	"basic-go/lmbook/pkg/wego"
	"basic-go/lmbook/tag/grpc"
	"basic-go/lmbook/tag/ioc"
	"basic-go/lmbook/tag/repository/cache"
	"basic-go/lmbook/tag/repository/dao"
	"basic-go/lmbook/tag/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitDB,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisTagCache,
		dao.NewGORMTagDAO,
		ioc.InitRepository,
		service.NewTagService,
		grpc.NewTagServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
