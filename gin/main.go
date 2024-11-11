package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建一个默认的gin实例
	server := gin.Default()
	// 使用两个中间件
	server.Use(func(context *gin.Context) {
		println("这是第一个 Middleware")
	}, func(context *gin.Context) {
		println("这是第二个 Middleware")
	})
	// 定义一个GET请求，路径为/hello
	server.GET("/hello", func(ctx *gin.Context) {
		// 返回字符串"hello, world"
		ctx.String(http.StatusOK, "hello, world")
	})

	// 参数路由，路径参数
	// 定义一个GET请求，路径为/users/:name
	server.GET("/users/:name", func(ctx *gin.Context) {
		// 获取路径参数name
		name := ctx.Param("name")
		// 返回字符串"hello, "+name
		ctx.String(http.StatusOK, "hello, "+name)
	})

	// 查询参数
	// GET /order?id=123
	// 定义一个GET请求，路径为/order
	server.GET("/order", func(ctx *gin.Context) {
		// 获取请求参数id
		id := ctx.Query("id")
		// 返回订单ID
		ctx.String(http.StatusOK, "订单 ID 是 "+id)
	})

	// 定义一个GET请求，路径为/views/*.html
	server.GET("/views/*.html", func(ctx *gin.Context) {
		// 获取路径参数.html
		view := ctx.Param(".html")
		// 返回字符串"view 是 "+view
		ctx.String(http.StatusOK, "view 是 "+view)
	})

	//server.GET("/star/*", func(ctx *gin.Context) {
	//	view := ctx.Param(".html")
	//	ctx.String(http.StatusOK, "view 是 "+view)
	//})
	//server.GET("/star/*/abc", func(ctx *gin.Context) {
	//	view := ctx.Param(".html")
	//	ctx.String(http.StatusOK, "view 是 "+view)
	//})

	// 定义一个POST请求，路径为/login
	server.POST("/login", func(ctx *gin.Context) {
		// 返回字符串"hello, login"
		ctx.String(http.StatusOK, "hello, login")
	})

	//go func() {
	//	server1 := gin.Default()
	//
	//	server1.Run(":8081")
	//}()

	// 如果你不传参数，那么实际上监听的是 8080 端口
	server.Run(":8080")
	// 这种写法是错的
	//  missing port in address
	//server.Run("8080")
}
