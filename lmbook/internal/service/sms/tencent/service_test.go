package tencent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
	"testing"
)

// 这个需要手动跑，也就是你需要在本地搞好这些环境变量
func TestSender(t *testing.T) {
	// 从环境变量中获取SMS_SECRET_ID
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	// 如果环境变量SMS_SECRET_ID不存在，测试失败
	if !ok {
		t.Fatal()
	}
	// 从环境变量中获取SMS_SECRET_KEY
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")

	// 创建一个新的SMS客户端
	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",               // 设置地域为南京
		profile.NewClientProfile()) // 创建一个新的客户端配置
	if err != nil {
		t.Fatal(err) // 如果创建客户端失败，测试失败
	}

	// 创建一个新的服务实例
	s := NewService(c, "1400842696", "妙影科技")

	// 定义测试用例
	testCases := []struct {
		name    string   // 测试用例名称
		tplId   string   // 模板ID
		params  []string // 模板参数
		numbers []string // 手机号码
		wantErr error    // 预期的错误
	}{
		{
			name:   "发送验证码",            // 测试用例名称
			tplId:  "1877556",          // 模板ID
			params: []string{"123456"}, // 模板参数
			// 改成你的手机号码
			numbers: []string{""}, // 手机号码
		},
	}
	// 遍历测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 调用Send方法发送短信
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			// 断言返回的错误与预期的错误是否一致
			assert.Equal(t, tc.wantErr, er)
		})
	}
}
