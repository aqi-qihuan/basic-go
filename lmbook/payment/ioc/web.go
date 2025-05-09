package ioc

import (
	"basic-go/lmbook/payment/web"
	"basic-go/lmbook/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func InitGinServer(hdl *web.WechatHandler) *ginx.Server {
	engine := gin.Default()
	hdl.RegisterRoutes(engine)
	addr := viper.GetString("http.addr")
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "daming_geektime",
		Subsystem: "lmbook_payment",
		Name:      "http",
	})
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}
