package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
}

// Send 方法用于发送短信
func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 创建一个新的 SendSmsRequest 请求
	request := sms.NewSendSmsRequest()
	// 设置请求的上下文
	request.SetContext(ctx)
	// 设置短信应用 ID
	request.SmsSdkAppId = s.appId
	// 设置短信签名
	request.SignName = s.signName
	// 设置模板 ID
	request.TemplateId = ekit.ToPtr[string](tplId)
	// 设置模板参数
	request.TemplateParamSet = s.toPtrSlice(args)
	// 设置接收短信的电话号码
	request.PhoneNumberSet = s.toPtrSlice(numbers)
	// 发送短信并获取响应
	response, err := s.client.SendSms(request)
	// 异常处理
	if err != nil {
		return err
	}
	// 遍历响应中的发送状态
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			// 不可能进来这里
			continue
		}
		// 获取发送状态
		status := *statusPtr
		// 如果发送状态码不为 "Ok"，则表示发送失败
		if status.Code == nil || *(status.Code) != "Ok" {
			// 发送失败
			return fmt.Errorf("发送短信失败, code: %s, msg:%s", *status.Code,
				*status.Message)
		}

	}
	// 发送成功
	return nil
}

// toPtrSlice 方法将字符串切片转换为字符串指针切片
func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data,
		func(idx int, src string) *string {
			return &src
		})
}

// NewService 方法创建一个新的 Service 实例
func NewService(client *sms.Client, appId, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}

}
