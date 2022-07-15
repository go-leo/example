package main

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-leo/leo"
	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	middlewarecontext "github.com/go-leo/leo/middleware/context"
	middlewarelog "github.com/go-leo/leo/middleware/log"
	"github.com/go-leo/leo/middleware/requestid"
	"github.com/go-leo/leo/registry/factory"
	"github.com/go-leo/leo/runner/net/http/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-leo/example/api/helloworld"
)

var APPVersion string

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.PlainText())
	uri, err := url.Parse("consul://localhost:8500")
	if err != nil {
		panic(err)
	}
	registrar, err := factory.NewRegistrar(uri)
	if err != nil {
		panic(err)
	}

	app := leo.NewApp(
		leo.Service(helloworld.GreeterServiceDesc(NewGreeter())),
		leo.Name("registrydemo"),
		leo.Version(APPVersion),
		leo.Logger(logger),
		leo.HTTP(&leo.HttpOptions{
			Routers: []server.Router{
				{
					HTTPMethod:   http.MethodGet,
					Path:         "/time",
					HandlerFuncs: []gin.HandlerFunc{Time},
				},
			},
		}),
		leo.GRPC(&leo.GRPCOptions{
			UnaryServerInterceptors: []grpc.UnaryServerInterceptor{
				requestid.GRPCServerMiddleware(),
				middlewarecontext.GRPCServerMiddleware(func(ctx context.Context) context.Context {
					traceID, _ := requestid.FromContext(ctx)
					return log.NewContext(ctx, logger.With(log.F{K: "TraceID", V: traceID}))
				}),
				middlewarelog.GRPCServerMiddleware(
					func(ctx context.Context) log.Logger { return log.FromContextOrDiscard(ctx) },
					middlewarelog.WithSkip("/grpc.health.v1.Health/Check"),
				),
			},
		}),
		leo.Management(&leo.ManagementOptions{}),
		leo.Registrar(registrar),
	)
	err = app.Run(ctx)
	if err != nil {
		panic(err)
	}
}

type Greeter struct {
	helloworld.UnimplementedGreeterServer
}

func NewGreeter() *Greeter {
	return &Greeter{}
}

func (ctrl *Greeter) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(time.RFC3339Nano))
		_ = grpc.SetTrailer(ctx, trailer)
	}()
	header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(time.RFC3339Nano)})
	_ = grpc.SendHeader(ctx, header)
	return &helloworld.HelloReply{Message: "hello " + request.GetName()}, nil
}

func Time(c *gin.Context) {
	c.String(http.StatusOK, "current time is %s", time.Now().Local().Format("2006-01-02 15:04:05"))
}
