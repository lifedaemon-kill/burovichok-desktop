package main

import (
	"context"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	importerService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/importer"
	uiService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/ui"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
	"log"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	inmemory "github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
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
	blocksStore := inmemory.NewStore()

	guidebooksStore, err := sqlite.New(conf.DB)
	if err != nil {
		return errors.Wrap(err, "sqlite.New")
	}

	ui := uiService.NewService(conf.UI.Name, conf.UI.Width, conf.UI.Height, zLog, importer, blocksStore)
	if err = ui.Run(); err != nil {
		zLog.Errorw("UI service failed", "error", err)
		return err
	}

	zLog.Infow("Application shutting down...")
	zLog.Infow("Shutdown completed")

	return nil
}
