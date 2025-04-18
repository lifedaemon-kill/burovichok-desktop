package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

func NewGuidebookStorage(db *sqlx.DB) (GuidebooksStorage, error) {
	return &sqlite{
		DB: db,
	}, nil
}

func (s sqlite) AddOilPlaces(fields []models.OilField) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddInstrumentType(types []models.InstrumentType) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddProductiveHorizon(horizons []models.ProductiveHorizon) error {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetAllOilPlaces() ([]models.OilField, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetAllInstrumentType() ([]models.InstrumentType, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetAllProductiveHorizon() ([]models.ProductiveHorizon, error) {
	//TODO implement me
	panic("implement me")
}
