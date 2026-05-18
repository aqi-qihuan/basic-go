package ioc

import (
	events2 "basic-go/lmbook/interactive/events"
	"basic-go/lmbook/interactive/repository/dao"
	"basic-go/lmbook/pkg/logger"
	"basic-go/lmbook/pkg/migrator/events"
	"basic-go/lmbook/pkg/migrator/events/fixer"
	"basic-go/lmbook/pkg/saramax"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	return p
}

func InitConsumers(c1 *events2.MySQLBinlogConsumer[dao.Interactive], fixConsumer *fixer.Consumer[dao.Interactive]) []saramax.Consumer {
	return []saramax.Consumer{c1, fixConsumer}
}

// InitInteractiveReadEventConsumer 创建 InteractiveReadEventConsumer
func InitInteractiveReadEventConsumer(
	client sarama.Client,
	l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	producer events.Producer,
) *events2.MySQLBinlogConsumer[dao.Interactive] {
	return events2.NewMySQLBinlogConsumer[dao.Interactive](
		client, l, "interactives", src, dst, producer,
	)
}
