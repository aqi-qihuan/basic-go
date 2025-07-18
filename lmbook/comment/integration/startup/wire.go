//go:build wireinject

package startup

import (
	grpc2 "basic-go/lmbook/comment/grpc"
	"basic-go/lmbook/comment/repository"
	"basic-go/lmbook/comment/repository/dao"
	"basic-go/lmbook/comment/service"
	"basic-go/lmbook/pkg/logger"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewGrpcServer,
)

var thirdProvider = wire.NewSet(
	logger.NewNoOpLogger,
	InitTestDB,
)

func InitGRPCServer() *grpc2.CommentServiceServer {
	wire.Build(thirdProvider, serviceProviderSet)
	return new(grpc2.CommentServiceServer)
}
