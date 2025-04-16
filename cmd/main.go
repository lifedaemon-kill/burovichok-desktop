package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-backend/internal/config"
	"github.com/lifedaemon-kill/burovichok-backend/internal/pkg/logger"
	uiService "github.com/lifedaemon-kill/burovichok-backend/internal/pkg/service/ui"
)

func main() {
	ctx := context.Background()
	if err := bootstrap(ctx); err != nil {
		log.Fatalf("[main] bootstrap failed: %v", err)
	}
}

func bootstrap(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.Load(config.PathConfig)
	if err != nil {
		return errors.Wrap(err, "config.Load")
	}

	zLog, err := logger.NewLogger(conf.Logger.Env)
	if err != nil {
		return errors.Wrap(err, "logger.NewLogger")
	}

	zLog.Infow("Logger and config initialized successfully")

	ui := uiService.NewService("burovichok", 800, 400)
	if err = ui.Run(); err != nil {
		zLog.Errorw("UI service failed", "error", err)
		return err
	}

	zLog.Infow("Application shutting down...")
	zLog.Infow("Shutdown completed")

	return nil
}
