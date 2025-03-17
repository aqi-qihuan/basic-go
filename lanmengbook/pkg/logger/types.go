package logger

// 定义一个 Logger 接口，包含四种日志级别的方法：Debug, Info, Warn, Error
// 这些方法都接受一个字符串消息和一个可变参数列表，参数列表可以是任意类型
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// example 函数演示了如何使用 Logger 接口
// 这里定义了一个 Logger 类型的变量 l，并调用其 Info 方法记录一条信息
func example() {
	var l Logger
	// 调用 Info 方法，记录一条信息，其中包含一个整数参数 123
	l.Info("用户的微信 id %d", 123)
}

// 定义一个 LoggerV1 接口，与 Logger 接口类似，但参数列表中的字段必须是 Field 类型
type LoggerV1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}

// Field 结构体用于表示键值对，其中 Key 是字符串类型的键，Val 是任意类型的值
type Field struct {
	Key string
	Val any
}

// exampleV1 函数演示了如何使用 LoggerV1 接口
// 这里定义了一个 LoggerV1 类型的变量 l，并调用其 Info 方法记录一条信息
func exampleV1() {
	var l LoggerV1
	// 调用 Info 方法，记录一条信息，其中包含一个 Field 类型的参数，表示新用户的 union_id
	// 这是一个新用户 union_id=123
	l.Info("这是一个新用户", Field{Key: "union_id", Val: 123})
}

// 定义一个 LoggerV2 接口，与 Logger 接口类似，但参数列表中的字段必须是偶数个
// 参数列表必须以 key1,value1,key2,value2 的形式传递
type LoggerV2 interface {
	// 它要去 args 必须是偶数，并且是以 key1,value1,key2,value2 的形式传递
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
