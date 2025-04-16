package main

import (
	"fmt"
	"internal/pkg/logger/z/logger.go"
)

const configPath = "config/config.yaml"

func main() {
	fmt.Println("start app")
	z.InitLogger()
	defer Log.Sync()
	z.Log.Info("zap logger init success")

	conf := config.Load(configPath)
	z.Log.Info("config load success", "conf", conf)

	z.Log.Info("end")
}
