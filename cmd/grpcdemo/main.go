package main

import (
	"context"

	"github.com/go-leo/leo"
	"google.golang.org/grpc"

	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	middlewarecontext "github.com/go-leo/leo/middleware/context"
	middlewarelog "github.com/go-leo/leo/middleware/log"
	"github.com/go-leo/leo/middleware/requestid"

	"github.com/go-leo/example/api/helloworld"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())
	// 初始化app
	app := leo.NewApp(
		leo.Service(helloworld.GreeterServiceDesc(new(GreeterService))),
		leo.Name("grpcdemo"),
		leo.Logger(logger),
		leo.GRPC(&leo.GRPCOptions{
			Port: 9090,
			UnaryServerInterceptors: []grpc.UnaryServerInterceptor{
				requestid.GRPCServerMiddleware(),
				middlewarecontext.GRPCServerMiddleware(func(ctx context.Context) context.Context {
					traceID, _ := requestid.FromContext(ctx)
					return log.NewContext(ctx, logger.Clone().With(log.F{K: "TraceID", V: traceID}))
				}),
				middlewarelog.GRPCServerMiddleware(func(ctx context.Context) log.Logger { return log.FromContextOrDiscard(ctx) }),
			},
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
