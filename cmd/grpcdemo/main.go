package main

import (
	"context"

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
		leo.Name("grpcdemo"), // 服务名
		leo.Logger(logger),   // 日志组件
		leo.Service(helloworld.GreeterServiceDesc(new(GreeterService))), // 服务
		leo.GRPC(&leo.GRPCOptions{
			Port:                    9090,                            // grpc端口号
			UnaryServerInterceptors: []grpc.UnaryServerInterceptor{}, // grpc 拦截器
			TLSConf:                 nil,                             // tls 配置
			GRPCServerOptions:       []grpc.ServerOption{},           // grpc其他服务参数
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
