package web

import (
	"gitee.com/geekbang/basic-go/lanmengbook/internal/domain"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	// 定义正则表达式，用于匹配邮箱地址
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 定义正则表达式，用于匹配密码
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

// UserHandler 结构体，包含处理用户请求所需的正则表达式和服务
type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            *service.UserService
}

// NewUserHandler 函数，用于创建一个新的 UserHandler 实例
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
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
	// POST /users/edit
	ug.POST("/edit", h.Edit)
	// GET /users/profile
	ug.GET("/profile", h.Profile)
}

// SignUp 方法，用于处理用户注册请求
func (h *UserHandler) SignUp(ctx *gin.Context) {
	// 定义一个结构体，用于接收用户注册请求的数据
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	// 声明一个 SignUpReq 类型的变量 req
	var req SignUpReq
	// 绑定请求体到 req 变量，如果绑定失败则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 使用正则表达式匹配邮箱地址，如果匹配失败则返回错误信息
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	// 检查两次输入的密码是否一致，如果不一致则返回错误信息
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	// 使用正则表达式匹配密码，如果匹配失败则返回错误信息
	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	// 调用服务层的 Signup 方法进行注册，如果注册失败则根据错误类型返回相应的错误信息
	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

// LoginJWT 方法，用于处理用户登录请求，并返回 JWT 令牌
func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	// 定义一个结构体，用于接收用户登录请求的数据
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// 声明一个 Req 类型的变量 req
	var req Req
	// 绑定请求体到 req 变量，如果绑定失败则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用服务层的 Login 方法进行登录，如果登录失败则根据错误类型返回相应的错误信息
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 创建一个 UserClaims 结构体，用于存储 JWT 令牌的声明信息
		uc := UserClaims{
			Uid:       u.Id,
			UserAgent: ctx.GetHeader("User-Agent"),
			RegisteredClaims: jwt.RegisteredClaims{
				// 1 分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		}
		// 创建一个新的 JWT 令牌，使用 HS512 签名算法，并将 UserClaims 结构体作为声明信息
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		// 使用 JWTKey 对 JWT 令牌进行签名，并将签名后的令牌字符串赋值给 tokenStr 变量
		tokenStr, err := token.SignedString(JWTKey)
		if err != nil {
			// 如果签名过程中发生错误，则返回系统错误信息
			ctx.String(http.StatusOK, "系统错误")
		}
		// 在响应头中添加 x-jwt-token 字段，用于存储 JWT 令牌
		ctx.Header("x-jwt-token", tokenStr)
		// 返回登录成功的信息
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		// 如果用户名或密码不正确，则返回相应的错误信息
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		// 如果发生其他错误，则返回系统错误信息
		ctx.String(http.StatusOK, "系统错误")
	}
}

// Login 方法，用于处理用户登录请求，并使用会话（session）来存储用户信息
func (h *UserHandler) Login(ctx *gin.Context) {
	// 定义一个结构体，用于接收用户登录请求的数据
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// 声明一个 Req 类型的变量 req
	var req Req
	// 绑定请求体到 req 变量，如果绑定失败则直接返回
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 调用服务层的 Login 方法进行登录，如果登录失败则根据错误类型返回相应的错误信息
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 从 gin 上下文中获取默认的会话对象
		sess := sessions.Default(ctx)
		// 在会话中设置 userId 字段，用于存储用户的 ID
		sess.Set("userId", u.Id)
		// 设置会话的过期时间为 10 分钟
		sess.Options(sessions.Options{
			// 十分钟
			MaxAge: 30,
		})
		// 保存会话信息
		err = sess.Save()
		if err != nil {
			// 如果保存会话信息时发生错误，则返回系统错误信息
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		// 返回登录成功的信息
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		// 如果用户名或密码不正确，则返回相应的错误信息
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		// 如果发生其他错误，则返回系统错误信息
		ctx.String(http.StatusOK, "系统错误")
	}
}

// Edit 方法，用于处理用户编辑请求
func (h *UserHandler) Edit(ctx *gin.Context) {
	// 嵌入一段刷新过期时间的代码
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	// 从 gin 上下文中获取用户信息，此处注释掉了，因为没有实际获取到用户信息
	//us := ctx.MustGet("user").(UserClaims)
	// 返回一个字符串，表示这是个人信息页面
	ctx.String(http.StatusOK, "这是 profile")
	// 嵌入一段刷新过期时间的代码
}

// JWTKey 是用于验证 JWT 令牌的密钥
var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

// UserClaims 是一个结构体，用于存储 JWT 令牌中的用户声明信息
type UserClaims struct {
	// 嵌入了 jwt.RegisteredClaims 结构体，用于存储标准的 JWT 声明信息
	jwt.RegisteredClaims
	// Uid 是用户的唯一标识符
	Uid int64
	// UserAgent 是用户的浏览器或客户端的标识信息
	UserAgent string
}
