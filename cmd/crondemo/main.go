package main

import (
	"context"

	"github.com/go-leo/cron"
	"github.com/go-leo/leo/v2"
	"github.com/go-leo/leo/v2/global"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	cronmiddleware "github.com/go-leo/leo/v2/middleware/cron"
)

func main() {
	ctx := context.Background()
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))
	// 初始化app
	task := cron.New([]*cron.Job{Print()},
		cron.Middleware(
			cronmiddleware.CronMiddleware(global.Logger().Clone()),
			cronmiddleware.SkipIfStillRunning(global.Logger().Clone()),
		))
	app := leo.NewApp(
		leo.Name("crondemo"),
		// 日志打印
		leo.Logger(global.Logger()),
		leo.Runnable(task),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

func Print() *cron.Job {
	return cron.NewJob("print", "@every 5s", func() {
		global.Logger().Info("this is from cron")
	})
}
