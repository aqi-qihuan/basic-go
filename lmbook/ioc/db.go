package ioc

import (
	"basic-go/lmbook/internal/repository/dao"
	"basic-go/lmbook/pkg/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
// 参数 l 是一个符合 logger.LoggerV1 接口的日志记录器
// 返回一个 *gorm.DB 类型的数据库连接实例
func InitDB(l logger.LoggerV1) *gorm.DB {
	// 定义一个 Config 结构体，用于存储数据库配置信息
	type Config struct {
		DSN string `yaml:"dsn"` // DSN 数据源名称，用于连接数据库
	}
	// 初始化一个 Config 实例，并设置默认的 DSN
	var cfg Config = Config{
		DSN: "root:root@tcp(localhost:3316)/webook",
	}
	// 使用 viper 解析配置文件中的 "db" 键，并将值填充到 cfg 中
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		// 如果解析配置文件出错，则抛出 panic
		panic(err)
	}
	// 使用 gorm 打开一个数据库连接
	// mysql.Open(cfg.DSN) 用于根据 DSN 创建 MySQL 的数据库连接
	// &gorm.Config{} 用于配置 gorm 的行为
	// Logger: glogger.New(goormLoggerFunc(l.Debug), glogger.Config{}) 用于配置 gorm 的日志记录器
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: glogger.New(goormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值设置为 0，表示不记录慢查询
			SlowThreshold: 0,
			// 日志级别设置为 Info，表示记录一般信息
			LogLevel: glogger.Info,
		}),
	})
	if err != nil {
		// 如果打开数据库连接出错，则抛出 panic
		panic(err)
	}
	// 初始化数据库表结构
	err = dao.InitTables(db)
	if err != nil {
		// 如果初始化表结构出错，则抛出 panic
		panic(err)
	}
	// 返回数据库连接实例
	return db
}

// goormLoggerFunc 是一个函数类型，用于适配 gorm 的日志记录器
// 参数 msg 是日志消息
// 参数 fields 是日志的附加字段
type goormLoggerFunc func(msg string, fields ...logger.Field)

// Printf 实现了 gorm.Logger 接口的 Printf 方法
// 参数 s 是格式化字符串
// 参数 i 是格式化字符串的参数
func (g goormLoggerFunc) Printf(s string, i ...interface{}) {
	// 调用 goormLoggerFunc 类型的函数，将格式化字符串和参数转换为日志消息
	g(s, logger.Field{Key: "args", Val: i})
}
