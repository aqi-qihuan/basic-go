//go:build wireinject

package main

import (
	"basic-go/lmbook/pkg/wego"
	"basic-go/lmbook/user/grpc"
	"basic-go/lmbook/user/ioc"
	"basic-go/lmbook/user/repository"
	"basic-go/lmbook/user/repository/cache"
	"basic-go/lmbook/user/repository/dao"
	"basic-go/lmbook/user/service"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisUserCache,
		dao.NewGORMUserDAO,
		repository.NewCachedUserRepository,
		service.NewUserService,
		grpc.NewUserServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
