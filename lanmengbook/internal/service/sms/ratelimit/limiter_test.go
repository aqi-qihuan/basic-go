package ratelimit

import (
	"basic-go/lanmengbook/internal/service/sms"
	smsmocks "basic-go/lanmengbook/internal/service/sms/mocks"
	"basic-go/lanmengbook/pkg/limiter"
	limitermocks "basic-go/lanmengbook/pkg/limiter/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestRateLimitSMSService_Send 测试RateLimitSMSService的Send方法
func TestRateLimitSMSService_Send(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter)

		// 一个输入都没有

		// 预期输出
		wantErr error
	}{
		{
			name: "不限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				// 创建sms.Service的mock对象
				svc := smsmocks.NewMockService(ctrl)
				// 创建limiter.Limiter的mock对象
				l := limitermocks.NewMockLimiter(ctrl)
				// 设置limiter.Limiter的Limit方法的预期行为，返回false表示不限流
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				// 设置sms.Service的Send方法的预期行为，返回nil表示发送成功
				svc.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil)
				return svc, l
			},
		},
		{
			name: "限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				// 创建sms.Service的mock对象
				svc := smsmocks.NewMockService(ctrl)
				// 创建limiter.Limiter的mock对象
				l := limitermocks.NewMockLimiter(ctrl)
				// 设置limiter.Limiter的Limit方法的预期行为，返回true表示限流
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return svc, l
			},
			wantErr: errLimited, // 预期错误为限流错误
		},
		{
			name: "限流器错误",
			mock: func(ctrl *gomock.Controller) (sms.Service, limiter.Limiter) {
				// 创建sms.Service的mock对象
				svc := smsmocks.NewMockService(ctrl)
				// 创建limiter.Limiter的mock对象
				l := limitermocks.NewMockLimiter(ctrl)
				// 设置limiter.Limiter的Limit方法的预期行为，返回false和错误表示限流器错误
				l.EXPECT().Limit(gomock.Any(), gomock.Any()).
					Return(false, errors.New("redis限流器错误"))
				return svc, l
			},
			wantErr: errors.New("redis限流器错误"), // 预期错误为限流器错误
		},
	}

	// 遍历测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建gomock.Controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish() // 确保测试结束后释放资源
			// 获取mock对象
			smsSvc, l := tc.mock(ctrl)
			// 创建RateLimitSMSService实例
			svc := NewRateLimitSMSService(smsSvc, l)
			// 调用Send方法
			err := svc.Send(context.Background(), "abc",
				[]string{"123"}, "123456")
			// 断言错误是否符合预期
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
