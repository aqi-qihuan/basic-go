//go:build wireinject

package main

import (
	"basic-go/lmbook/pkg/wego"
	"basic-go/lmbook/reward/grpc"
	"basic-go/lmbook/reward/ioc"
	"basic-go/lmbook/reward/repository"
	"basic-go/lmbook/reward/repository/cache"
	"basic-go/lmbook/reward/repository/dao"
	"basic-go/lmbook/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitEtcdClient,
	ioc.InitRedis)

func Init() *wego.App {
	wire.Build(thirdPartySet,
		service.NewWechatNativeRewardService,
		ioc.InitAccountClient,
		ioc.InitGRPCxServer,
		ioc.InitPaymentClient,
		repository.NewRewardRepository,
		cache.NewRewardRedisCache,
		dao.NewRewardGORMDAO,
		grpc.NewRewardServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
