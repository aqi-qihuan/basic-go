package web

import (
	"basic-go/lmbook/internal/service"
	"basic-go/lmbook/internal/service/oauth2/wechat"
	ijwt "basic-go/lmbook/internal/web/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

// OAuth2WechatHandler 是一个处理微信 OAuth2 认证的结构体
type OAuth2WechatHandler struct {
	// svc 是微信服务的实例，用于与微信平台进行交互
	svc wechat.Service
	// userSvc 是用户服务的实例，用于处理用户相关的业务逻辑
	userSvc service.UserService
	// Ijwt.Handler 是一个嵌入字段，表示 OAuth2WechatHandler 实现了 Ijwt.Handler 接口
	ijwt.Handler
	// key 是用于加密和解密数据的密钥
	key []byte
	// stateCookieName 是存储在客户端的 Cookie 名称，用于防止 CSRF 攻击
	stateCookieName string
}

// NewOAuth2WechatHandler 创建一个新的OAuth2WechatHandler实例
// 参数：
//   - svc: wechat.Service类型，用于处理微信相关的服务
//   - hdl: ijwt.Handler类型，用于处理JWT相关的操作
//   - userSvc: service.UserService类型，用于处理用户相关的服务
//
// 返回值：
//   - *OAuth2WechatHandler: 返回一个指向OAuth2WechatHandler实例的指针
func NewOAuth2WechatHandler(svc wechat.Service,
	hdl ijwt.Handler,
	userSvc service.UserService) *OAuth2WechatHandler {
	// 返回一个新的OAuth2WechatHandler实例，并初始化其字段
	return &OAuth2WechatHandler{
		svc:             svc,                                        // 微信服务实例
		userSvc:         userSvc,                                    // 用户服务实例
		key:             []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgB"), // 用于加密的密钥
		stateCookieName: "jwt-state",                                // 用于存储状态的Cookie名称
		Handler:         hdl,                                        // JWT处理实例
	}
}

// 定义一个名为 RegisterRoutes 的方法，该方法绑定到 OAuth2WechatHandler 类型
func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	// 创建一个路由组，路径前缀为 "/oauth2/wechat"
	g := server.Group("/oauth2/wechat")
	// 在该路由组中定义一个 GET 请求的路由，路径为 "/authurl"
	// 当访问 "/oauth2/wechat/authurl" 时，调用 o.Auth2URL 方法处理请求
	g.GET("/authurl", o.Auth2URL)
	// 在该路由组中定义一个任意方法（GET、POST、PUT、DELETE 等）的路由，路径为 "/callback"
	// 当访问 "/oauth2/wechat/callback" 时，调用 o.Callback 方法处理请求
	g.Any("/callback", o.Callback)
}

// Auth2URL 处理微信OAuth2认证的URL生成请求
func (o *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	// 生成一个随机的state值，用于防止CSRF攻击
	state := uuid.New()
	// 调用服务层的AuthURL方法，生成认证URL，并传递state值
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		// 如果生成认证URL失败，返回错误信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "构造跳转URL失败",
			Code: 5,
		})
		return
	}
	// 将state值存储在Cookie中，以便后续验证
	err = o.setStateCookie(ctx, state)
	if err != nil {
		// 如果存储Cookie失败，返回服务器异常信息
		ctx.JSON(http.StatusOK, Result{
			Msg:  "服务器异常",
			Code: 5,
		})
	}
	// 返回生成的认证URL
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

// Callback 处理微信OAuth2回调请求
func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	// 验证请求的state参数，确保请求的合法性
	err := o.verifyState(ctx)
	if err != nil {
		// 如果验证失败，返回非法请求的响应
		ctx.JSON(http.StatusOK, Result{
			Msg:  "非法请求",
			Code: 4,
		})
		return
	}
	// 你校验不校验都可以
	code := ctx.Query("code")
	// state := ctx.Query("state")
	wechatInfo, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "授权码有误",
			Code: 4,
		})
		return
	}
	u, err := o.userSvc.FindOrCreateByWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		return
	}
	err = o.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	return
}

// verifyState 验证 OAuth2 授权流程中的 state 参数，确保请求的合法性
func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	// 从请求中获取 state 参数
	state := ctx.Query("state")
	// 尝试从请求中获取名为 o.stateCookieName 的 cookie
	ck, err := ctx.Cookie(o.stateCookieName)
	if err != nil {
		// 如果获取 cookie 失败，返回错误信息
		return fmt.Errorf("无法获得 cookie %w", err)
	}
	// 定义一个 StateClaims 结构体用于存储解析后的 token 信息
	var sc StateClaims
	// 使用 jwt.ParseWithClaims 解析 cookie 中的 token
	// jwt.ParseWithClaims 函数会验证 token 的签名并解析其内容
	// 第二个参数是一个指向 StateClaims 的指针，用于存储解析后的数据
	// 第三个参数是一个回调函数，用于提供用于验证 token 签名的密钥
	_, err = jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		// 返回用于验证 token 签名的密钥
		return o.key, nil
	})
	if err != nil {
		// 如果解析 token 失败，返回错误信息
		return fmt.Errorf("解析 token 失败 %w", err)
	}
	// 比较请求中的 state 参数和解析后的 token 中的 state 值
	if state != sc.State {
		// state 不匹配，有人搞你
		return fmt.Errorf("state 不匹配")
	}
	return nil
}

// setStateCookie 为OAuth2WechatHandler设置状态Cookie
func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context,
	state string) error {
	// 创建StateClaims结构体实例，包含状态信息
	claims := StateClaims{
		State: state,
	}
	// 使用HS512算法创建一个新的JWT token，并携带claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	// 使用OAuth2WechatHandler的密钥对token进行签名，生成字符串
	tokenStr, err := token.SignedString(o.key)
	if err != nil {

		// 如果签名过程中出现错误，返回错误
		return err
	}
	// 在HTTP响应中设置Cookie，Cookie名为o.stateCookieName
	// tokenStr为Cookie的值，过期时间为600秒
	// 路径为"/oauth2/wechat/callback"，不设置域名，不设置HttpOnly，设置为Secure
	ctx.SetCookie(o.stateCookieName, tokenStr,
		600, "/oauth2/wechat/callback",
		"", false, true)
	// 返回nil表示操作成功
	return nil
}

// 定义一个名为 StateClaims 的结构体，用于存储 JWT 的声明信息
type StateClaims struct {
	// 嵌入 jwt.RegisteredClaims 结构体，继承其所有字段和方法
	// jwt.RegisteredClaims 是一个标准 JWT 声明结构体，包含常见的注册声明（如 exp, iat, sub 等）
	jwt.RegisteredClaims
	// State 字段，用于存储自定义的状态信息
	// 这个字段可以用于存储与 JWT 相关的任意状态数据，例如 OAuth 2.0 中的 state 参数
	State string
}
