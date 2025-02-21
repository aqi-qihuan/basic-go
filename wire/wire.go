//go:build wireinject

package wire

import (
	"basic-go/wire/repository"
	"basic-go/wire/repository/dao"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, InitDB, dao.NewUserDAO, InitDB)
	return &repository.UserRepository{}
}
