//go:build wireinject

package startup

import (
	pmtv1 "basic-go/lmbook/api/proto/gen/payment/v1"
	"basic-go/lmbook/reward/repository"
	"basic-go/lmbook/reward/repository/cache"
	"basic-go/lmbook/reward/repository/dao"
	"basic-go/lmbook/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(InitTestDB, InitLogger, InitRedis)

func InitWechatNativeSvc(client pmtv1.WechatPaymentServiceClient) *service.WechatNativeRewardService {
	wire.Build(service.NewWechatNativeRewardService,
		thirdPartySet,
		cache.NewRewardRedisCache,
		repository.NewRewardRepository, dao.NewRewardGORMDAO)
	return new(service.WechatNativeRewardService)
}
