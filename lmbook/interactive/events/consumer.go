package events

import (
	"basic-go/lmbook/pkg/canalx"
	"basic-go/lmbook/pkg/logger"
	"basic-go/lmbook/pkg/migrator"
	"basic-go/lmbook/pkg/migrator/events"
	"basic-go/lmbook/pkg/migrator/validator"
	"basic-go/lmbook/pkg/saramax"
	"context"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

type MySQLBinlogConsumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.LoggerV1
	table    string
	srcToDst *validator.CanalIncrValidator[T]
	dstToSrc *validator.CanalIncrValidator[T]
	dstFirst *atomic.Bool
}

func NewMySQLBinlogConsumer[T migrator.Entity](
	client sarama.Client,
	l logger.LoggerV1,
	table string,
	src *gorm.DB,
	dst *gorm.DB,
	p events.Producer) *MySQLBinlogConsumer[T] {
	srcToDst := validator.NewCanalIncrValidator[T](src, dst, "SRC", l, p)
	dstToSrc := validator.NewCanalIncrValidator[T](src, dst, "DST", l, p)
	return &MySQLBinlogConsumer[T]{
		client: client, l: l,
		dstFirst: &atomic.Bool{},
		srcToDst: srcToDst,
		dstToSrc: dstToSrc,
		table:    table}
}

func (r *MySQLBinlogConsumer[T]) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("migrator_incr",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"lmbook_binlog"},
			saramax.NewHandler[canalx.Message[T]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer[T]) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[T]) error {
	// 首先判定这个消息我要不要处理

	// 以源表为准还是以目标表为准
	dstFirst := r.dstFirst.Load()
	var v *validator.CanalIncrValidator[T]
	// db:
	//  src:
	//    dsn: "root:root@tcp(localhost:13316)/lmbook"
	//  dst:
	//    dsn: "root:root@tcp(localhost:13316)/lmbook_intr"
	// 出于保险，你可以进一步校验表名
	if dstFirst && val.Database == "lmbook_intr" {
		// 以目标表为准，过来的也恰好是目标表的 binlog
		// 要校验
		v = r.dstToSrc
	} else if !dstFirst && val.Database == "lmbook" {
		// 以源表为准，过来的也恰好是源表的 binlog
		// 要校验
		v = r.srcToDst
	}

	if v == nil {
		return nil
	}

	for _, data := range val.Data {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := v.Validate(ctx, data.ID())
		cancel()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MySQLBinlogConsumer[T]) DstFirst() {
	r.dstFirst.Store(true)
}
