package ioc

import (
	grpc2 "basic-go/lmbook/cronjob/grpc"
	"basic-go/lmbook/pkg/grpcx"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(cronJobGrpc *grpc2.CronJobServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}

	// 解析端口
	port := 8081 // 默认端口
	if cfg.Addr != "" {
		addrParts := strings.Split(cfg.Addr, ":")
		if len(addrParts) > 1 {
			if p, err := strconv.Atoi(addrParts[1]); err == nil {
				port = p
			}
		}
	}

	server := grpc.NewServer()
	cronJobGrpc.Register(server)
	return &grpcx.Server{
		Server: server,
		Port:   port,
	}
}
