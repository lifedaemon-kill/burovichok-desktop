package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	calcService "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/service/calc"
	converterService "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/service/convertor"
	importerService "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/service/importer"
	uiService "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/service/ui"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/storage/inmemory"
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

	converter := converterService.NewService()
	calc := calcService.NewService()
	importer := importerService.NewService(calc, converter)

	store := inmemory.NewStore()

	ui := uiService.NewService(conf.UI.Name, conf.UI.Width, conf.UI.Height, zLog, importer, converter, store)
	if err := ui.Run(); err != nil {
		zLog.Errorw("UI service failed", "error", err)
		return err
	}

	zLog.Infow("Application shutting down...")
	zLog.Infow("Shutdown completed")

	return nil
}
