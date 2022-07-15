package main

import (
	"context"

	"github.com/robfig/cron/v3"

	"github.com/go-leo/leo"
	"github.com/go-leo/leo/global"
	"github.com/go-leo/leo/log"
	"github.com/go-leo/leo/log/zap"
	cronmiddleware "github.com/go-leo/leo/middleware/cron"
	"github.com/go-leo/leo/middleware/recovery"
	crontask "github.com/go-leo/leo/runner/task/cron"
)

func main() {
	ctx := context.Background()
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))
	// 初始化app
	app := leo.NewApp(
		leo.Name("crondemo"),
		// 日志打印
		leo.Logger(global.Logger()),
		leo.Cron(&leo.CronOptions{
			Jobs: []*crontask.Job{Print()},
			Middlewares: []cron.JobWrapper{
				recovery.CronMiddleware(global.Logger().Clone()),
				cronmiddleware.SkipIfStillRunning(global.Logger().Clone()),
			},
		}),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

func Print() *crontask.Job {
	return crontask.NewJob("print", "@every 5s", func() {
		global.Logger().Info("this is from cron")
	})
}
