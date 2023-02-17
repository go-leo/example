package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-leo/grpc"
	"github.com/go-leo/grpcproxy"
	"github.com/go-leo/leo/v2"
	leohttp "github.com/go-leo/leo/v2/http"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-leo/example/v2/api/helloworld"
)

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON())

	grpcSrv, err := grpc.NewServer(
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

	// Set up a connection to the server.
	conn, err := googlegrpc.Dial("localhost:9090", googlegrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	engine := grpcproxy.AppendRoutes(gin.New(), helloworld.GreeterProxyRoutes(helloworld.NewGreeterClient(conn))...)
	httpSrv, err := leohttp.NewServer(8080, engine)
	if err != nil {
		panic(err)
	}
	// 初始化app
	app := leo.NewApp(
		leo.Name("grpcproxydemo"),
		leo.Logger(logger),
		leo.Runnable(grpcSrv),
		leo.HTTP(httpSrv),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
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
