package events

import (
	"basic-go/lmbook/pkg/logger"
	"basic-go/lmbook/search/service"
	"github.com/IBM/sarama"
)

type SyncDataEvent struct {
	IndexName string
	DocID     string
	// 这里应该是 BizTags
	Data string
}

type SyncDataEventConsumer struct {
	logger   logger.LoggerV1
	syncSvc  service.SyncService
	consumer sarama.Consumer
}
