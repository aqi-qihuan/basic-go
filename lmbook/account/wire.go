//go:build wireinject

package main

import (
	"basic-go/lmbook/account/grpc"
	"basic-go/lmbook/account/ioc"
	"basic-go/lmbook/account/repository"
	"basic-go/lmbook/account/repository/dao"
	"basic-go/lmbook/account/service"
	"basic-go/lmbook/pkg/wego"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitGRPCxServer,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
