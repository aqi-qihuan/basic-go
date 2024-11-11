package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 定义错误常量
var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// 定义UserDAO结构体
type UserDAO struct {
	db *gorm.DB
}

// 创建UserDAO实例
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// 插入用户信息
func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 获取当前时间戳
	now := time.Now().UnixMilli()
	// 设置创建时间和更新时间
	u.Ctime = now
	u.Utime = now
	// 插入用户信息
	err := dao.db.WithContext(ctx).Create(&u).Error
	// 判断是否为邮箱冲突错误
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

// 根据邮箱查找用户信息
func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	// 根据邮箱查找用户信息
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

// 定义User结构体
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// 时区，UTC 0 的毫秒数
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64

	// json 存储
	//Addr string
}

//type Address struct {
//	Uid
//}
