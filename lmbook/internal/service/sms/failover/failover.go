package failover

import (
	"basic-go/lmbook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

// FailOverSMSService 是一个短信服务失败转移的服务结构体
type FailOverSMSService struct {
	svcs []sms.Service // 存储多个短信服务实例，用于在发送短信时进行失败转移
	// 　v1 的字段
	// 当前服务商下标，使用原子操作保证并发安全，用于循环选择不同的短信服务
	idx uint64
}

// NewFailOverSMSService 创建一个新的 FailOverSMSService 实例
// 参数 svcs 是一个包含多个短信服务实例的切片
func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

// Send 方法尝试依次使用所有的短信服务发送短信
// 只要有一个服务发送成功，就返回 nil；如果所有服务都失败，则返回错误
func (f *FailOverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 遍历所有的短信服务实例
	for _, svc := range f.svcs {
		// 调用当前短信服务的 Send 方法
		err := svc.Send(ctx, tplId, args, numbers...)
		// 如果发送成功，返回 nil
		if err == nil {
			return nil
		}
		// 记录发送失败的错误信息
		log.Println(err)
	}
	// 所有服务都发送失败，返回错误
	return errors.New("轮询了所有服务服务商，但是都发送失败了")
}

// 起始下标轮询
// 并且出错也轮询
// SendV1 方法使用原子操作更新当前服务商下标，从该下标开始轮询短信服务
// 处理特定错误（如 context.Canceled 和 context.DeadlineExceeded），并在所有服务都失败时返回错误
func (f *FailOverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 原子操作增加当前服务商下标
	idx := atomic.AddUint64(&f.idx, 1)
	// 获取短信服务实例的数量
	length := uint64(len(f.svcs))
	// 迭代所有的短信服务
	for i := 0; i < int(length); i++ {
		// 取余数来计算当前要使用的短信服务下标
		svc := f.svcs[idx%length]
		// 调用当前短信服务的 Send 方法
		err := svc.Send(ctx, tplId, args, numbers...)
		// 根据不同的错误情况进行处理
		switch err {
		case nil:
			// 发送成功，返回 nil
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 前者是被取消，后者超时，直接返回错误
			return err
		}
		// 记录发送失败的错误信息
		log.Println(err)
	}
	// 所有服务都发送失败，返回错误
	return errors.New("轮询了所有服务服务商，但是都发送失败了")
}
