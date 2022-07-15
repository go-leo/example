package main

import (
	"context"
	"time"

	"github.com/go-leo/leo"
	"github.com/go-leo/leo/global"
	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	"github.com/go-leo/leo/runner"
)

func main() {
	ctx := context.Background()
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))
	// 初始化app
	app := leo.NewApp(
		leo.Name("callabledemo"),
		// 日志打印
		leo.Logger(global.Logger()),
		leo.Callable(new(CallableDemo)),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

var _ runner.Callable = new(CallableDemo)

type CallableDemo struct{}

func (c *CallableDemo) String() string {
	return "CallableDemo"
}

func (c *CallableDemo) Invoke(ctx context.Context) error {
	global.Logger().Info("start invoke")
	defer global.Logger().Info("stop invoke")
	global.Logger().Info("will sleep 30s")
	select {
	case <-ctx.Done():
	case <-time.After(30 * time.Second):
	}
	return nil
}
