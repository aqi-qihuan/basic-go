package ioc

import (
	"basic-go/lmbook/interactive/repository/dao"
	"basic-go/lmbook/pkg/ginx"
	"basic-go/lmbook/pkg/gormx/connpool"
	"basic-go/lmbook/pkg/logger"
	"basic-go/lmbook/pkg/migrator/events"
	"basic-go/lmbook/pkg/migrator/events/fixer"
	"basic-go/lmbook/pkg/migrator/scheduler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

// InitGinxServer 管理后台的 server
func InitGinxServer(l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	pool *connpool.DoubleWritePool,
	producer events.Producer) *ginx.Server {
	engine := gin.Default()
	group := engine.Group("/migrator")
	ginx.InitCounter(prometheus2.CounterOpts{
		Namespace: "geektime_daming",
		Subsystem: "lmbook_intr_admin",
		Name:      "biz_code",
		Help:      "统计业务错误码",
	})
	sch := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	sch.RegisterRoutes(group)
	return &ginx.Server{
		Engine: engine,
		Addr:   viper.GetString("migrator.http.addr"),
	}
}

func InitInteractiveProducer(p sarama.SyncProducer) events.Producer {
	return events.NewSaramaProducer("inconsistent_interactive", p)
}

func InitFixerConsumer(client sarama.Client,
	l logger.LoggerV1,
	src SrcDB,
	dst DstDB) *fixer.Consumer[dao.Interactive] {
	res, err := fixer.NewConsumer[dao.Interactive](client, l, "inconsistent_interactive", src, dst)
	if err != nil {
		panic(err)
	}
	return res
}
