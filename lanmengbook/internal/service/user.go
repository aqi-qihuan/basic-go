package service

import (
	"basic-go/lanmengbook/internal/domain"
	"basic-go/lanmengbook/internal/repository"
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 定义错误常量
var (
	// ErrDuplicateEmail 表示用户邮箱重复错误
	ErrDuplicateEmail = repository.ErrDuplicateUser
	// ErrInvalidUserOrPassword 表示用户不存在或者密码不对
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

// UserService 定义了一个用户服务接口
type UserService interface {
	// Signup 方法用于注册用户
	Signup(ctx context.Context, u domain.User) error
	// Login 方法用于用户登录
	Login(ctx context.Context, email string, password string) (domain.User, error)
	// UpdateNonSensitiveInfo 方法用于更新用户的非敏感信息
	UpdateNonSensitiveInfo(ctx context.Context,
		user domain.User) error
	// FindById 方法用于根据用户ID查找用户
	FindById(ctx context.Context,
		uid int64) (domain.User, error)
	// FindOrCreate 方法用于查找或创建用户
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	// FindOrCreateByWechat 是一个方法，用于根据微信信息查找或创建用户。
	// ctx 参数是上下文对象，用于控制请求的截止时间、取消信号和其他请求范围的值。
	// info 参数是 domain.WechatInfo 类型的对象，包含了微信用户的信息。
	// 返回值是一个 domain.User 对象，表示找到或创建的用户，以及一个 error 对象，表示操作过程中可能发生的错误。
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

// 定义一个userService结构体，用于处理用户相关的业务逻辑
type userService struct {
	// repo字段用于存储用户数据访问层的接口实例
	repo repository.UserRepository
	// logger字段用于记录日志，目前被注释掉了
	//logger *zap.Logger
}

// NewUserService函数用于创建一个新的userService实例
// 参数repo是用户数据访问层的接口实例
// 返回值是UserService接口的实现，即userService结构体的指针
func NewUserService(repo repository.UserRepository) UserService {
	// 返回一个新的userService实例，并将传入的repo赋值给userService的repo字段
	// logger字段目前被注释掉了，没有赋值
	return &userService{
		repo: repo,
		//logger: zap.L(),
	}
}

// Signup 方法，用于处理用户注册的 HTTP 请求
func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		// 如果加密失败，则返回错误
		return err
	}
	// 将加密后的密码存储在用户对象中
	u.Password = string(hash)
	// 创建用户
	return svc.repo.Create(ctx, u)
}

// 用户登录
func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 根据邮箱查找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		// 如果用户不存在，返回错误
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		// 如果发生其他错误，返回错误
		return domain.User{}, err
	}
	// 检查密码对不对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 如果密码不对，返回错误
		return domain.User{}, ErrInvalidUserOrPassword
	}
	// 返回用户信息
	return u, nil
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context,
	user domain.User) error {
	// UpdateNicknameAndXXAnd
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

// FindById 方法，用于根据用户ID查找用户
func (svc *userService) FindById(ctx context.Context,
	uid int64) (domain.User, error) {
	// 调用 repo 的 FindById 方法来查找用户
	return svc.repo.FindById(ctx, uid)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User,
	error) {
	// 先找一下，我们认为，大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err == repository.ErrUserNotFound {
		// 有两种情况
		// err == nil, u 是可用的
		// err!= nil，系统错误，
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err!= nil，系统错误
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)

}
func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 这边就是意味着是一个新用户
	// JSON 格式的 wechatInfo
	zap.L().Info("新用户", zap.Any("wechatInfo", wechatInfo))
	//svc.logger.Info("新用户", zap.Any("wechatInfo", wechatInfo))
	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: wechatInfo,
	})
	if err != nil && err != repository.ErrDuplicateUser {
		return domain.User{}, err
	}
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
}
