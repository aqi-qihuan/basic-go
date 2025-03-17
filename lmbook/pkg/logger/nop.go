package logger

// 定义一个名为NopLogger的结构体，用于实现一个无操作的日志记录器
type NopLogger struct {
}

// NewNopLogger函数用于创建并返回一个新的NopLogger实例
func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

// Debug方法用于记录调试级别的日志信息，接收一个字符串消息和可变数量的Field参数
// 由于是NopLogger，该方法不执行任何操作
func (n *NopLogger) Debug(msg string, args ...Field) {
}

// Info方法用于记录信息级别的日志信息，接收一个字符串消息和可变数量的Field参数
// 由于是NopLogger，该方法不执行任何操作
func (n *NopLogger) Info(msg string, args ...Field) {
}

// Warn方法用于记录警告级别的日志信息，接收一个字符串消息和可变数量的Field参数
// 由于是NopLogger，该方法不执行任何操作
func (n *NopLogger) Warn(msg string, args ...Field) {
}

// Error方法用于记录错误级别的日志信息，接收一个字符串消息和可变数量的Field参数
// 由于是NopLogger，该方法不执行任何操作
func (n *NopLogger) Error(msg string, args ...Field) {
}
