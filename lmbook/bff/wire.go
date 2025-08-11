//go:build wireinject

package main

import (
	"basic-go/lmbook/bff/ioc"
	"basic-go/lmbook/bff/web"
	"basic-go/lmbook/bff/web/jwt"
	"basic-go/lmbook/pkg/wego"
	"github.com/google/wire"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitEtcdClient,

		web.NewArticleHandler,
		web.NewUserHandler,
		web.NewRewardHandler,
		jwt.NewRedisHandler,

		ioc.InitUserClient,
		ioc.InitIntrClient,
		ioc.InitRewardClient,
		ioc.InitCodeClient,
		ioc.InitArticleClient,
		ioc.InitGinServer,
		wire.Struct(new(wego.App), "WebServer"),
	)
	return new(wego.App)
}
