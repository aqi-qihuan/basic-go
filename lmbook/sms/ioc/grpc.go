package ioc

import (
	"basic-go/lmbook/pkg/grpcx"
	"basic-go/lmbook/pkg/logger"
	grpc2 "basic-go/lmbook/sms/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(smsServer *grpc2.SmsServiceServer,
	ecli *clientv3.Client,
	l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		Port     int    `yaml:"port"`
		EtcdAddr string `yaml:"etcdAddr"`
		EtcdTTL  int64  `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	smsServer.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		Name:       "sms",
		L:          l,
		EtcdTTL:    cfg.EtcdTTL,
		EtcdClient: ecli,
	}
}
