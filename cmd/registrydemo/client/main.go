package main

import (
	"context"
	"net/url"

	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	middlewarelog "github.com/go-leo/leo/middleware/log"
	"github.com/go-leo/leo/registry/factory"
	"github.com/go-leo/leo/runner/net/grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-leo/example/api/helloworld"
)

var APPVersion string

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.PlainText())
	target := "consul://localhost:8500/registrydemo"
	uri, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	chain := []grpc.UnaryClientInterceptor{
		middlewarelog.GRPCClientMiddleware(func(ctx context.Context) log.Logger {
			return logger.Clone()
		}),
	}
	discovery, _ := factory.NewDiscovery(uri)
	dialOptions := []grpc.DialOption{
		grpc.WithResolvers(client.NewResolverBuilder(discovery)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(chain...),
	}

	cc, err := grpc.DialContext(ctx, target, dialOptions...)
	if err != nil {
		panic(err)
	}
	client := helloworld.NewGreeterClient(cc)

	hello, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "Tom"})
	if err != nil {
		panic(err)
	}
	logger.Info(hello.Message)

	hello, err = client.SayHello(ctx, &helloworld.HelloRequest{Name: "Jack"})
	if err != nil {
		panic(err)
	}
	logger.Info(hello.Message)

	hello, err = client.SayHello(ctx, &helloworld.HelloRequest{Name: "Kitty"})
	if err != nil {
		panic(err)
	}
	logger.Info(hello.Message)
}
