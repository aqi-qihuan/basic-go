package sms

import "context"

// Service 定义了一个短信服务接口
type Service interface {
	// Send 方法用于发送短信
	Send(ctx context.Context, tplId string,
		args []string, numbers ...string) error
}
