package grpc

import (
	articlev1 "basic-go/lmbook/api/proto/gen/article/v1"
	"basic-go/lmbook/article/domain"
	"basic-go/lmbook/article/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ArticleServiceServer struct {
	articlev1.UnimplementedArticleServiceServer
	service service.ArticleService
}

func NewArticleServiceServer(svc service.ArticleService) *ArticleServiceServer {
	return &ArticleServiceServer{
		service: svc,
	}
}
