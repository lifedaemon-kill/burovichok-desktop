package database

import (
	"context"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/postgres"
)

// Service реализует Service через Postgres хранилище
type Service struct {
	pg  *postgres.Postgres
	log logger.Logger
}

// NewService создаёт Service, используя Postgres
func NewService(pg *postgres.Postgres, zLog logger.Logger) Service {
	zLog.Infow("Postgres repository initialized successfully")
	return Service{pg: pg, log: zLog}
}

// GetAllReports возвращает все TableFive
func (d *Service) GetAllReports() ([]models.TableFive, error) {
	reports, err := d.pg.GetAllTableFive(context.Background())
	if err != nil {
		d.log.Errorw("GetAllReports failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllReports succeeded", "count", len(reports))
	return reports, nil
}

// SaveReport сохраняет TableFive и возвращает ID
func (d *Service) SaveReport(tableFive models.TableFive) (int64, error) {
	id, err := d.pg.AddBlockFive(context.Background(), tableFive)
	if err != nil {
		d.log.Errorw("SaveReport failed", "error", err)
		return 0, err
	}
	d.log.Debugw("SaveReport succeeded", "id", id)
	return id, nil
}

// GetAllInstrumentTypes возвращает все InstrumentType
func (d *Service) GetAllInstrumentTypes() ([]models.InstrumentType, error) {
	items, err := d.pg.GetAllInstrumentType(context.Background())
	if err != nil {
		d.log.Errorw("GetAllInstrumentTypes failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllInstrumentTypes succeeded", "count", len(items))
	return items, nil
}

// GetAllProductiveHorizons возвращает все ProductiveHorizon
func (d *Service) GetAllProductiveHorizons() ([]models.ProductiveHorizon, error) {
	items, err := d.pg.GetAllProductiveHorizon(context.Background())
	if err != nil {
		d.log.Errorw("GetAllProductiveHorizons failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllProductiveHorizons succeeded", "count", len(items))
	return items, nil
}

// GetAllOilFields возвращает все OilField
func (d *Service) GetAllOilFields() ([]models.OilField, error) {
	items, err := d.pg.GetAllOilField(context.Background())
	if err != nil {
		d.log.Errorw("GetAllOilFields failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllOilFields succeeded", "count", len(items))
	return items, nil
}

// SaveInstrumentType сохраняет набор InstrumentType
func (d *Service) SaveInstrumentType(types []models.InstrumentType) error {
	if err := d.pg.AddInstrumentType(context.Background(), types); err != nil {
		d.log.Errorw("SaveInstrumentType failed", "error", err)
		return err
	}
	d.log.Debugw("SaveInstrumentType succeeded", "count", len(types))
	return nil
}

// SaveOilField сохраняет набор OilField
func (d *Service) SaveOilField(fields []models.OilField) error {
	if err := d.pg.AddOilField(context.Background(), fields); err != nil {
		d.log.Errorw("SaveOilField failed", "error", err)
		return err
	}
	d.log.Debugw("SaveOilField succeeded", "count", len(fields))
	return nil
}

// SaveProductiveHorizon сохраняет набор ProductiveHorizon
func (d *Service) SaveProductiveHorizon(horizons []models.ProductiveHorizon) error {
	if err := d.pg.AddProductiveHorizon(context.Background(), horizons); err != nil {
		d.log.Errorw("SaveProductiveHorizon failed", "error", err)
		return err
	}
	d.log.Debugw("SaveProductiveHorizon succeeded", "count", len(horizons))
	return nil
}
