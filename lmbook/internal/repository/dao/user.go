package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	// ErrDuplicateEmail 表示邮箱冲突错误
	ErrDuplicateEmail = errors.New("邮箱冲突")
	// ErrRecordNotFound 表示记录未找到错误
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// UserDAO 是用户数据访问对象的接口
type UserDAO interface {
	// Insert 插入一个新用户
	Insert(ctx context.Context, u User) error
	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (User, error)
	// UpdateById 根据用户ID更新用户信息
	UpdateById(ctx context.Context, entity User) error
	// FindById 根据用户ID查找用户
	FindById(ctx context.Context, uid int64) (User, error)
	// FindByPhone 根据电话号码查找用户
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechat(ctx context.Context, openId string) (User, error)
}

// GORMUserDAO 是基于 GORM 的用户数据访问对象实现
type GORMUserDAO struct {
	// db 是 GORM 数据库对象
	db *gorm.DB
}

func (dao *GORMUserDAO) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id=?", openId).First(&u).Error
	return u, err
}

// NewUserDAO 创建一个新的 UserDAO 实例
func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

// Insert 插入一个新用户
func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 用户冲突，邮箱冲突
			return ErrDuplicateEmail
		}
	}
	return err
}

// FindByEmail 根据邮箱查找用户
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

// UpdateById 根据用户ID更新用户信息
func (dao *GORMUserDAO) UpdateById(ctx context.Context, entity User) error {
	// 这种写法依赖于 GORM 的零值和主键更新特性
	// Update 非零值 WHERE id =?
	//return dao.db.WithContext(ctx).Updates(&entity).Error
	return dao.db.WithContext(ctx).Model(&entity).Where("id =?", entity.Id).
		Updates(map[string]any{
			"utime":    time.Now().UnixMilli(),
			"nickname": entity.Nickname,
			"birthday": entity.Birthday,
			"about_me": entity.AboutMe,
		}).Error
}

// FindById 根据用户ID查找用户
func (dao *GORMUserDAO) FindById(ctx context.Context, uid int64) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("id =?", uid).First(&res).Error
	return res, err
}

// FindByPhone 根据电话号码查找用户
func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var res User
	err := dao.db.WithContext(ctx).Where("phone =?", phone).First(&res).Error
	return res, err
}

// User 是用户结构体
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 代表这是一个可以为 NULL 的列
	//Email    *string
	Email    sql.NullString `gorm:"unique"`
	Password string

	Nickname string `gorm:"type=varchar(128)"`
	// YYYY-MM-DD
	Birthday int64
	AboutMe  string `gorm:"type=varchar(4096)"`

	// 代表这是一个可以为 NULL 的列
	Phone sql.NullString `gorm:"unique"`

	// 1 如果查询要求同时使用 openid 和 unionid，就要创建联合唯一索引
	// 2 如果查询只用 openid，那么就在 openid 上创建唯一索引，或者 <openid, unionId> 联合索引
	// 3 如果查询只用 unionid，那么就在 unionid 上创建唯一索引，或者 <unionid, openid> 联合索引

	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString

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
