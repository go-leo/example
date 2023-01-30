package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"

	"github.com/go-leo/leo/v2"
	"github.com/go-leo/leo/v2/global"
	"github.com/go-leo/leo/v2/log"
	"github.com/go-leo/leo/v2/log/zap"
	"github.com/go-leo/leo/v2/runner"
)

func main() {
	file, err := os.OpenFile("/tmp/cpu_percent", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	global.SetLogger(zap.New(zap.LevelAdapt(log.Debug), zap.Console(true), zap.JSON()))
	// 初始化app
	app := leo.NewApp(
		leo.Name("runnabledemo"),
		// 日志打印
		leo.Logger(global.Logger()),
		leo.Runnable(NewRunnableDemo(file)),
	)
	// 运行app
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

var _ runner.Runnable = new(RunnableDemo)

type RunnableDemo struct {
	w     io.WriteCloser
	exitC chan struct{}
}

func NewRunnableDemo(w io.WriteCloser) *RunnableDemo {
	return &RunnableDemo{w: w, exitC: make(chan struct{})}
}

func (h *RunnableDemo) String() string { return "RunnableDemo" }

func (h *RunnableDemo) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-h.exitC:
			ticker.Stop()
			return h.w.Close()
		case t := <-ticker.C:
			percent, err := cpu.Percent(time.Second, true)
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintf(h.w, "%s cpu percent is %f\n", t.String(), percent); err != nil {
				return err
			}
		}
	}
}

func (h *RunnableDemo) Stop(ctx context.Context) error {
	close(h.exitC)
	return nil
}
