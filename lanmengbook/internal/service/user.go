package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/domain"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository"
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
}

// 定义UserService结构体
type userService struct {
	// repo 是用户仓库的实例
	repo repository.UserRepository
}

// NewUserService 函数创建一个新的 UserService 实例
func NewUserService(repo repository.UserRepository) UserService {
	// 返回一个 userService 结构体实例，该实例实现了 UserService 接口
	return &userService{
		// 初始化 userService 实例的 repo 字段，该字段是一个 UserRepository 接口的实现
		repo: repo,
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
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err!= nil，系统错误
	if err != nil && !errors.Is(err, repository.ErrDuplicateUser) {
		return domain.User{}, err

	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)

}
