//go:build wireinject

package main

import (
	"basic-go/lmbook/payment/grpc"
	"basic-go/lmbook/payment/ioc"
	"basic-go/lmbook/payment/repository"
	"basic-go/lmbook/payment/repository/dao"
	"basic-go/lmbook/payment/web"
	"basic-go/lmbook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitEtcdClient,
		ioc.InitKafka,
		ioc.InitProducer,
		ioc.InitWechatClient,
		dao.NewPaymentGORMDAO,
		ioc.InitDB,
		repository.NewPaymentRepository,
		grpc.NewWechatServiceServer,
		ioc.InitWechatNativeService,
		ioc.InitWechatConfig,
		ioc.InitWechatNotifyHandler,
		ioc.InitGRPCServer,
		web.NewWechatHandler,
		ioc.InitGinServer,
		ioc.InitLogger,
		wire.Struct(new(wego.App), "WebServer", "GRPCServer"))
	return new(wego.App)
}
