package simpleim

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/IBM/sarama"
)

type IMService struct {
	producer sarama.SyncProducer
}

func (s *IMService) Receive(ctx context.Context, sender int64, msg Message) error {
	// 转发到 Kafka 里面，以通知别的网关节点
	// 1. 查找目标
	members := s.findMembers()
	// 2. 通知 Kafka，让别的节点能够订阅到消息
	for _, mem := range members {
		if mem == sender {
			// 本人就不用转发了
			continue
		}
		msgJson, err := json.Marshal(Event{Receiver: mem, Msg: msg})
		if err != nil {
			continue
		}
		_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
			Topic: eventName,
			Key:   sarama.ByteEncoder(strconv.FormatInt(mem, 10)),
			Value: sarama.ByteEncoder(msgJson),
		})
		if err != nil {
			log.Println("发送消息失败", err)
			continue
		}
	}
	return nil
}

func (s *IMService) findMembers() []int64 {
	// 这里就是查询 IM 中的群组服务拿到成员
	// 模拟拿到的结果（模拟数据库查询的结果）
	return []int64{1, 2, 3, 4}
}
