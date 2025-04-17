package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	importerService "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/service/importer"
	uiService "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/service/ui"
	inmemory "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/storage/inmemory"
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

	importer := importerService.NewService()
	store := inmemory.NewStore()

	ui := uiService.NewService(conf.UI.Name, conf.UI.Width, conf.UI.Height, zLog, importer, store)
	if err := ui.Run(); err != nil {
		zLog.Errorw("UI service failed", "error", err)
		return err
	}

	zLog.Infow("Application shutting down...")
	zLog.Infow("Shutdown completed")

	return nil
}
