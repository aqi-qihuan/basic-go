package service

import (
	"basic-go/lanmengbook/internal/domain"
	"basic-go/lanmengbook/internal/repository"
	repomocks "basic-go/lanmengbook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// TestPasswordEncrypt 测试密码加密功能
func TestPasswordEncrypt(t *testing.T) {
	// 原始密码，包含明文密码和额外的字符串
	password := []byte("123456#hello")
	// 使用bcrypt生成加密后的密码，bcrypt.DefaultCost是默认的加密成本
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	// 断言加密过程中没有错误
	assert.NoError(t, err)
	// 打印加密后的密码
	println(string(encrypted))
	// 比较加密后的密码和原始密码，验证加密是否正确
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("123456#hello"))
	// 断言比较过程中没有错误
	assert.NoError(t, err)
}

// Test_userService_Login 测试用户登录功能
func Test_userService_Login(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		// 测试用例名称
		name string

		// mock函数，用于生成模拟的用户仓库
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 预期输入
		ctx      context.Context
		email    string
		password string

		// 预期输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				// 创建一个新的模拟用户仓库
				repo := repomocks.NewMockUserRepository(ctrl)
				// 设置预期调用，查找邮箱为"123@qq.com"的用户
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					// 返回预期的用户信息，包括加密后的密码
					Return(domain.User{
						Email: "123@qq.com",
						// 你在这边拿到的密码，就应该是一个正确的密码
						// 加密后的正确的密码
						Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
						Phone:    "15212345678",
					}, nil)
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",

			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
				Phone:    "15212345678",
			},
		},

		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 设置预期调用，查找邮箱为"123@qq.com"的用户，返回用户未找到错误
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",
			wantErr:  ErrInvalidUserOrPassword,
		},

		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 设置预期调用，查找邮箱为"123@qq.com"的用户，返回数据库错误
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db错误"))
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",
			wantErr:  errors.New("db错误"),
		},

		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				// 设置预期调用，查找邮箱为"123@qq.com"的用户，返回用户信息
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						// 你在这边拿到的密码，就应该是一个正确的密码
						// 加密后的正确的密码
						Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
						Phone:    "15212345678",
					}, nil)
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#helloABCde",

			wantErr: ErrInvalidUserOrPassword,
		},
	}
	// 遍历所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个新的gomock控制器
			ctrl := gomock.NewController(t)
			// 确保在测试结束时调用Finish方法
			defer ctrl.Finish()
			// 调用测试用例的mock函数创建模拟仓库
			repo := tc.mock(ctrl)
			// 创建一个新的用户服务实例
			svc := NewUserService(repo)
			// 调用用户服务的Login方法进行登录操作
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			// 断言实际错误与预期错误是否一致
			assert.Equal(t, tc.wantErr, err)
			// 断言实际用户与预期用户是否一致
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
