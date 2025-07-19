//go:build wireinject

package startup

import (
	"basic-go/lmbook/follow/grpc"
	"basic-go/lmbook/follow/repository"
	"basic-go/lmbook/follow/repository/cache"
	"basic-go/lmbook/follow/repository/dao"
	"basic-go/lmbook/follow/service"
	"github.com/google/wire"
)

func InitServer() *grpc.FollowServiceServer {
	wire.Build(
		InitRedis,
		InitLog,
		InitTestDB,
		dao.NewGORMFollowRelationDAO,
		cache.NewRedisFollowCache,
		repository.NewFollowRelationRepository,
		service.NewFollowRelationService,
		grpc.NewFollowRelationServiceServer,
	)
	return new(grpc.FollowServiceServer)
}
