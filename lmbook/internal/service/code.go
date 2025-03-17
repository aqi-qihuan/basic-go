package service

import (
	"basic-go/lmbook/internal/repository"
	"basic-go/lmbook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

// 定义一个变量 ErrCodeSendTooMany，并将其赋值为 repository 包中的 ErrCodeSendTooMany 错误码
var ErrCodeSendTooMany = repository.ErrCodeSendTooMany

// CodeService 是一个接口，定义了发送验证码和验证输入的验证码的方法。
type CodeService interface {
	// Send 方法用于发送验证码。
	// 参数 ctx 是上下文，用于控制请求的取消和超时。
	// 参数 biz 是业务类型，用于区分不同的验证码业务。
	// 参数 phone 是接收验证码的手机号码。
	// 返回值 error 表示发送过程中可能出现的错误。
	Send(ctx context.Context, biz, phone string) error
	// Verify 方法用于验证输入的验证码。
	// 参数 ctx 是上下文，用于控制请求的取消和超时。
	// 参数 biz 是业务类型，用于区分不同的验证码业务。
	// 参数 phone 是接收验证码的手机号码。
	// 参数 inputCode 是用户输入的验证码。
	// 返回值 bool 表示验证结果，true 表示验证成功，false 表示验证失败。
	// 返回值 error 表示验证过程中可能出现的错误。
	Verify(ctx context.Context,
		biz, phone, inputCode string) (bool, error)
}

// 定义一个名为codeService的结构体
type codeService struct {
	// repo字段，类型为repository.CodeRepository，用于存储代码仓库的接口
	repo repository.CodeRepository
	// sms字段，类型为sms.Service，用于存储短信服务的接口
	sms sms.Service
}

// NewCodeService 是一个工厂函数，用于创建一个新的 CodeService 实例。
// 它接收两个参数：一个是代码仓库接口 repository.CodeRepository 的实例 repo，
// 另一个是短信服务接口 sms.Service 的实例 smsSvc。
// 返回值是一个 CodeService 接口的实例。
func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	// 创建一个 codeService 结构体的实例，并将传入的 repo 和 smsSvc 赋值给结构体的相应字段。
	// 这里使用 & 符号获取结构体的指针，因为 CodeService 接口可能需要指针接收者。
	return &codeService{
		repo: repo,   // 将传入的代码仓库实例赋值给 codeService 结构体的 repo 字段。
		sms:  smsSvc, // 将传入的短信服务实例赋值给 codeService 结构体的 sms 字段。
	}
}

// Send 是 codeService 结构体的一个方法，用于发送验证码
func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	// 调用 generate 方法生成验证码
	code := svc.generate()
	// 调用 repo 的 Set 方法将生成的验证码存储到数据库或缓存中
	// 参数 ctx 是上下文，用于控制请求的取消和超时
	// 参数 biz 是业务类型，用于区分不同的业务场景
	// 参数 phone 是接收验证码的手机号
	// 参数 code 是生成的验证码
	err := svc.repo.Set(ctx, biz, phone, code)
	// 你在这儿，是不是要开始发送验证码了？
	if err != nil {
		return err
	}
	const codeTplId = "1877556"
	return svc.sms.Send(ctx, codeTplId, []string{code}, phone)
}

// Verify 是 codeService 结构体的一个方法，用于验证输入的验证码是否正确。
// 参数 ctx 是上下文，用于控制请求的截止时间、取消信号等。
// 参数 biz 是业务类型，用于区分不同的验证码业务。
// 参数 phone 是手机号，用于标识验证码的接收者。
// 参数 inputCode 是用户输入的验证码，需要与系统生成的验证码进行比对。
// 返回值 bool 表示验证结果，true 表示验证通过，false 表示验证失败。
// 返回值 error 表示验证过程中发生的错误，如果验证通过且无错误，则返回 nil。
func (svc *codeService) Verify(ctx context.Context,
	biz, phone, inputCode string) (bool, error) {
	// 调用 svc.repo.Verify 方法进行验证码的验证，传入上下文、业务类型、手机号和用户输入的验证码。
	// 该方法返回两个值：ok 表示验证结果，err 表示验证过程中可能发生的错误。
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	// 检查验证过程中是否发生了特定的错误 repository.ErrCodeVerifyTooMany，
	// 该错误表示验证次数过多，可能是因为用户多次输入错误的验证码。
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *codeService) generate() string {
	// 0-999999
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
