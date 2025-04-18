package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
)

type importer interface {
	GetAllInstrumentsTypes() ([]models.InstrumentType, error)
	GetAllProductiveHorizons() ([]models.ProductiveHorizon, error)
	GetAllOilFields() ([]models.OilField, error)
}

type exporter interface {
	SaveReport(
		tableOne []models.TableOne,
		tableTwo []models.TableTwo,
		tableThree []models.TableThree,
		tableFour []models.TableFour,
		tableFive []models.TableFive) error

	SaveInstrumentType([]models.InstrumentType) error
	SaveOilField([]models.OilField) error
	SaveProductiveHorizon([]models.ProductiveHorizon) error
}

// Service отвечает за работу с базой данных
type Service interface {
	importer
	exporter
}

type dbService struct {
	g sqlite.GuidebooksStorage
	b sqlite.BlocksStorage
}

func (d dbService) SaveReport(tableOne []models.TableOne, tableTwo []models.TableTwo, tableThree []models.TableThree, tableFour []models.TableFour, tableFive []models.TableFive) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) SaveInstrumentType(types []models.InstrumentType) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) SaveOilField(fields []models.OilField) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) SaveProductiveHorizon(horizons []models.ProductiveHorizon) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) GetAllInstrumentsTypes() ([]models.InstrumentType, error) {
	//TODO implement me
	panic("implement me")
}

func (d dbService) GetAllProductiveHorizons() ([]models.ProductiveHorizon, error) {
	//TODO implement me
	panic("implement me")
}

func (d dbService) GetAllOilFields() ([]models.OilField, error) {
	//TODO implement me
	panic("implement me")
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
