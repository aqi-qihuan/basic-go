//go:build wireinject

package startup

import (
	"basic-go/lmbook/interactive/grpc"
	repository2 "basic-go/lmbook/interactive/repository"
	cache2 "basic-go/lmbook/interactive/repository/cache"
	dao2 "basic-go/lmbook/interactive/repository/dao"
	service2 "basic-go/lmbook/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB,
	//InitSaramaClient,
	//InitSyncProducer,
	InitLogger,
)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitInteractiveService() *grpc.InteractiveServiceServer {
	wire.Build(thirdPartySet, interactiveSvcSet, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
