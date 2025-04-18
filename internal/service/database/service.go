package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
)

// Service отвечает за работу с базой данных
type Service interface {
	importer
	exporter
	updater
	deleter
}

type dbService struct {
	g sqlite.GuidebooksStorage
	b sqlite.BlocksStorage
}

func NewService(db *sqlx.DB, zLog logger.Logger) (Service, error) {
	//Репозитории для работы с данными
	guidebooksRepository, err := sqlite.NewGuidebookStorage(db)
	if err != nil {
		zLog.Errorw("sqlite.NewGuidebookStorage", "error", err)
		return nil, err
	}
	blocksRepository, err := sqlite.NewBlockStorage(db)
	if err != nil {
		zLog.Errorw("sqlite.NewBlockStorage", "error", err)
		return nil, err
	}
	zLog.Infow("Repositories initialized successfully")

	return &dbService{
		g: guidebooksRepository,
		b: blocksRepository,
	}, nil
}
