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
	// ErrDuplicateEmail 表示邮箱已存在
	ErrDuplicateEmail = repository.ErrDuplicateEmail
	// ErrInvalidUserOrPassword 表示用户不存在或者密码不对
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

// 定义UserService结构体
type UserService struct {
	// repo 是用户仓库的实例
	repo *repository.UserRepository
}

// 创建UserService实例
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// 注册用户
func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	// 对密码进行加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 创建用户
	return svc.repo.Create(ctx, u)
}

// 用户登录
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
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
