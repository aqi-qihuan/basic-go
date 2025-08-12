//go:build wireinject

package main

import (
	"basic-go/lmbook/pkg/wego"
	"basic-go/lmbook/sms/grpc"
	"basic-go/lmbook/sms/ioc"
	"github.com/google/wire"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitSmsTencentService,
		grpc.NewSmsServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
