package grpc

import (
	smsv1 "basic-go/lmbook/api/proto/gen/sms/v1"
	"basic-go/lmbook/sms/service"
	"context"
	"google.golang.org/grpc"
)

type SmsServiceServer struct {
	smsv1.UnimplementedSmsServiceServer
	svc service.Service
}

func NewSmsServiceServer(svc service.Service) *SmsServiceServer {
	return &SmsServiceServer{
		svc: svc,
	}
}

func (s *SmsServiceServer) Register(server grpc.ServiceRegistrar) {
	smsv1.RegisterSmsServiceServer(server, s)
}

func (s *SmsServiceServer) Send(ctx context.Context, req *smsv1.SmsSendRequest) (*smsv1.SmsSendResponse, error) {
	err := s.svc.Send(ctx, req.TplId, req.Args, req.Numbers...)
	return &smsv1.SmsSendResponse{}, err
}
