package z

import (
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/models"
	"go.uber.org/zap"
)

var (
	Log *zap.SugaredLogger
)

func InitLogger(conf config.LoggerConf) (err error) {
	var l *zap.Logger

	switch conf.Env {
	case models.ENV_PROD:
		l, err = zap.NewProduction()
	case models.ENV_DEV:
		l, err = zap.NewDevelopment()
	}

	if err != nil || l == nil {
		return err
	}

	Log = l.Sugar()
	return nil
}
