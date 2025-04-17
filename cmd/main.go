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
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
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

	db, err := sqlite.NewDB(conf.DB)
	if err != nil {
		return errors.Wrapf(err, "не удалось открыть бд")
	}
	// Применяем миграции
	if err = goose.Up(db.DB, "migrations"); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	//Репозитории для работы с данными
	guidebooksRepository, err := sqlite.NewGuidebookRepository(db)
	if err != nil {
		return errors.Wrap(err, "sqlite.NewGuidebookRepository")
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
