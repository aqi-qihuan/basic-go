//go:build wireinject

package startup

import (
	"basic-go/lmbook/interactive/grpc"
	"basic-go/lmbook/interactive/repository"
	"basic-go/lmbook/interactive/repository/cache"
	"basic-go/lmbook/interactive/repository/dao"
	"basic-go/lmbook/interactive/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	InitRedis, InitTestDB,
	InitLog,
	InitKafka,
)

func InitGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		grpc.NewInteractiveServiceServer,
		thirdProvider,
		dao.NewGORMInteractiveDAO,
		cache.NewRedisInteractiveCache,
		repository.NewCachedInteractiveRepository,
		service.NewInteractiveService,
	)
	return new(grpc.InteractiveServiceServer)
}
