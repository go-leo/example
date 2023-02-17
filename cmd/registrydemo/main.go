package main

import (
	"context"
	"net/url"

	gonicgin "github.com/gin-gonic/gin"
	"github.com/go-leo/gin"
	leogrpc "github.com/go-leo/grpc"
	"github.com/go-leo/grpcproxy"
	"github.com/go-leo/leo/v2"
	"github.com/go-leo/leo/v2/global"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	"github.com/go-leo/leo/v2/registry/factory"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-leo/example/v2/api/helloworld"
)

var APPVersion string

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.PlainText())
	global.SetLogger(logger)
	uri, err := url.Parse("consul://localhost:8500?health_check_path=/health/check")
	if err != nil {
		panic(err)
	}
	registrar, err := factory.NewRegistrar(uri)
	if err != nil {
		panic(err)
	}

	services := []leogrpc.Service{{Impl: new(GreeterService), Desc: helloworld.Greeter_ServiceDesc}}
	grpcSrv, err := leogrpc.NewServer(0, services, leogrpc.Registrar(registrar))
	if err != nil {
		panic(err)
	}

	// Set up a connection to the server.
	conn, err := googlegrpc.Dial("localhost:9090", googlegrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	engine := grpcproxy.AppendRoutes(gonicgin.New(), helloworld.GreeterProxyRoutes(helloworld.NewGreeterClient(conn))...)
	httpSrv, err := gin.NewServer(0, engine, gin.Registrar(registrar))
	if err != nil {
		panic(err)
	}
	// 初始化app
	app := leo.NewApp(
		leo.Name("grpcproxydemo"),
		leo.Logger(logger),
		leo.Runnable(grpcSrv),
		leo.Runnable(httpSrv),
	)
	err = app.Run(ctx)
	if err != nil {
		panic(err)
	}
}

type GreeterService struct {
	helloworld.UnimplementedGreeterServer
}

func (ctrl *GreeterService) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	global.Logger().Info("hello " + request.GetName())
	return &helloworld.HelloReply{Message: "hello " + request.GetName()}, nil
}
