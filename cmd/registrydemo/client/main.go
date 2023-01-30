package main

import (
	"context"
	"net/url"

	leogrpc "github.com/go-leo/leo/v2/grpc"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	middlewarelog "github.com/go-leo/leo/v2/middleware/log"
	"github.com/go-leo/leo/v2/registry/factory"
	googlegprc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-leo/example/v2/api/helloworld"
)

var APPVersion string

func main() {
	ctx := context.Background()
	logger := zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.PlainText())
	target := "consul://localhost:8500/grpcproxydemo"
	uri, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	chain := []googlegprc.UnaryClientInterceptor{
		middlewarelog.GRPCClientMiddleware(func(ctx context.Context) log.Logger {
			return logger.Clone()
		}),
	}
	discovery, _ := factory.NewDiscovery(uri)
	dialOptions := []googlegprc.DialOption{
		googlegprc.WithResolvers(leogrpc.NewResolverBuilder(discovery)),
		googlegprc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		googlegprc.WithTransportCredentials(insecure.NewCredentials()),
		googlegprc.WithChainUnaryInterceptor(chain...),
	}

	cc, err := googlegprc.DialContext(ctx, target, dialOptions...)
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
