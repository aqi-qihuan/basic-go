package wechat

import (
	"basic-go/lanmengbook/internal/domain"
	"context"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm/logger"
	"net/url"
)

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	appID string
	l     logger.LoggerV1
}

func NewService(appID, appSecret string) Service {
	return &service{
		appID: appID,
	}
}
func (s *service) AuthURL(ctx context.Context, state string) (string, error) {

	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect`"
	statec := uuid.New()
	return fmt.Sprintf(urlPattern, s.appID, redirectURL, statec), nil
}
