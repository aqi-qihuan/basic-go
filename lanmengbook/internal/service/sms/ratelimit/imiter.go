package ratelimit

import (
	"basic-go/lanmengbook/internal/service/sms"
	"basic-go/lanmengbook/pkg/limiter"
	"context"
	"errors"
)

// 定义一个错误变量，表示触发限流
var errLimited = errors.New("触发限流")

// 确保RateLimitSMSService实现了sms.Service接口
var _ sms.Service = &RateLimitSMSService{}

// 定义一个结构体RateLimitSMSService，用于装饰sms.Service，添加限流功能
type RateLimitSMSService struct {
	// 被装饰的sms.Service实例
	svc sms.Service
	// 限流器实例
	limiter limiter.Limiter
	// 限流器的键
	key string
}

// 定义一个结构体RateLimitSMSServiceV1，使用嵌入字段简化代码
type RateLimitSMSServiceV1 struct {
	// 嵌入sms.Service接口，直接使用其方法
	sms.Service
	// 限流器实例
	limiter limiter.Limiter
	// 限流器的键
	key string
}

// Send方法实现了sms.Service接口的Send方法，添加了限流逻辑
func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 调用限流器的Limit方法，检查是否触发限流
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		// 如果Limit方法返回错误，直接返回该错误
		return err
	}
	if limited {
		// 如果触发限流，返回自定义的限流错误
		return errLimited
	}
	// 如果未触发限流，调用被装饰的sms.Service的Send方法
	return r.svc.Send(ctx, tplId, args, numbers...)
}

// NewRateLimitSMSService函数用于创建RateLimitSMSService实例
func NewRateLimitSMSService(svc sms.Service,
	l limiter.Limiter) *RateLimitSMSService {
	// 返回RateLimitSMSService实例，传入被装饰的sms.Service实例和限流器实例，并设置限流器的键
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "sms-limiter",
	}
}
