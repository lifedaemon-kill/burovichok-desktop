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
func (d *Service) GetAllReports(ctx context.Context) ([]models.TableFive, error) {
	reports, err := d.pg.GetAllTableFive(ctx)
	if err != nil {
		d.log.Errorw("GetAllReports failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllReports succeeded", "count", len(reports))

	return reports, nil
}

// SaveReport сохраняет TableFive и возвращает ID
func (d *Service) SaveReport(ctx context.Context, tableFive models.TableFive) (int64, error) {
	id, err := d.pg.AddBlockFive(ctx, tableFive)
	if err != nil {
		d.log.Errorw("SaveReport failed", "error", err)
		return 0, err
	}
	d.log.Debugw("SaveReport succeeded", "id", id)

	return id, nil
}

// GetAllInstrumentTypes возвращает все InstrumentType
func (d *Service) GetAllInstrumentTypes(ctx context.Context) ([]models.InstrumentType, error) {
	items, err := d.pg.GetAllInstrumentType(ctx)
	if err != nil {
		d.log.Errorw("GetAllInstrumentTypes failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllInstrumentTypes succeeded", "count", len(items))

	return items, nil
}

// GetAllProductiveHorizons возвращает все ProductiveHorizon
func (d *Service) GetAllProductiveHorizons(ctx context.Context) ([]models.ProductiveHorizon, error) {
	items, err := d.pg.GetAllProductiveHorizon(ctx)
	if err != nil {
		d.log.Errorw("GetAllProductiveHorizons failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllProductiveHorizons succeeded", "count", len(items))

	return items, nil
}

// GetAllOilFields возвращает все OilField
func (d *Service) GetAllOilFields(ctx context.Context) ([]models.OilField, error) {
	items, err := d.pg.GetAllOilField(ctx)
	if err != nil {
		d.log.Errorw("GetAllOilFields failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllOilFields succeeded", "count", len(items))

	return items, nil
}

// GetAllResearchTypes возвращает все ResearchTypes
func (d *Service) GetAllResearchTypes(ctx context.Context) ([]models.ResearchType, error) {
	items, err := d.pg.GetAllResearchType(ctx)
	if err != nil {
		d.log.Errorw("GetAllResearchTypes failed", "error", err)
		return nil, err
	}
	d.log.Debugw("GetAllResearchTypes succeeded", "count", len(items))

	return items, nil
}

// SaveInstrumentTypes сохраняет набор InstrumentType
func (d *Service) SaveInstrumentTypes(ctx context.Context, types []models.InstrumentType) error {
	if err := d.pg.AddInstrumentType(ctx, types); err != nil {
		d.log.Errorw("SaveInstrumentTypes failed", "error", err)
		return err
	}
	d.log.Debugw("SaveInstrumentTypes succeeded", "count", len(types))

	return nil
}

// SaveOilFields сохраняет набор OilField
func (d *Service) SaveOilFields(ctx context.Context, fields []models.OilField) error {
	if err := d.pg.AddOilField(ctx, fields); err != nil {
		d.log.Errorw("SaveOilFields failed", "error", err)
		return err
	}
	d.log.Debugw("SaveOilFields succeeded", "count", len(fields))

	return nil
}

// SaveProductiveHorizons сохраняет набор ProductiveHorizon
func (d *Service) SaveProductiveHorizons(ctx context.Context, horizons []models.ProductiveHorizon) error {
	if err := d.pg.AddProductiveHorizon(ctx, horizons); err != nil {
		d.log.Errorw("SaveProductiveHorizons failed", "error", err)
		return err
	}
	d.log.Debugw("SaveProductiveHorizons succeeded", "count", len(horizons))

	return nil
}

// SaveResearchTypes сохраняет набор ResearchType
func (d *Service) SaveResearchTypes(ctx context.Context, researches []models.ResearchType) error {
	if err := d.pg.AddResearchType(ctx, researches); err != nil {
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
		d.log.Errorw("SaveBlockFour failed", "error", err)
		return uuid.Nil, err
	}
	d.log.Debugw("SaveBlockFour succeeded", "count", len(items))

	return id, nil
}
