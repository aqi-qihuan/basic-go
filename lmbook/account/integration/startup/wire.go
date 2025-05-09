//go:build wireinject

package startup

import (
	"basic-go/lmbook/account/grpc"
	"basic-go/lmbook/account/repository"
	"basic-go/lmbook/account/repository/dao"
	"basic-go/lmbook/account/service"
	"github.com/google/wire"
)

func InitAccountService() *grpc.AccountServiceServer {
	wire.Build(InitTestDB,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer)
	return new(grpc.AccountServiceServer)
}
