package main

import (
	"fmt"

	"github.com/go-leo/leo/v2/config"
	"github.com/go-leo/leo/v2/config/medium/file"
	"github.com/go-leo/leo/v2/config/parser"
	"github.com/go-leo/leo/v2/config/valuer"
)

func main() {
	manager := config.NewManager(
		config.WithLoader(file.NewLoader("./cmd/configdemo/config/config.yaml")),
		config.WithParser(parser.NewYamlParser()),
		config.WithValuer(valuer.NewTrieTreeValuer()),
		config.WithWatcher(file.NewWatcher("./cmd/configdemo/config/config.yaml")),
	)
	err := manager.ReadConfig()
	if err != nil {
		panic(err)
	}
	asMap := manager.AsMap()
	fmt.Println(asMap)
}
