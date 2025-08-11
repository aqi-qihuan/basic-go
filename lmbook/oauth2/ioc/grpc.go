package ioc

import (
	grpc2 "basic-go/lmbook/oauth2/grpc"
	"basic-go/lmbook/pkg/grpcx"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(oauth2Server *grpc2.Oauth2ServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	oauth2Server.Register(server)
	return &grpcx.Server{
		Server: server,
	}
}
