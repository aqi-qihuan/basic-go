package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/service/sms"
	"math/rand"
)

var ErrCodeSendTooMany = repository.ErrCodeVerifyTooMany

// CodeService 定义了一个短信验证码服务接口
type CodeService interface {
	// Send 方法用于发送短信验证码
	Send(ctx context.Context, biz, phone string) error
	// Verify 方法用于验证短信验证码
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}

// codeService 是 CodeService 接口的实现
type codeService struct {
	// repo 是验证码存储库
	repo repository.CodeRepository
	// sms 是短信服务
	sms sms.Service
}

// NewCodeService 创建一个新的 CodeService 实例
func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo: repo,
		sms:  smsSvc,
	}
}

// Send 方法用于发送短信验证码
func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	// 生成一个随机的 6 位数字验证码
	code := svc.generate()
	// 将验证码存储到存储库中
	err := svc.repo.Set(ctx, biz, phone, code)
	// 如果存储失败，返回错误
	if err != nil {
		return err
	}
	// 发送短信验证码
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

// Verify 方法用于验证短信验证码
func (svc *codeService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	// 验证短信验证码
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	// 如果验证次数过多，返回错误
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	// 返回验证结果
	return ok, err
}

// generate 方法用于生成一个随机的 6 位数字验证码
func (svc *codeService) generate() string {
	// 生成一个 0-999999 的随机数
	code := rand.Intn(1000000)
	// 将随机数格式化为 6 位数字字符串
	return fmt.Sprintf("%06d", code)
}
