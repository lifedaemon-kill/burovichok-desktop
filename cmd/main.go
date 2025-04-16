package main

import (
	"fmt"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/logger/z"
)

const configPath = "config/config.yaml"

func main() {
	fmt.Println("start app")
	conf := config.Load(configPath)

	z.InitLogger()
	defer z.Log.Sync()
	z.Log.Infow("zap logger init success")

	z.Log.Infow("config load success", "conf", conf)

	fmt.Println("end")
}
