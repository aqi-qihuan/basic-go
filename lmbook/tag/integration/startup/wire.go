//go:build wireinject

package startup

import (
	"basic-go/lmbook/tag/events"
	"basic-go/lmbook/tag/grpc"
	"basic-go/lmbook/tag/repository/cache"
	"basic-go/lmbook/tag/repository/dao"
	"basic-go/lmbook/tag/service"
	"github.com/google/wire"
)

func InitGRPCService(p events.Producer) *grpc.TagServiceServer {
	wire.Build(InitTestDB, InitRedis,
		InitLog,
		dao.NewGORMTagDAO,
		InitRepository,
		cache.NewRedisTagCache,
		service.NewTagService,
		grpc.NewTagServiceServer,
	)
	return new(grpc.TagServiceServer)
}
