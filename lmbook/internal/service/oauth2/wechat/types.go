package wechat

import (
	"basic-go/lmbook/internal/domain"
	"basic-go/lmbook/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// 定义一个Service接口，包含两个方法：AuthURL和VerifyCode
type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

// 定义一个全局变量redirectURL，用于存储回调URL，并进行URL编码
var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

// 定义一个service结构体，包含appID、appSecret、http客户端和日志记录器
type service struct {
	appID     string
	appSecret string
	client    *http.Client
	l         logger.LoggerV1
}

// NewService函数用于创建一个新的service实例
func NewService(appID string, appSecret string, l logger.LoggerV1) Service {
	return &service{
		appID:     appID,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

// VerifyCode方法用于验证微信授权码，并获取用户信息
func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	// 构建获取access token的URL
	accessTokenUrl := fmt.Sprintf(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code`,
		s.appID, s.appSecret, code)
	// 创建一个新的HTTP GET请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accessTokenUrl, nil)
	if err != nil {
		return domain.WechatInfo{}, err // 如果请求创建失败，返回错误
	}
	// 发送HTTP请求
	httpResp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err // 如果请求发送失败，返回错误
	}

	// 定义一个Result结构体用于存储响应结果
	var res Result
	// 将HTTP响应解析为Result结构体
	err = json.NewDecoder(httpResp.Body).Decode(&res)
	if err != nil {
		//转json为结构体失败
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		// 获取 access token 失败
		return domain.WechatInfo{}, fmt.Errorf("调用微信接口失败 errcode %d, errmsg %s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		UnionId: res.UnionId,
		OpenId:  res.OpenId,
	}, nil
}

// AuthURL 生成用于微信登录的授权URL
func (s *service) AuthURL(ctx context.Context, state string) (string, error) {

	// 定义授权URL的模板，其中包含多个占位符
	const authURLPattern = `https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`
	// 使用fmt.Sprintf函数将占位符替换为实际的值
	// s.appID: 服务实例中的微信应用ID
	// redirectURL: 重定向URL，这里应该是预先定义好的
	// state: 随机生成的状态字符串，用于防止CSRF攻击
	return fmt.Sprintf(authURLPattern, s.appID, redirectURL, state), nil
}

type Result struct {
	AccessToken string `json:"access_token"`
	// access_token接口调用凭证超时时间，单位（秒）
	ExpiresIn int64 `json:"expires_in"`
	// 用户刷新access_token
	RefreshToken string `json:"refresh_token"`
	// 授权用户唯一标识
	OpenId string `json:"openid"`
	// 用户授权的作用域，使用逗号（,）分隔
	Scope string `json:"scope"`
	// 当且仅当该网站应用已获得该用户的userinfo授权时，才会出现该字段。
	UnionId string `json:"unionid"`

	// 错误返回
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
