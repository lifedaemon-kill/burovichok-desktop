package database

import (
	"context"

	"github.com/google/uuid"

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
func NewService(pg *postgres.Postgres, zLog logger.Logger) *Service {
	zLog.Infow("Postgres repository initialized successfully")
	return &Service{pg: pg, log: zLog}
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

// GetAllResearchTypes возвращает все ResearchTypes
func (d *Service) GetAllResearchTypes() ([]models.ResearchType, error) {
	items, err := d.pg.GetAllResearchType(context.Background())
	if err != nil {
		d.log.Errorw("GetAllResearchTypes failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllResearchTypes succeeded", "count", len(items))

	return items, nil
}

// SaveInstrumentTypes сохраняет набор InstrumentType
func (d *Service) SaveInstrumentTypes(types []models.InstrumentType) error {
	if err := d.pg.AddInstrumentType(context.Background(), types); err != nil {
		d.log.Errorw("SaveInstrumentTypes failed", "error", err)
		return err
	}
	d.log.Debugw("SaveInstrumentTypes succeeded", "count", len(types))

	return nil
}

// SaveOilFields сохраняет набор OilField
func (d *Service) SaveOilFields(fields []models.OilField) error {
	if err := d.pg.AddOilField(context.Background(), fields); err != nil {
		d.log.Errorw("SaveOilFields failed", "error", err)
		return err
	}
	d.log.Debugw("SaveOilFields succeeded", "count", len(fields))

	return nil
}

// SaveProductiveHorizons сохраняет набор ProductiveHorizon
func (d *Service) SaveProductiveHorizons(horizons []models.ProductiveHorizon) error {
	if err := d.pg.AddProductiveHorizon(context.Background(), horizons); err != nil {
		d.log.Errorw("SaveProductiveHorizons failed", "error", err)
		return err
	}
	d.log.Debugw("SaveProductiveHorizons succeeded", "count", len(horizons))

	return nil
}

// SaveResearchTypes сохраняет набор ResearchType
func (d *Service) SaveResearchTypes(researches []models.ResearchType) error {
	if err := d.pg.AddResearchType(context.Background(), researches); err != nil {
		d.log.Errorw("SaveResearchTypes failed", "error", err)
		return err
	}
	d.log.Debugw("SaveResearchTypes succeeded", "count", len(researches))

	return nil
}

func (d *Service) GetBlockFourByResearchID(ctx context.Context, researchID uuid.UUID) ([]models.TableFour, error) {
	res, err := d.pg.GetBlockFourByID(ctx, researchID)
	if err != nil {
		d.log.Errorw("GetBlockFourByWell failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetBlockFourByWell succeeded")

	return res, nil
}

func (d *Service) SaveBlockFour(ctx context.Context, items []models.TableFour) (uuid.UUID, error) {
	id, err := d.pg.AddBlockFour(ctx, items)
	if err != nil {
		d.log.Errorw("SaveResearchTypes failed", "error", err)
		return uuid.Nil, err
	}
	d.log.Debugw("SaveResearchTypes succeeded", "count", len(items))

	return id, nil
}
