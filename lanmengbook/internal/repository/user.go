package repository

import (
	"basic-go/lanmengbook/internal/domain"
	"basic-go/lanmengbook/internal/repository/cache"
	"basic-go/lanmengbook/internal/repository/dao"
	"context"
	"database/sql"
	"log"
	"time"
)

var (
	// ErrDuplicateUser 表示用户重复错误
	ErrDuplicateUser = dao.ErrDuplicateEmail
	// ErrUserNotFound 表示用户未找到错误
	ErrUserNotFound = dao.ErrRecordNotFound
)

// UserRepository 是用户存储库的接口
type UserRepository interface {
	// Create 创建一个新用户
	Create(ctx context.Context, u domain.User) error
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	// UpdateNonZeroFields 更新用户的非零字段
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
	// FindByPhone 根据电话号码查找用户
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	// FindById 根据用户ID查找用户
	FindById(ctx context.Context, uid int64) (domain.User, error)
}

// CachedUserRepository 是基于缓存的用户存储库实现
type CachedUserRepository struct {
	// dao 是用户数据访问对象
	dao dao.UserDAO
	// cache 是用户缓存对象
	cache cache.UserCache
}

// NewCachedUserRepository 函数创建一个新的 CachedUserRepository 实例
func NewCachedUserRepository(dao dao.UserDAO,
	c cache.UserCache) UserRepository {
	// 返回一个 CachedUserRepository 结构体实例，该实例实现了 UserRepository 接口
	return &CachedUserRepository{
		// 初始化 CachedUserRepository 实例的 dao 字段，该字段是一个 UserDAO 接口的实现
		dao: dao,
		// 初始化 CachedUserRepository 实例的 cache 字段，该字段是一个 UserCache 接口的实现
		cache: c,
	}
}

// Create 方法用于创建一个新用户
func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

// FindByEmail 方法用于根据邮箱查找用户
func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// 调用 dao 层的 FindByEmail 方法，根据邮箱地址查找用户
	u, err := repo.dao.FindByEmail(ctx, email)
	// 如果查找过程中发生错误，返回一个空的用户对象和该错误
	if err != nil {
		return domain.User{}, err
	}
	// 如果没有错误，将从数据库中查找到的用户数据转换为领域对象，并返回
	return repo.toDomain(u), nil
}

// toDomain 方法用于将数据访问对象转换为领域对象
func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

// toEntity 方法用于将领域对象转换为数据访问对象
func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

// UpdateNonZeroFields 方法用于更新用户的非零字段
func (repo *CachedUserRepository) UpdateNonZeroFields(ctx context.Context,
	user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

// FindById 方法用于根据用户ID查找用户
func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err 为 nil，就返回
	if err == nil {
		return du, nil
	}

	// err 不为 nil，就要查询数据库
	// err 有两种可能
	// 1. key 不存在，说明 redis 是正常的
	// 2. 访问 redis 有问题。可能是网络有问题，也可能是 redis 本身就崩溃了

	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	//go func() {
	//	err = repo.cache.Set(ctx, du)
	//	if err!= nil {
	//		log.Println(err)
	//	}
	//}()

	err = repo.cache.Set(ctx, du)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		log.Println(err)
	}
	return du, nil
}

// FindByIdV1 方法用于根据用户ID查找用户，处理缓存未命中的情况
func (repo *CachedUserRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err 为 nil，就返回
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		u, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		//go func() {
		//	err = repo.cache.Set(ctx, du)
		//	if err!= nil {
		//		log.Println(err)
		//	}
		//}()

		err = repo.cache.Set(ctx, du)
		if err != nil {
			// 网络崩了，也可能是 redis 崩了
			log.Println(err)
		}
		return du, nil
	default:
		// 接近降级的写法
		return domain.User{}, err
	}

}

// FindByPhone 方法用于根据电话号码查找用户
func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	// 调用 dao 层的 FindByPhone 方法，根据电话号码查找用户
	u, err := repo.dao.FindByPhone(ctx, phone)
	// 如果查找过程中发生错误，返回一个空的用户对象和该错误
	if err != nil {
		return domain.User{}, err
	}
	// 如果没有错误，将查找到的用户对象转换为领域对象，并返回
	return repo.toDomain(u), nil
}
