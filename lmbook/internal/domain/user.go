package domain

import "time"

// 定义一个用户结构体
type User struct {
	Id       int64
	Email    string
	Password string

	Nickname string
	// YYYY-MM-DD
	Birthday time.Time
	AboutMe  string

	Phone string

	// UTC 0 的时区
	Ctime time.Time

	WechatInfo WechatInfo
	//Addr Address
}

// TodeayIsBirthday 判断今天是否是生日
func (u User) TodayIsBirthday() bool {
	now := time.Now()
	return u.Birthday.Year() == now.Year() &&
		u.Birthday.Month() == now.Month() &&
		u.Birthday.Day() == now.Day()
}

//type Address struct {
//	Province string
//	Region   string
//}

//func (u User) ValidateEmail() bool {
// 在这里用正则表达式校验
//return u.Email
//}
