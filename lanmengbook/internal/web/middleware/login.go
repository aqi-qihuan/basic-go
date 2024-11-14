package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// 定义一个LoginMiddlewareBuilder结构体
type LoginMiddlewareBuilder struct {
}

// CheckLogin方法用于检查用户是否已经登录
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	//注册一下这个类型
	gob.Register(time.Now())
	// 返回一个gin.HandlerFunc类型的函数
	return func(ctx *gin.Context) {
		// 获取请求的URL路径
		path := ctx.Request.URL.Path
		// 如果路径是/users/signup或/users/login，则不需要登录校验
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		// 获取session
		sess := sessions.Default(ctx)

		// 尝试获取session中的userId
		userId := sess.Get("userId")
		// 如果session中没有userId，则表示用户未登录
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()

		//我怎么知道我刷新了？
		//假如说，我们的策略是没分组刷新一次，我怎么知道已经过了一分钟？
		const updateTimeKey = "update_time"
		//试着拿出上一次的刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || (!ok) || now.Sub(lastUpdateTime) > time.Second*10 {
			//你这是第一次进来
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				//打印日志
				fmt.Println(err)
			}
		}

	}
}
