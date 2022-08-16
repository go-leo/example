package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-leo/leo"
	"google.golang.org/grpc"

	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	middlewarecontext "github.com/go-leo/leo/middleware/context"
	middlewarelog "github.com/go-leo/leo/middleware/log"
	"github.com/go-leo/leo/middleware/requestid"
	httpserver "github.com/go-leo/leo/runner/net/http/server"

	"github.com/go-leo/example/api/helloworld"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())
	// 初始化app
	app := leo.NewApp(
		leo.Service(helloworld.GreeterServiceDesc(new(GreeterService))),
		leo.Name("proxydemo"),
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
		leo.HTTP(&leo.HttpOptions{
			Port: 8080,
			GRPCDialOptions: []grpc.DialOption{
				grpc.WithChainUnaryInterceptor(
					requestid.GRPCClientMiddleware(),
					middlewarelog.GRPCClientMiddleware(func(ctx context.Context) log.Logger { return log.FromContextOrDiscard(ctx) }),
				),
			},
			// 全局中间件
			GinMiddlewares: []gin.HandlerFunc{
				requestid.GinMiddleware(),
				middlewarecontext.GinMiddleware(func(ctx context.Context) context.Context {
					traceID, _ := requestid.FromContext(ctx)
					return log.NewContext(ctx, logger.Clone().With(log.F{K: "TraceID", V: traceID}))
				}),
				middlewarelog.GinMiddleware(func(ctx context.Context) log.Logger { return log.FromContextOrDiscard(ctx) }),
			},
			Routers: []httpserver.Router{
				{
					HTTPMethods:  []string{http.MethodGet},
					Path:         "/time",
					HandlerFuncs: []gin.HandlerFunc{Time},
				},
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

func Time(c *gin.Context) {
	c.String(http.StatusOK, time.Now().Format(time.RFC3339))
}

type GreeterService struct {
	helloworld.UnimplementedGreeterServer
}

func (ctrl *GreeterService) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello " + request.GetName()}, nil
}
