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

// TestTimeoutFailoverSMSService_Send 测试TimeoutFailoverSMSService的Send方法
func TestTimeoutFailoverSMSService_Send(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name      string                                      // 测试用例名称
		mock      func(ctrl *gomock.Controller) []sms.Service // 模拟的SMS服务
		threshold int32                                       // 超时阈值
		idx       int32                                       // 当前服务索引
		cnt       int32                                       // 超时计数

		wantErr error // 期望的错误
		wantCnt int32 // 期望的超时计数
		wantIdx int32 // 期望的服务索引
	}{
		{
			name: "咩有触发切换", // 测试用例名称
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil) // 设置期望的Send方法返回nil
				return []sms.Service{svc0} // 返回包含一个服务的切片
			},
			idx:       0,  // 当前服务索引
			cnt:       12, // 超时计数
			threshold: 15, // 超时阈值
			wantIdx:   0,  // 期望的服务索引
			// 成功了，重置超时计数
			wantCnt: 0,   // 期望的超时计数
			wantErr: nil, // 期望的错误
		},
		{

			name: "触发切换，成功", // 测试用例名称
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc1 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil) // 设置期望的Send方法返回nil
				return []sms.Service{svc0, svc1} // 返回包含两个服务的切片
			},
			idx:       0,  // 当前服务索引
			cnt:       15, // 超时计数
			threshold: 15, // 超时阈值
			// 触发了切换
			wantIdx: 1,   // 期望的服务索引
			wantCnt: 0,   // 期望的超时计数
			wantErr: nil, // 期望的错误
		},
		{
			name: "触发切换,失败", // 测试用例名称
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc1 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(errors.New("发送失败")) // 设置期望的Send方法返回错误
				return []sms.Service{svc0, svc1} // 返回包含两个服务的切片
			},
			idx:       1,                  // 当前服务索引
			cnt:       15,                 // 超时计数
			threshold: 15,                 // 超时阈值
			wantIdx:   0,                  // 期望的服务索引
			wantCnt:   0,                  // 期望的超时计数
			wantErr:   errors.New("发送失败"), // 期望的错误
		},
		{
			name: "触发切换,超时", // 测试用例名称
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc1 := smsmocks.NewMockService(ctrl) // 创建模拟的SMS服务
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded) // 设置期望的Send方法返回超时错误
				return []sms.Service{svc0, svc1} // 返回包含两个服务的切片
			},
			idx:       1,  // 当前服务索引
			cnt:       15, // 超时计数
			threshold: 15, // 超时阈值
			// 触发了切换
			wantIdx: 0,                        // 期望的服务索引
			wantCnt: 1,                        // 期望的超时计数
			wantErr: context.DeadlineExceeded, // 期望的错误
		},
	}
	// 遍历测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)                                            // 创建gomock控制器
			defer ctrl.Finish()                                                        // 确保测试结束后释放资源
			svc := NewTimeoutFailoverSMSService(tc.mock(ctrl), tc.threshold)           // 创建TimeoutFailoverSMSService实例
			svc.cnt = tc.cnt                                                           // 设置超时计数
			svc.idx = tc.idx                                                           // 设置当前服务索引
			err := svc.Send(context.Background(), "12", []string{"34"}, "15312345678") // 调用Send方法
			assert.Equal(t, tc.wantErr, err)                                           // 断言错误
			assert.Equal(t, tc.wantCnt, svc.cnt)                                       // 断言超时计数
			assert.Equal(t, tc.wantIdx, svc.idx)                                       // 断言服务索引
		})
	}
}
