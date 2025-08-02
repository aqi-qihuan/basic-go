package test

import (
	feedv1 "basic-go/lmbook/api/proto/gen/feed/v1"
	followMocks "basic-go/lmbook/api/proto/gen/follow/v1/mocks"
	"basic-go/lmbook/feed/grpc"
	"basic-go/lmbook/feed/ioc"
	"basic-go/lmbook/feed/repository"
	"basic-go/lmbook/feed/repository/cache"
	"basic-go/lmbook/feed/repository/dao"
	"basic-go/lmbook/feed/service"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func InitGrpcServer(t *testing.T) (feedv1.FeedSvcServer, *followMocks.MockFollowServiceClient, *gorm.DB) {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	feedPullEventDAO := dao.NewFeedPullEventDAO(db)
	feedPushEventDAO := dao.NewFeedPushEventDAO(db)
	cmdable := ioc.InitRedis()
	feedEventCache := cache.NewFeedEventCache(cmdable)
	feedEventRepo := repository.NewFeedEventRepo(feedPullEventDAO, feedPushEventDAO, feedEventCache)
	mockCtrl := gomock.NewController(t)
	followClient := followMocks.NewMockFollowServiceClient(mockCtrl)
	v := ioc.RegisterHandler(feedEventRepo, followClient)
	feedService := service.NewFeedService(feedEventRepo, v)
	feedEventGrpcSvc := grpc.NewFeedEventGrpcSvc(feedService)
	return feedEventGrpcSvc, followClient, db
}
