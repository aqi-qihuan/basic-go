//go:build wireinject

package main

import (
	"basic-go/lmbook/cronjob/grpc"
	"basic-go/lmbook/cronjob/ioc"
	"basic-go/lmbook/cronjob/repository"
	"basic-go/lmbook/cronjob/repository/dao"
	"basic-go/lmbook/cronjob/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewCronJobServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
