package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-leo/leo/v2/middleware/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-leo/example/v2/api/helloworld"
)

type A int

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithWriteBufferSize(1024 * 1024),
		grpc.WithReadBufferSize(1024 * 1024),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024),
			grpc.MaxCallSendMsgSize(1024*1024),
		),
		grpc.WithChainUnaryInterceptor(
			trace.GRPCClientMiddleware(),
			//requestid.GRPCClientMiddleware(),
		),
	}
	clientConn, err := grpc.DialContext(ctx, "localhost:9090", dialOptions...)
	if err != nil {
		panic(err)
	}
	client := helloworld.NewGreeterClient(clientConn)
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "xxxx"})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply)
}
