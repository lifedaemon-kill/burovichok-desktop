package z

import "go.uber.org/zap"

var (
	Log *zap.SugaredLogger
)

func InitLogger() {
	notsugar, err := zap.NewProduction()
	Log = notsugar.Sugar()
	if err != nil {
		panic(err)
	}
}
