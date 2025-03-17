package logger

import "go.uber.org/zap"

// ZapLogger 是一个包装了 zap.Logger 的结构体，用于提供更方便的日志记录方法
type ZapLogger struct {
	l *zap.Logger // l 是一个指向 zap.Logger 的指针，用于实际的日志记录操作
}

// NewZapLogger 是一个构造函数，用于创建一个新的 ZapLogger 实例
// 参数 l 是一个指向 zap.Logger 的指针
// 返回一个新的 ZapLogger 实例
func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{
		l: l, // 将传入的 zap.Logger 赋值给 ZapLogger 结构体的 l 字段
	}
}

// Debug 方法用于记录调试级别的日志
// 参数 msg 是日志消息
// 参数 args 是一个可变参数列表，用于传递额外的日志字段
func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toArgs(args)...) // 调用 zap.Logger 的 Debug 方法，并传递转换后的日志字段
}

// Info 方法用于记录信息级别的日志
// 参数 msg 是日志消息
// 参数 args 是一个可变参数列表，用于传递额外的日志字段
func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toArgs(args)...) // 调用 zap.Logger 的 Info 方法，并传递转换后的日志字段
}

// Warn 方法用于记录警告级别的日志
// 参数 msg 是日志消息
// 参数 args 是一个可变参数列表，用于传递额外的日志字段
func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toArgs(args)...) // 调用 zap.Logger 的 Warn 方法，并传递转换后的日志字段
}

// Error 方法用于记录错误级别的日志
// 参数 msg 是日志消息
// 参数 args 是一个可变参数列表，用于传递额外的日志字段
func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toArgs(args)...) // 调用 zap.Logger 的 Error 方法，并传递转换后的日志字段
}

// toArgs 方法用于将自定义的 Field 切片转换为 zap.Field 切片
// 参数 args 是一个自定义的 Field 切片
// 返回一个 zap.Field 切片
func (z *ZapLogger) toArgs(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args)) // 创建一个与 args 长度相同的 zap.Field 切片
	for _, arg := range args {             // 遍历 args 切片
		res = append(res, zap.Any(arg.Key, arg.Val)) // 将每个 Field 转换为 zap.Field 并添加到 res 切片中
	}
	return res // 返回转换后的 zap.Field 切片
}
