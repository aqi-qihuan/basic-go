package ioc

import (
	wechat2 "basic-go/lmbook/oauth2/service/wechat"
	"basic-go/lmbook/pkg/logger"
	"github.com/spf13/viper"
)

func InitService(logv1 logger.LoggerV1) wechat2.Service {
	type Config struct {
		AppID     string `yaml:"appId"`
		AppSecret string `yaml:"appSecret"`
	}
	var cfg Config
	err := viper.UnmarshalKey("weChatConf", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat2.NewService(cfg.AppID, cfg.AppSecret, logv1)
}
