package main

import (
	"context"
	"time"

	"github.com/go-leo/leo"
	"google.golang.org/grpc"

	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"

	"github.com/go-leo/example/api/helloworld"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())
	// 初始化app
	app := leo.NewApp(
		leo.Name("grpcdemo"),
		leo.Logger(logger),
		leo.Service(helloworld.GreeterServiceDesc(new(GreeterService))),
		leo.GRPC(&leo.GRPCOptions{
			Port: 9090,
			UnaryServerInterceptors: []grpc.UnaryServerInterceptor{
				func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
					ctx, cancelFunc := context.WithTimeout(ctx, time.Second)
					defer cancelFunc()
					return handler(ctx, req)
				},
			},
			TLSConf:           nil,
			GRPCServerOptions: []grpc.ServerOption{},
		}),
		leo.Management(&leo.ManagementOptions{
			Port: 16060,
		}),
	)
	// 运行app
	err := app.Run(ctx)
	if err != nil {
		panic(err)
	}
}

type GreeterService struct {
	helloworld.UnimplementedGreeterServer
}

func (ctrl *GreeterService) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello " + request.GetName()}, nil
}
