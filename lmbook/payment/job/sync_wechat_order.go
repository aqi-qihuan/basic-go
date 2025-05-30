package job

import (
	"basic-go/lmbook/payment/service/wechat"
	"basic-go/lmbook/pkg/logger"
	"context"
	"time"
)

type SyncWechatOrderJob struct {
	svc *wechat.NativePaymentService
	l   logger.LoggerV1
}

func (s *SyncWechatOrderJob) Name() string {
	return "sync_wechat_order_job"
}

// 我这个定时任务，多久运行一次？
// 不必特别频繁，比如说一分钟运行一次
func (s *SyncWechatOrderJob) Run() error {
	// 定时找到超时的微信支付订单，然后发起同步
	// 针对过期订单
	t := time.Now().Add(-time.Minute * 31)
	//t := time.Now().Add(-time.Minute * 5)
	offset := 0
	const limit = 100
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		pmts, err := s.svc.FindExpiredPayment(ctx, offset, limit, t)
		cancel()
		if err != nil {
			// 如果不中断
			return err
		}
		for _, pmt := range pmts {
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
			err = s.svc.SyncWechatInfo(ctx, pmt.BizTradeNO)
			cancel()
			if err != nil {
				s.l.Error("同步微信订单状态失败", logger.Error(err),
					logger.String("biz_trade_no", pmt.BizTradeNO))
			}
		}
		if len(pmts) < limit {
			return nil
		}
		offset = offset + len(pmts)
	}
}
