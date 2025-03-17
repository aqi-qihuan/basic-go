package web

import (
	"basic-go/lmbook/internal/domain"
	"basic-go/lmbook/internal/service"
	ijwt "basic-go/lmbook/internal/web/jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

const (
	// 定义正则表达式，用于验证邮箱格式
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 定义正则表达式，用于验证密码格式，密码必须包含字母、数字、特殊字符，并且不少于八位
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	// 定义业务名称，用于标识登录业务
	bizLogin = "login"
)

// UserHandler 结构体，用于处理用户相关的 HTTP 请求
type UserHandler struct {
	// ijwt.Handler 是一个接口类型，用于定义处理 JWT（JSON Web Token）相关操作的处理器。
	// 该接口可以包含多个方法，用于处理不同类型的 JWT 操作，例如生成、验证、解析等。
	ijwt.Handler
	// 用于验证邮箱格式的正则表达式
	emailRexExp *regexp.Regexp
	// 用于验证密码格式的正则表达式
	passwordRexExp *regexp.Regexp
	// 用户服务接口，用于处理用户相关的业务逻辑
	svc service.UserService
	// 验证码服务接口，用于处理验证码相关的业务逻辑
	codeSvc service.CodeService
}

// NewUserHandler 创建一个新的用户处理器实例
// 参数 svc 是用户服务接口的实现，用于处理用户相关的业务逻辑
// 参数 hdl 是 JWT 处理器接口的实现，用于处理 JWT 相关的操作
// 参数 codeSvc 是验证码服务接口的实现，用于处理验证码相关的业务逻辑
// 返回一个指向 UserHandler 实例的指针
func NewUserHandler(svc service.UserService,
	hdl ijwt.Handler,
	codeSvc service.CodeService) *UserHandler {
	// 返回一个新的 UserHandler 实例，并初始化其字段
	return &UserHandler{
		// 使用正则表达式编译器编译电子邮件正则表达式模式
		// emailRegexPattern 是一个字符串常量，定义了电子邮件的格式
		// regexp.None 表示不使用任何特殊的匹配选项
		emailRexExp: regexp.MustCompile(emailRegexPattern, regexp.None),
		// 使用正则表达式编译器编译密码正则表达式模式
		// passwordRegexPattern 是一个字符串常量，定义了密码的格式
		// regexp.None 表示不使用任何特殊的匹配选项
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		// 将传入的用户服务接口实现赋值给 UserHandler 的 svc 字段
		svc: svc,
		// 将传入的验证码服务接口实现赋值给 UserHandler 的 codeSvc 字段
		codeSvc: codeSvc,
		// 将传入的 JWT 处理器接口实现赋值给 UserHandler 的 Handler 字段
		Handler: hdl,
	}
}

// RegisterRoutes 方法，用于注册用户相关的路由
func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// REST 风格
	//server.POST("/user", h.SignUp)
	//server.PUT("/user", h.SignUp)
	//server.GET("/users/:username", h.Profile)
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", h.SignUp)
	// POST /users/login
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	ug.POST("/logout", h.LogoutJWT)
	// POST /users/edit
	ug.POST("/edit", h.Edit)
	// GET /users/profile
	ug.GET("/profile", h.Profile)
	ug.GET("/refresh_token", h.RefreshToken)

	// 手机验证码登录相关功能
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

// LoginSMS 方法，用于处理手机验证码登录的 HTTP 请求
func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	// 定义一个结构体 Req，包含两个字段：Phone 和 Code，分别用于接收请求中的手机号和验证码
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	// 声明一个 Req 类型的变量 req，用于接收解析后的请求参数
	var req Req
	// 调用 ctx.Bind 方法，将请求参数绑定到 req 变量上，如果发生错误则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 调用 codeSvc.Verify 方法，验证手机号和验证码是否匹配，如果发生错误则返回错误信息
	ok, err := h.codeSvc.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		// 返回 JSON 格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		zap.L().Error("验证码校验失败",
			// 在生产环境绝对不能打
			// 开发环境你可以随便打
			//zap.String("phone", req.Phone),
			zap.Error(err))
		return
	}
	// 如果验证码验证失败，则返回错误信息
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码不对，请重新输入",
		})
		return
	}
	// 调用 svc.FindOrCreate 方法，根据手机号查找或创建用户，如果发生错误则返回错误信息
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 返回 JSON 格式的成功信息，状态码为 200，表示请求成功，且业务逻辑处理成功
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

// SendSMSLoginCode 方法，用于处理发送手机验证码的 HTTP 请求
func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	// 定义一个结构体 Req，包含一个字段：Phone，用于接收请求中的手机号
	type Req struct {
		Phone string `json:"phone"`
	}
	// 声明一个 Req 类型的变量 req，用于接收解析后的请求参数
	var req Req
	// 调用 ctx.Bind 方法，将请求参数绑定到 req 变量上，如果发生错误则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	// 如果手机号为空，则返回错误信息
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	// 调用 codeSvc.Send 方法，发送手机验证码，如果发生错误则返回错误信息
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		// 如果发送成功，则返回 JSON 格式的成功信息，状态码为 200，表示请求成功，且业务逻辑处理成功
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		// 如果发送太频繁，则返回 JSON 格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		zap.L().Error("频繁发送验证码")
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		// 如果发生其他错误，则返回 JSON 格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 补日志的
		log.Println(err)
	}
}

// SignUp 方法，用于处理用户注册的 HTTP 请求
func (h *UserHandler) SignUp(ctx *gin.Context) {
	// 定义一个结构体 SignUpReq，包含三个字段：Email、Password 和 ConfirmPassword，分别用于接收请求中的邮箱、密码和确认密码
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	// 声明一个 SignUpReq 类型的变量 req，用于接收解析后的请求参数
	var req SignUpReq
	// 调用 ctx.Bind 方法，将请求参数绑定到 req 变量上，如果发生错误则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 调用 emailRexExp.MatchString 方法，验证邮箱格式是否合法，如果发生错误则返回错误信息
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		// 返回字符串格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 如果邮箱格式不合法，则返回错误信息
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	// 如果密码和确认密码不一致，则返回错误信息
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	// 调用 passwordRexExp.MatchString 方法，验证密码格式是否合法，如果发生错误则返回错误信息
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		// 返回字符串格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 如果密码格式不合法，则返回错误信息
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	// 调用 svc.Signup 方法，注册用户，如果发生错误则返回错误信息
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		// 如果注册成功，则返回字符串格式的成功信息，状态码为 200，表示请求成功，且业务逻辑处理成功
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		// 如果邮箱已存在，则返回字符串格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		// 如果发生其他错误，则返回字符串格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.String(http.StatusOK, "系统错误")
	}
}

// LoginJWT 处理用户JWT登录请求
func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	// 定义请求结构体，用于接收JSON格式的用户名和密码
	type Req struct {
		Email    string `json:"email"`    // 用户邮箱
		Password string `json:"password"` // 用户密码
	}
	var req Req
	// 将请求体绑定到req结构体中，如果绑定失败则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用服务层的Login方法进行用户验证
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	// 根据不同的错误类型进行不同的处理
	switch err {
	case nil:
		// 如果验证成功，设置登录Token
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			// 如果设置Token失败，返回系统错误
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		// 返回登录成功信息
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		// 如果用户名或密码错误，返回错误信息
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		// 其他错误情况，返回系统错误
		ctx.String(http.StatusOK, "系统错误")
	}
}

// setJWTToken 方法，用于为用户设置 JWT 令牌
//func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
//	// 创建一个 UserClaims 结构体实例，用于存储 JWT 声明
//	uc := UserClaims{
//		// 设置用户 ID
//		Uid: uid,
//		// 从请求头中获取 User-Agent 信息
//		UserAgent: ctx.GetHeader("User-Agent"),
//		// 设置注册声明
//		RegisteredClaims: jwt.RegisteredClaims{
//			// 设置过期时间为当前时间加 30 分钟
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
//		},
//	}
//	// 创建一个新的 JWT 对象，使用 HS512 签名方法和自定义声明
//	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
//	// 对 JWT 进行签名，生成 token 字符串
//	tokenStr, err := token.SignedString(JWTKey)
//	// 如果发生错误，返回系统错误信息
//	if err != nil {
//		ctx.String(http.StatusOK, "系统错误")
//	}
//	// 在响应头中设置 x-jwt-token 字段，值为生成的 token 字符串
//	ctx.Header("x-jwt-token", tokenStr)
//}

// Login 方法，用于处理用户通过邮箱和密码登录的 HTTP 请求
func (h *UserHandler) Login(ctx *gin.Context) {
	// 定义一个结构体 Req，包含两个字段：Email 和 Password，分别用于接收请求中的邮箱和密码
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// 声明一个 Req 类型的变量 req，用于接收解析后的请求参数
	var req Req
	// 调用 ctx.Bind 方法，将请求参数绑定到 req 变量上，如果发生错误则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用 svc.Login 方法，根据邮箱和密码登录用户，如果发生错误则返回错误信息
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 如果登录成功，则获取默认的 session 对象
		sess := sessions.Default(ctx)
		// 在 session 中设置 userId 字段，值为用户的 ID
		sess.Set("userId", u.Id)
		// 设置 session 的过期时间为 30 分钟
		sess.Options(sessions.Options{
			// 十分钟
			MaxAge: 30,
		})
		// 保存 session
		err = sess.Save()
		if err != nil {
			// 如果保存 session 时发生错误，则返回系统错误信息
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		// 返回字符串格式的成功信息，状态码为 200，表示请求成功，且业务逻辑处理成功
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		// 如果用户名或密码不正确，则返回字符串格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		// 如果发生其他错误，则返回字符串格式的错误信息，状态码为 200，表示请求成功，但业务逻辑处理失败
		ctx.String(http.StatusOK, "系统错误")
	}
}

// Edit 方法，用于处理用户编辑个人信息的 HTTP 请求
func (h *UserHandler) Edit(ctx *gin.Context) {
	// 定义一个结构体 Req，包含三个字段：Nickname、Birthday 和 AboutMe，分别用于接收请求中的昵称、生日和个人简介
	type Req struct {
		// 改邮箱，密码，或者能不能改手机号
		Nickname string `json:"nickname"`
		// YYYY-MM-DD
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	// 声明一个 Req 类型的变量 req，用于接收解析后的请求参数
	var req Req
	// 调用 ctx.Bind 方法，将请求参数绑定到 req 变量上，如果发生错误则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 从请求上下文中获取用户信息，如果获取失败则返回未授权状态码
	uc, ok := ctx.MustGet("user").(ijwt.UserClaims)
	if !ok {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 解析用户输入的生日字符串，如果格式不正确则返回错误信息
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	// 调用 svc.UpdateNonSensitiveInfo 方法，更新用户的非敏感信息，如果发生错误则返回系统异常信息
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Uid,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	// 如果更新成功，则返回字符串格式的成功信息，状态码为 200，表示请求成功，且业务逻辑处理成功
	ctx.String(http.StatusOK, "更新成功")
}

// Profile 方法，用于处理用户查看个人资料的 HTTP 请求
func (h *UserHandler) Profile(ctx *gin.Context) {
	// 尝试从上下文中获取用户信息，如果获取失败则返回未授权状态码
	uc, ok := ctx.MustGet("user").(ijwt.UserClaims)
	if !ok {
		//ctx.String(http.StatusOK, "系统错误")
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 根据用户 ID 查询用户信息，如果发生错误则返回系统异常信息
	u, err := h.svc.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	// 定义一个结构体 User，用于封装用户信息
	type User struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"aboutMe"`
		Birthday string `json:"birthday"`
	}
	// 将用户信息封装成 JSON 格式并返回，如果成功则返回状态码 200
	ctx.JSON(http.StatusOK, User{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	})
}

// JWTKey 是用于验证 JWT 的密钥
//var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	//约定，前端在 Authorization 里面带上这个 refresh_token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// token 无效或者 redis 有问题
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

// UserClaims 结构体定义了 JWT 中包含的用户声明
//type UserClaims struct {
//	// 嵌入 jwt.RegisteredClaims 以包含标准的 JWT 声明
//	jwt.RegisteredClaims
//	// Uid 是用户的唯一标识符
//	Uid int64
//	// UserAgent 是用户的客户端信息
//	UserAgent string
//}

// LogoutJWT 处理用户退出登录的请求，使用JWT进行身份验证
func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	// 调用ClearToken方法清除用户的Token
	err := h.ClearToken(ctx)
	// 如果清除Token过程中发生错误
	if err != nil {
		// 返回状态码200和错误信息，Code: 5表示系统异常
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统异常"})
		return
	}
	// 如果清除Token成功，返回状态码200和成功信息
	ctx.JSON(http.StatusOK, Result{Msg: "退出登录成功"})
}
