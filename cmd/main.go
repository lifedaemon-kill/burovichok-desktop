package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	calcService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/calc"
	converterService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/convertor"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/service/database"
	importerService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/importer"
	uiService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/ui"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/inmemory"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/postgres"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"

	chartService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/chart"


	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	archiverService "github.com/lifedaemon-kill/burovichok-desktop/internal/service/archiver" // Добавим импорт для нового сервиса
)

func main() {
	ctx := context.Background()
	if err := bootstrap(ctx); err != nil {
		log.Fatalf("[main] bootstrap failed: %v", err)
	}

	<-ctx.Done()
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func bootstrap(ctx context.Context) error {
	// 1. Контекст с отменой на сигналы
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// 2. Загрузка конфига
	conf, err := config.Load(config.PathConfig)
	if err != nil {
		return errors.Wrap(err, "config.Load")
	}

	// 3. Инициализация логгера
	zLog, err := logger.NewLogger(conf.Logger.Env)
	if err != nil {
		return errors.Wrap(err, "logger.NewLogger")
	}
	zLog.Infow("Logger and config initialized successfully")

	// 4. Подключение к PostgreSQL и миграции
	pg, err := postgres.New(ctx, conf.DB, zLog)
	if err != nil {
		zLog.Errorw("postgres.New", "error", err)
		return err
	}
	zLog.Infow("Database initialized successfully")

	if err = goose.SetDialect("postgres"); err != nil {
		zLog.Errorw("goose.SetDialect", "error", err)
		return err
	}
	if err = goose.RunContext(ctx, "up", pg.GetSqlDB(), conf.DB.MigrationsPath); err != nil {
		zLog.Errorw("goose.Run up", "error", err)
		return err
	}
	zLog.Infow("Migrations applied successfully")


	minioClient, err := minio.New(conf.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Minio.AccessKey, conf.Minio.SecretKey, ""),
		Secure: conf.Minio.UseSSL,
	})
	if err != nil {
		zLog.Errorw("minio.New failed", "error", err)
		return errors.Wrap(err, "failed to initialize minio client")
	}


	bucketName := conf.Minio.BucketName
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		zLog.Errorw("minio.BucketExists failed", "bucket", bucketName, "error", err)
		return errors.Wrapf(err, "failed to check if minio bucket '%s' exists", bucketName)
	}
	if !exists {
		zLog.Infow("MinIO bucket does not exist, attempting to create", "bucket", bucketName)
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}) // Можно указать регион, если нужно
		if err != nil {
			zLog.Errorw("minio.MakeBucket failed", "bucket", bucketName, "error", err)
			return errors.Wrapf(err, "failed to create minio bucket '%s'", bucketName)
		}
		zLog.Infow("MinIO bucket created successfully", "bucket", bucketName)
	} else {
		zLog.Infow("MinIO bucket already exists", "bucket", bucketName)
	}
	zLog.Infow("MinIO client initialized successfully")


	// 5. Создание сервиса работы с БД
	dbService := database.NewService(pg, zLog)

	// 6. Инициализация доменных сервисов
	converter := converterService.NewService()
	calc := calcService.NewService(dbService, zLog)
	importer := importerService.NewService(calc, converter)
	chartSvc := chartService.NewService()
	inMemoryStorage := inmemory.NewInMemoryBlocksStorage()


	archiver := archiverService.NewService(inMemoryStorage, minioClient, bucketName, zLog)


	// 7. Запуск UI
	err = os.Setenv("LANG", "ru_RU.UTF-8")
	if err != nil {
		return errors.Wrap(err, "os.SetEnv")
	}

	ui := uiService.NewService(
		conf.UI,
		zLog,
		importer,
		converter,
		inMemoryStorage,
		dbService,
		chartSvc,
		calc,
		archiver,
	)
	
	if err = ui.Run(ctx); err != nil {
		zLog.Errorw("UI service failed", "error", err)
		return err
	}

	// 8. Грейсфул-шатдаун
	zLog.Infow("Application shutting down...")
	if err = pg.DB.Close(); err != nil {
		return errors.Wrap(err, "pg.DB.Close")
	}
	zLog.Infow("Shutdown completed")

	return nil
}
