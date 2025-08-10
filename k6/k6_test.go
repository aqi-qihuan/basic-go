package k6

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHello(t *testing.T) {
	server := gin.Default()
	server.POST("/hello", func(ctx *gin.Context) {
		var u User
		ctx.Bind(&u)
		r := rand.Int31n(1000)
		time.Sleep(time.Millisecond * time.Duration(r))
		// 这里我们模拟一下错误
		// 模拟 10% 比例的错误
		if r%100 < 10 {
			ctx.String(http.StatusInternalServerError, "系统错误")
		} else {
			ctx.String(http.StatusOK, u.Name)
		}
	})
	server.Run(":8080")
}

type User struct {
	Name string `json:"name"`
}
