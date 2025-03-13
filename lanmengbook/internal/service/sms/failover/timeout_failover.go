package failover

import (
	"basic-go/lanmengbook/internal/service/sms"
	"context"
	"sync/atomic"
)

// TimeoutFailoverSMSService 是一个带有超时切换机制的短信服务
type TimeoutFailoverSMSService struct {
	svcs []sms.Service // 存储多个短信服务实例
	// 当前正在使用节点
	idx int32
	// 连续几个超时了
	cnt int32
	// 切换的间隔 只读的
	threshold int32
}

// NewTimeoutFailoverSMSService 创建一个新的 TimeoutFailoverSMSService 实例
func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,      // 初始化短信服务列表
		threshold: threshold, // 初始化切换阈值
	}
}

// Send 发送短信，并在超时时切换到下一个短信服务
func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx) // 原子加载当前使用的短信服务索引
	cnt := atomic.LoadInt32(&t.cnt) // 原子加载连续超时计数
	// 超过阈值，执行切换
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))             // 计算新的索引，循环切换
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) { // 原子比较并交换索引
			// 重置这个 cnt 计数
			atomic.StoreInt32(&t.cnt, 0) // 原子存储重置计数
		}
		idx = newIdx // 更新当前索引
	}
	svc := t.svcs[idx]                            // 获取当前短信服务实例
	err := svc.Send(ctx, tplId, args, numbers...) // 调用当前短信服务发送短信
	switch err {
	case nil:
		//连续超时，所以不超时的时候重置到0
		atomic.StoreInt32(&t.cnt, 0) // 原子存储重置计数
		return nil
	case context.DeadlineExceeded:
		// 超时了，cnt + 1
		atomic.AddInt32(&t.cnt, 1) // 原子增加计数
	default:
		// 遇到了错误，但是又不是超时错误，这个时候，你要考虑怎么搞
		// 我可以增加，也可以不增加
		// 如果强调一定是超时，那么就不增加
		// 如果是 EOF 之类的错误，你还可以考虑直接切换
	}
	return err // 返回错误
}
