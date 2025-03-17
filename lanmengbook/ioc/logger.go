package ioc

import (
	"basic-go/lanmengbook/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitLogger 初始化日志记录器
func InitLogger() logger.LoggerV1 {
	// 创建一个新的开发环境配置
	cfg := zap.NewDevelopmentConfig()
	// 使用Viper库从配置文件中解析"log"键的值，并将其反序列化到cfg中
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		// 如果解析配置时发生错误，则抛出panic
		panic(err)
	}
	// 使用配置构建一个新的日志记录器
	l, err := cfg.Build()
	if err != nil {
		// 如果构建日志记录器时发生错误，则抛出panic
		panic(err)
	}
	// 返回一个新的Zap日志记录器实例
	return logger.NewZapLogger(l)
}
