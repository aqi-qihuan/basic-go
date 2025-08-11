//go:build wireinject

package main

import (
	"basic-go/lmbook/code/grpc"
	"basic-go/lmbook/code/ioc"
	"basic-go/lmbook/code/repository"
	"basic-go/lmbook/code/repository/cache"
	"basic-go/lmbook/code/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc.InitSmsRpcClient,
		cache.NewRedisCodeCache,
		repository.NewCachedCodeRepository,
		service.NewSMSCodeService,
		grpc.NewCodeServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
