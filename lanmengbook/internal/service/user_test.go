package service

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// 测试密码加密
func TestUserRegister(t *testing.T) {
	// 定义密码
	password := []byte("123456#hello")
	// 使用bcrypt库对密码进行加密
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	// 断言没有错误
	assert.NoError(t, err)
	// 打印加密后的密码
	println(string(encrypted))
	// 使用bcrypt库对加密后的密码进行解密
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#hello"))
	// 断言没有错误
	assert.NoError(t, err)
}
