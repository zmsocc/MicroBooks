package ioc

import (
	"github.com/spf13/viper"
	grpc2 "github.com/zmsocc/practice/webook/interactive/grpc"
	"github.com/zmsocc/practice/webook/pkg/grpcx"
	"google.golang.org/grpc"
)

func InitGRPCxServer(intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Addr string `yaml:"addr"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	intrServer.Register(server)
	return &grpcx.Server{
		Server: server,
		Addr:   cfg.Addr,
	}
}
