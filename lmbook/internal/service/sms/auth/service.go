package auth

import (
	"basic-go/lmbook/internal/service/sms"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

// SMSService 是一个结构体，用于发送短信服务
type SMSService struct {
	// svc 是一个短信服务接口
	svc sms.Service
	// key 是用于解析 JWT 的密钥
	key []byte
}

// Send 方法用于发送短信
// ctx 是上下文，用于控制请求的取消和超时
// tplToken 是模板的 JWT 令牌
// args 是模板参数
// numbers 是接收短信的手机号码列表
func (s *SMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {

	// 定义一个 SMSClaims 结构体变量
	var claims SMSClaims
	// 使用 jwt.ParseWithClaims 解析 JWT 令牌
	// tplToken 是待解析的 JWT 令牌
	// &claims 是用于存储解析结果的变量
	// func(token *jwt.Token) (interface{}, error) 是一个回调函数，用于提供解析 JWT 所需的密钥
	// 该回调函数返回短信服务的密钥 s.key
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	// 如果解析过程中出现错误，返回该错误
	if err != nil {
		return err
	}
	// 调用短信服务接口的 Send 方法发送短信
	// ctx 是上下文
	// claims.TplID 是模板 ID
	// args 是模板参数
	// numbers 是接收短信的手机号码列表
	return s.svc.Send(ctx, claims.Tpl, args, numbers...)

}

// SMSClaims 是一个结构体，用于存储 JWT 令牌中的声明
type SMSClaims struct {
	// RegisteredClaims 是 JWT 标准声明
	jwt.RegisteredClaims
	// TplID 是模板 ID
	Tpl string
	// 额外加字段
}
