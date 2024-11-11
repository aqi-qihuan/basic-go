package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/domain"
	"gitee.com/geekbang/basic-go/lanmengbook/internal/repository/dao"
)

// 定义错误常量
var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

// 定义UserRepository结构体
type UserRepository struct {
	dao *dao.UserDAO
}

// 创建UserRepository实例
func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

// 创建用户
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	// 将domain.User转换为dao.User
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

// 根据邮箱查找用户
func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// 调用dao层查找用户
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		// 如果找不到用户，返回错误
		return domain.User{}, err
	}
	// 将dao.User转换为domain.User
	return repo.toDomain(u), nil
}

// 将dao.User转换为domain.User
func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}
