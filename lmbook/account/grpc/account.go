package grpc

import (
	"basic-go/lmbook/account/domain"
	"basic-go/lmbook/account/service"
	accountv1 "basic-go/lmbook/api/proto/gen/account/v1"
	"context"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
)

type AccountService struct {
	accountV1.UnimplementedAccountServiceServer
	svc service.AccountService
}

func NewAccountService(svc service.AccountService) *AccountService {

}
