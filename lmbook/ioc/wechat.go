package ioc

import (
	"basic-go/lmbook/internal/service/oauth2/wechat"
	"basic-go/lmbook/pkg/logger"
	"os"
)

// InitWechatService 初始化微信服务
// 参数 l 是一个日志记录器接口，用于记录日志
// 返回一个微信服务实例
func InitWechatService(l logger.LoggerV1) wechat.Service {
	// 从环境变量中获取 WECHAT_APP_ID
	appID, ok := os.LookupEnv("WECHAT_APP_ID")
	// 检查 WECHAT_APP_ID 是否存在
	if !ok {
		// 如果不存在，抛出 panic，终止程序
		panic("找不到环境变量 WECHAT_APP_ID")
	}
	// 从环境变量中获取 WECHAT_APP_SECRET
	appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	// 检查 WECHAT_APP_SECRET 是否存在
	if !ok {
		// 如果不存在，抛出 panic，终止程序
		panic("找不到环境变量 WECHAT_APP_SECRET")
	}
	// 使用获取到的 appID、appSecret 和日志记录器 l 创建一个新的微信服务实例
	return wechat.NewService(appID, appSecret, l)
}
