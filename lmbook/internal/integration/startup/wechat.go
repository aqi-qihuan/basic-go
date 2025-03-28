package startup

import (
	"basic-go/lmbook/internal/service/oauth2/wechat"
	"basic-go/lmbook/pkg/logger"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
