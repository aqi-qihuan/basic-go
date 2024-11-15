package main

import (
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/service"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/web"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/web/middleware"
	"gitee.com/geekbang/basic-go/lanmengbook/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 这个是允许前端访问你的后端响应中带的头部
		ExposeHeaders: []string{"x-jwt-token"},
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

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	server.Use(ratelimit.NewBuilder(redisClient,
		time.Second, 1).Build())

	useJWT(server)
	//useSession(server)
	return server
}

func useJWT(server *gin.Engine) {
	login := middleware.LoginJWTMiddlewareBuilder{}
	server.Use(login.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("secret"))
	// 基于内存的实现
	//store := memstore.NewStore([]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
	//	[]byte("eF1`yQ9>yT1`tH1,sJ0.zD8;mZ9~nC6("))
	//store, err := redis.NewStore(16, "tcp",
	//	"localhost:6379", "",
	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgK"),
	//	[]byte("k6CswdUm75WKcbM68UQUuxVsHSpTCwgA"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
