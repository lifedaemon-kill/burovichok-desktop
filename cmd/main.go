package main

import (
	"fmt"
	"internal/pkg/logger"
)

const configPath = "config/config.yaml"

func main() {
	fmt.Println("Start")
	conf := config.Load(configPath)
	zlog.InitLogger()
	defer W.Sync()
	fmt.Println(conf)

	zlog.W.Info("Hello zap logger!")
	fmt.Println("End")
}
