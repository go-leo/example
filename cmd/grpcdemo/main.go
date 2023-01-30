package main

import (
	"context"

	"github.com/go-leo/leo/v2"
	"github.com/go-leo/leo/v2/grpc"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"

	"github.com/go-leo/example/v2/api/helloworld"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())

	srv, err := grpc.NewServer(
		9090,
		[]grpc.Service{
			{
				Impl: new(GreeterService),
				Desc: helloworld.Greeter_ServiceDesc,
			},
		},
	)
	if err != nil {
		panic(err)
	}
	// 初始化app
	app := leo.NewApp(
		leo.Name("grpcdemo"), // 服务名
		leo.Logger(logger),   // 日志组件
		leo.GRPC(srv),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

type GreeterService struct {
	helloworld.UnimplementedGreeterServer
}

func (ctrl *GreeterService) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello " + request.GetName()}, nil
}
