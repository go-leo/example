package main

import (
	"context"
	"fmt"

	"github.com/go-leo/cobra"
	"github.com/go-leo/leo/v2"
	"github.com/go-leo/leo/v2/global"
	spf13cobra "github.com/spf13/cobra"
)

func main() {
	root := &spf13cobra.Command{
		Use:   "",
		Short: "",
		RunE: func(cmd *spf13cobra.Command, args []string) error {
			fmt.Println("root", args)
			return nil
		},
	}
	version := &spf13cobra.Command{
		Use:   "version",
		Short: "version",
		RunE: func(cmd *spf13cobra.Command, args []string) error {
			fmt.Println("version: ", args)
			return nil
		},
	}
	root.AddCommand(version)
	command := cobra.New(root)

	app := leo.NewApp(
		leo.Name("cobrademo"),
		// 日志打印
		leo.Logger(global.Logger()),
		leo.Callable(command),
	)
	// 运行app
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
