package ioc

import (
	"basic-go/lanmengbook/internal/service/sms"
	"basic-go/lanmengbook/internal/service/sms/localsms"
	"basic-go/lanmengbook/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSMSService() sms.Service {
	// 注释：创建一个带有速率限制的短信服务，使用本地短信服务和Redis滑动窗口限流器
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter())
	// 注释：返回一个新的本地短信服务实例
	return localsms.NewService()
	// 如果有需要，就可以用这个
	//return initTencentSMSService()
}

// 初始化腾讯短信服务
func initTencentSMSService() sms.Service {
	// 从环境变量中获取腾讯 SMS 的 secret id
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	// 如果环境变量中没有找到 secret id，则抛出 panic
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	// 从环境变量中获取腾讯 SMS 的 secret key
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	// 如果环境变量中没有找到 secret key，则抛出 panic
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	// 创建腾讯 SMS 客户端
	c, err := tencentSMS.NewClient(
		// 使用 secret id 和 secret key 创建凭证
		common.NewCredential(secretId, secretKey),
		// 设置地域为南京
		"ap-nanjing",
		// 创建客户端配置
		profile.NewClientProfile(),
	)
	// 如果创建客户端时发生错误，则抛出 panic
	if err != nil {
		panic(err)
	}
	// 返回腾讯短信服务实例
	return tencent.NewService(c, "1400842696", "妙影科技")
}
