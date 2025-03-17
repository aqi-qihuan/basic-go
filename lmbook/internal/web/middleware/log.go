package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

// LogMiddlewareBuilder 是一个用于构建日志中间件的构建器结构体。
// 它包含了配置日志中间件所需的字段和方法。
type LogMiddlewareBuilder struct {
	// logFn 是一个函数类型字段，用于记录访问日志。
	// 它接收两个参数：一个是 context.Context，用于传递请求的上下文信息；
	// 另一个是 AccessLog，用于存储访问日志的具体内容。
	logFn func(ctx context.Context, l AccessLog)
	// allowReqBody 是一个布尔类型字段，用于指示是否允许记录请求体。
	// 如果为 true，则日志中间件会记录请求体的内容；如果为 false，则不会记录。
	allowReqBody bool
	// allowRespBody 是一个布尔类型字段，用于指示是否允许记录响应体。
	// 如果为 true，则日志中间件会记录响应体的内容；如果为 false，则不会记录。
	allowRespBody bool
}

// NewLogMiddlewareBuilder 是一个工厂函数，用于创建一个新的 LogMiddlewareBuilder 实例。
// 该函数接收一个 logFn 参数，这是一个函数类型，用于记录访问日志。
// logFn 函数接收两个参数：一个是 context.Context 类型，用于传递请求上下文；
// 另一个是 AccessLog 类型，用于传递访问日志信息。
// 函数返回一个指向 LogMiddlewareBuilder 实例的指针。
func NewLogMiddlewareBuilder(logFn func(ctx context.Context, l AccessLog)) *LogMiddlewareBuilder {
	// 返回一个新的 LogMiddlewareBuilder 实例，并将传入的 logFn 赋值给实例的 logFn 字段。
	return &LogMiddlewareBuilder{
		logFn: logFn,
	}
}

// 定义一个名为 AllowReqBody 的方法，它属于 LogMiddlewareBuilder 类型
func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	// 将 LogMiddlewareBuilder 实例的 allowReqBody 字段设置为 true
	l.allowReqBody = true
	// 返回当前 LogMiddlewareBuilder 实例，以便进行链式调用
	return l
}

// 定义一个方法 AllowRespBody，它属于 LogMiddlewareBuilder 类型
func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	// 将 LogMiddlewareBuilder 实例的 allowRespBody 字段设置为 true
	// 这表示允许记录响应体
	l.allowRespBody = true
	// 返回当前 LogMiddlewareBuilder 实例，以便进行链式调用
	return l
}

// LogMiddlewareBuilder 是一个用于构建日志中间件的构建器
func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	// 返回一个 gin.HandlerFunc 类型的函数，用于处理 HTTP 请求
	return func(ctx *gin.Context) {
		// 获取请求的路径
		path := ctx.Request.URL.Path
		// 如果路径长度超过 1024 个字符，则截取前 1024 个字符
		if len(path) > 1024 {
			path = path[:1024]
		}
		// 获取请求的方法（如 GET、POST 等）
		method := ctx.Request.Method
		// 创建一个 AccessLog 对象，用于记录访问日志
		al := AccessLog{
			Path:   path,
			Method: method,
		}
		// 如果允许记录请求体
		if l.allowReqBody {
			// Request.Body 是一个 Stream 对象，只能读一次
			body, _ := ctx.GetRawData()
			// 如果请求体长度超过 2048 个字符，则截取前 2048 个字符
			if len(body) > 2048 {
				al.ReqBody = string(body[:2048])
			} else {
				al.ReqBody = string(body)
			}
			// 放回去
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			//ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		start := time.Now()

		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}

		defer func() {
			al.Duration = time.Since(start)
			//duration := time.Now().Sub(start)
			l.logFn(ctx, al)
		}()

		// 直接执行下一个 middleware...直到业务逻辑
		ctx.Next()
		// 在这里，你就拿到了响应
	}
}

type AccessLog struct {
	Path     string        `json:"path"`
	Method   string        `json:"method"`
	ReqBody  string        `json:"req_body"`
	Status   int           `json:"status"`
	RespBody string        `json:"resp_body"`
	Duration time.Duration `json:"duration"`
}

// 定义一个结构体 responseWriter，它嵌入了一个 gin.ResponseWriter 接口和一个 *AccessLog 指针
type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

// Write 方法实现了 io.Writer 接口，用于写入数据到响应体
// 参数 data 是要写入的数据，返回写入的字节数和可能发生的错误
func (w *responseWriter) Write(data []byte) (int, error) {
	// 将写入的数据转换为字符串并赋值给 al.RespBody，用于记录响应体
	w.al.RespBody = string(data)
	// 调用嵌入的 gin.ResponseWriter 的 Write 方法，实际执行数据写入操作
	return w.ResponseWriter.Write(data)
}

// WriteHeader 方法用于设置 HTTP 响应的状态码
// 参数 statusCode 是要设置的状态码
func (w *responseWriter) WriteHeader(statusCode int) {
	// 将状态码赋值给 al.Status，用于记录响应状态码
	w.al.Status = statusCode
	// 调用嵌入的 gin.ResponseWriter 的 WriteHeader 方法，实际设置 HTTP 响应的状态码
	w.ResponseWriter.WriteHeader(statusCode)
}
