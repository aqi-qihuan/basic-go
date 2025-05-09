package startup

import (
	"basic-go/lmbook/payment/ioc"
	"basic-go/lmbook/payment/repository"
	"basic-go/lmbook/payment/repository/dao"
	"basic-go/lmbook/payment/service/wechat"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitLogger, InitTestDB)

var wechatNativeSvcSet = wire.NewSet(
	ioc.InitWechatClient,
	dao.NewPaymentGORMDAO,
	repository.NewPaymentRepository,
	ioc.InitWechatNativeService,
	ioc.InitWechatConfig)

func InitWechatNativeService() *wechat.NativePaymentService {
	wire.Build(wechatNativeSvcSet, thirdPartySet)
	return new(wechat.NativePaymentService)
}
