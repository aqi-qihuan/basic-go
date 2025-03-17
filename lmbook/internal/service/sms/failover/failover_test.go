package failover

import (
	"basic-go/lmbook/internal/service/sms"
	smsmocks "basic-go/lmbook/internal/service/sms/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

// TestFailOverSMSService_Send 测试 FailOverSMSService 的 Send 方法
func TestFailOverSMSService_Send(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		// 测试用例名称，用于标识每个测试场景
		name string
		// 用于生成模拟对象的函数，返回一个 sms.Service 切片
		mock func(ctrl *gomock.Controller) []sms.Service
		// 预期的错误，若发送成功则为 nil
		wantErr error
	}{
		{
			name: "一次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				// 创建第一个模拟的 SMS 服务
				svc0 := smsmocks.NewMockService(ctrl)
				// 设置期望的 Send 方法调用，返回 nil 表示成功
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0}
			},
		},
		{
			name: "第二次发送成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				// 创建第一个模拟的 SMS 服务
				svc0 := smsmocks.NewMockService(ctrl)
				// 设置期望的 Send 方法调用，返回错误表示失败
				svc0.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				// 创建第二个模拟的 SMS 服务
				svc1 := smsmocks.NewMockService(ctrl)
				// 设置期望的 Send 方法调用，返回 nil 表示成功
				svc1.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				return []sms.Service{svc0, svc1}
			},
		},
		{
			name: "全部失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				// 创建第一个模拟的 SMS 服务
				svc0 := smsmocks.NewMockService(ctrl)
				// 设置期望的 Send 方法调用，返回错误表示失败
				svc0.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				// 创建第二个模拟的 SMS 服务
				svc1 := smsmocks.NewMockService(ctrl)
				// 设置期望的 Send 方法调用，返回错误表示失败
				svc1.EXPECT().Send(gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				return []sms.Service{svc0, svc1}
			},
			// 预期错误，当所有服务都发送失败时会返回此错误
			wantErr: errors.New("轮询了所有服务服务商，但是都发送失败了"),
		},
	}

	// 遍历测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个新的 mock 控制器，用于管理模拟对象
			ctrl := gomock.NewController(t)
			// 确保测试结束后释放资源，调用 Finish 方法来验证模拟对象的期望是否满足
			defer ctrl.Finish()
			// 创建 FailOverSMSService 实例，传入模拟的服务切片
			svc := NewFailOverSMSService(tc.mock(ctrl))
			// 调用 Send 方法，传入上下文、模板 ID、参数和电话号码
			err := svc.Send(context.Background(), "123", []string{"123"}, "12345")
			// 断言实际错误与预期错误相等，验证测试结果
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
