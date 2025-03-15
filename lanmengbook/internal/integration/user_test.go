package integration

import (
	"basic-go/lanmengbook/internal/integration/startup"
	"basic-go/lanmengbook/internal/web"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// TestUserHandler_SendSMSCode 测试发送短信验证码的功能
func TestUserHandler_SendSMSCode(t *testing.T) {
	// 初始化Redis客户端
	rdb := startup.InitRedis()
	// 初始化Web服务器
	server := startup.InitWebServer()
	// 定义测试用例
	testCases := []struct {
		name string
		// before 测试用例执行前的准备工作
		before func(t *testing.T)
		// after 测试用例执行后的清理工作
		after func(t *testing.T)

		// phone 测试用例中的手机号码
		phone string

		// wantCode 期望的HTTP响应状态码
		wantCode int
		// wantBody 期望的HTTP响应体
		wantBody web.Result
	}{
		{
			// 测试用例的名称
			name: "发送成功的用例",
			// 测试用例执行前的钩子函数
			before: func(t *testing.T) {

				// 这里可以放置一些测试前的准备工作，当前为空
			},
			// 测试用例执行后的钩子函数
			after: func(t *testing.T) {
				// 创建一个带有超时的上下文，超时时间为10秒
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				// 在函数结束时取消上下文，释放资源
				defer cancel()
				// 定义Redis键名，用于存储验证码
				key := "phone_code:login:15212345678"
				// 从Redis中获取存储的验证码
				code, err := rdb.Get(ctx, key).Result()
				// 断言没有错误
				assert.NoError(t, err)
				// 断言获取到的验证码长度大于0
				assert.True(t, len(code) > 0)
				// 获取验证码的剩余生存时间
				dur, err := rdb.TTL(ctx, key).Result()
				// 断言没有错误
				assert.NoError(t, err)
				// 断言验证码的剩余生存时间大于9分钟50秒
				assert.True(t, dur > time.Minute*9+time.Second+50)
				// 删除Redis中的验证码
				err = rdb.Del(ctx, key).Err()
				// 断言没有错误
				assert.NoError(t, err)
			},
			// 测试用例中的电话号码
			phone: "15212345678",
			// 期望的HTTP响应状态码
			wantCode: http.StatusOK,
			// 期望的HTTP响应体
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
		{
			// 测试用例的名称，描述当前测试的场景
			name: "未输入手机号码",
			// 测试用例执行前的钩子函数，用于准备测试环境
			before: func(t *testing.T) {

				// 当前为空，表示在测试开始前不需要进行任何操作
			},
			after:    func(t *testing.T) {},
			wantCode: http.StatusOK,
			wantBody: web.Result{
				Code: 4,
				Msg:  "请输入手机号码",
			},
		},
		{
			// 测试用例的名称
			name: "发送太频繁",
			// 测试用例的前置处理函数
			before: func(t *testing.T) {
				// 创建一个带有10秒超时的上下文
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				// 在函数结束时取消上下文，释放资源
				defer cancel()
				// 定义存储在Redis中的键
				key := "phone_code:login:15212345678"
				// 将键值对存储到Redis中，设置过期时间为9分钟50秒
				err := rdb.Set(ctx, key, "123456", time.Minute*9+time.Second*50).Err()
				// 断言操作没有错误
				assert.NoError(t, err)
			},
			// 测试用例的后置处理函数
			after: func(t *testing.T) {
				// 创建一个带有10秒超时的上下文
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				// 在函数结束时取消上下文，释放资源
				defer cancel()
				// 定义要获取并删除的键
				key := "phone_code:login:15212345678"
				// 从Redis中获取并删除该键的值
				code, err := rdb.GetDel(ctx, key).Result()
				// 断言操作没有错误
				assert.NoError(t, err)
				// 断言获取的值与预期值一致
				assert.Equal(t, "123456", code)
			},
			// 测试用例中的手机号码
			phone: "15212345678",
			// 预期的HTTP响应状态码
			wantCode: http.StatusOK,
			// 预期的HTTP响应体内容
			wantBody: web.Result{
				Code: 4,
				Msg:  "短信发送太频繁，请稍后再试",
			},
		},
		{
			// 测试用例名称
			name: "系统错误",
			// 测试用例的前置处理函数
			before: func(t *testing.T) {
				// 创建一个带有10秒超时的上下文
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				// 在函数结束时取消上下文，释放资源
				defer cancel()
				// 定义一个键，用于存储验证码
				key := "phone_code:login:15212345678"
				// 将验证码"123456"存储到Redis中，设置过期时间为0（永不过期）
				err := rdb.Set(ctx, key, "123456", 0).Err()
				// 断言存储操作没有错误
				assert.NoError(t, err)
			},
			// 测试用例的后置处理函数
			after: func(t *testing.T) {
				// 创建一个带有10秒超时的上下文
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
				// 在函数结束时取消上下文，释放资源
				defer cancel()
				// 定义一个键，用于获取验证码
				key := "phone_code:login:15212345678"
				// 从Redis中获取并删除该键对应的值
				code, err := rdb.GetDel(ctx, key).Result()
				// 断言获取操作没有错误
				assert.NoError(t, err)
				// 断言获取的验证码值为"123456"
				assert.Equal(t, "123456", code)
			},
			// 测试用例中的手机号码
			phone: "15212345678",
			// 期望的HTTP响应状态码
			wantCode: http.StatusOK,
			// 期望的HTTP响应体内容
			wantBody: web.Result{
				// 期望的响应码
				Code: 5,
				// 期望的响应消息
				Msg: "系统错误",
			},
		},
	}

	// 遍历测试用例
	for _, tc := range testCases {
		// 使用测试用例的名称作为子测试的名称
		t.Run(tc.name, func(t *testing.T) {
			// 执行测试用例的准备工作
			tc.before(t)
			// 在函数结束时执行测试用例的清理工作
			defer tc.after(t)

			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)
			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res web.Result
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}

}
