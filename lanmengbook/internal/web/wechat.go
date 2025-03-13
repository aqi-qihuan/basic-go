package web

import (
	"basic-go/lanmengbook/internal/service"
	"basic-go/lanmengbook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
)

type OAuth2WechatHandler struct {
	svc         wechat.Service
	userService service.UserService
	//ijwt.Handler
	key             []byte
	stateCookieName string
}

func NewOAuth2WechatHandler(svc wechat.Service, userService service.UserService, key []byte, stateCookieName string) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:             svc,
		userService:     userService,
		key:             key,
		stateCookieName: stateCookieName,
	}

}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/auth", h.AuthURL)
	g.GET("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(context *gin.Context) {

}

func (h *OAuth2WechatHandler) Callback(context *gin.Context) {

}

//type OAuth2Handler struct {
//}
//func (h *OAuth2Handler) RegisterRoutes(server *gin.Engine) {
//g := server.Group("/oauth2")
//g.GET("/:platform/auth", h.AuthURL)
//g.GET("/:platform/callback", h.Callback)
//}
//
//func (h *OAuth2Handler) AuthURL(context *gin.Context) {
//	pla
//}
//
//func (h *OAuth2Handler) Callback(context *gin.Context) {
//
//
//}
