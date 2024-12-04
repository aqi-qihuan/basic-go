package localsms

import (
	"context"
	"log"
)

// SmsService 定义了一个短信服务
type SmsService struct {
}

// NewService 创建一个新的 SmsService 实例
func NewService() *SmsService {
	return &SmsService{}
}

// Send 方法用于发送短信
func (s *SmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 打印日志，显示验证码
	log.Println("验证码是", args)
	// 返回 nil，表示发送成功
	return nil
}
