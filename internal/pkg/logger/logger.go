package zlog

import "go.uber.org/zap"

var (
	W *zap.Logger
)

func InitLogger() {
	var err error
	W, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
