package main

import (
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/service"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/web"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	// 初始化数据库
	db := initDB()

	// 初始化web服务器
	server := initWebServer()

	// 初始化用户处理
	initUserHdl(db, server)

	// 启动服务器
	server.Run(":8080")
}

func initUserHdl(db *gorm.DB, server *gin.Engine) {
	// 创建用户数据访问对象
	ud := dao.NewUserDAO(db)

	// 创建用户仓库
	ur := repository.NewUserRepository(ud)

	// 创建用户服务
	us := service.NewUserService(ur)

	// 创建用户处理器
	hdl := web.NewUserHandler(us)

	// 注册路由
	hdl.RegisterRoutes(server)
}

func initDB() *gorm.DB {
	// 连接数据库
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/lanmengbook"))
	if err != nil {
		panic(err)
	}

	// 初始化表
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	// 创建gin引擎
	server := gin.Default()

	// 跨域配置
	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,

		AllowHeaders: []string{"Content-Type"},
		//AllowHeaders: []string{"content-type"},
		//AllowMethods: []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的 Middleware")
	})

	// 登录中间件
	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
	return server
}
