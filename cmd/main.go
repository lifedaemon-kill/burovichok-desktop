package main

import (
	"context"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	calcService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/calc"
	converterService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/convertor"
	importerService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/importer"
	uiService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/ui"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
	"log"
	"os/signal"
	"syscall"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
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

	converter := converterService.NewService()
	calc := calcService.NewService()
	importer := importerService.NewService(calc, converter)

	inMemoryBlocksStorage := inmemory.NewInMemoryBlocksStorage()

	//Запускаем sqlite
	db, err := sqlite.NewDB(conf.DB)
	if err != nil {
		zLog.Errorw("sqlite.NewDB", "error", err)
		return err
	}
	zLog.Infow("Database initialized successfully")

	// Применяем миграции
	if err = goose.Up(db.DB, conf.DB.MigrationsPath); err != nil {
		zLog.Errorw("goose.Up", "error", err)
		return err
	}
	zLog.Infow("Migrations initialized successfully")

	//Репозитории для работы с данными
	guidebooksRepository, err := sqlite.NewGuidebookRepository(db)
	if err != nil {
		zLog.Errorw("sqlite.NewGuidebookRepository", "error", err)
		return err
	}
	blocksRepository, err := sqlite.NewBlockRepository(db)
	if err != nil {
		zLog.Errorw("sqlite.NewBlockRepository", "error", err)
		return err
	}
	zLog.Infow("Repositories initialized successfully")

	ui := uiService.NewService(conf.UI.Name, conf.UI.Width, conf.UI.Height, zLog, importer, converter, inMemoryBlocksStorage, blocksRepository, guidebooksRepository)
	if err = ui.Run(); err != nil {
		zLog.Errorw("UI service failed", "error", err)
		return err
	}

	zLog.Infow("Application shutting down...")
	zLog.Infow("Shutdown completed")

	return nil
}
