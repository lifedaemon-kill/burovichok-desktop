package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
)

type Service interface {
	exporter
	importer
}
type exporter interface {
	SaveReport(tableFive models.TableFive) (id int32, err error)

	SaveInstrumentType([]models.InstrumentType) error
	SaveOilField([]models.OilField) error
	SaveProductiveHorizon([]models.ProductiveHorizon) error
}

type importer interface {
	GetAllReports() ([]models.TableFive, error)

	GetAllInstrumentsTypes() ([]models.InstrumentType, error)
	GetAllProductiveHorizons() ([]models.ProductiveHorizon, error)
	GetAllOilFields() ([]models.OilField, error)
}

type dbService struct {
	g   sqlite.GuidebooksStorage
	b   sqlite.BlocksStorage
	log logger.Logger
}

func (d dbService) GetAllReports() ([]models.TableFive, error) {
	reports, err := d.b.GetAllTableFive()
	if err != nil {
		d.log.Errorw("Get all reports", "error", err)
		return nil, err
	}
	d.log.Debugw("Get all reports", "reports", reports)
	return reports, nil
}

func (d dbService) SaveReport(tableFive models.TableFive) (int32, error) {
	id, err := d.b.AddBlockFive(tableFive)
	if err != nil {
		d.log.Errorw("failed to save report", "tableFive", tableFive, "error", err)
		return 0, err
	}
	d.log.Debugw("Save report", "tableFive", tableFive, "id", id)
	return id, nil
}

func (d dbService) SaveInstrumentType(types []models.InstrumentType) error {
	if err := d.g.AddInstrumentType(types); err != nil {
		d.log.Errorw("saving instrument type", "error", err)
		return err
	}
	d.log.Debugw("saved instrument type", "type", types)
	return nil
}

func (d dbService) SaveOilField(fields []models.OilField) error {
	if err := d.g.AddOilField(fields); err != nil {
		d.log.Errorw("saving oil field", "error", err)
		return err
	}
	d.log.Debugw("saved oil field", "fields", fields)
	return nil
}

func (d dbService) SaveProductiveHorizon(horizons []models.ProductiveHorizon) error {
	if err := d.g.AddProductiveHorizon(horizons); err != nil {
		d.log.Errorw("saving productive horizon", "error", err)
		return err

	}
	d.log.Debugw("saved productive horizon", "horizons", horizons)
	return nil
}

func (d dbService) GetAllInstrumentsTypes() ([]models.InstrumentType, error) {
	allInstrumentType, err := d.g.GetAllInstrumentType()
	if err != nil {
		d.log.Errorw("getting instrument type", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllInstrumentsTypes", "type", allInstrumentType)
	return allInstrumentType, nil
}

func (d dbService) GetAllProductiveHorizons() ([]models.ProductiveHorizon, error) {
	allPH, err := d.g.GetAllProductiveHorizon()
	if err != nil {
		d.log.Errorw("getting productive horizon", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllProductiveHorizons", "horizons", allPH)
	return allPH, nil
}

func (d dbService) GetAllOilFields() ([]models.OilField, error) {
	allOF, err := d.g.GetAllOilField()
	if err != nil {
		d.log.Errorw("getting oil field", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllOilFields", "fields", allOF)
	return allOF, nil
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
		g:   guidebooksRepository,
		b:   blocksRepository,
		log: zLog,
	}, nil
}
