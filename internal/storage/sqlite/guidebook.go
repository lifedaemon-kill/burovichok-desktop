package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
)

func NewGuidebookRepository(db *sqlx.DB) (GuidebooksStorage, error) {
	return &sqlite{
		DB: db,
	}, nil
}

func (s sqlite) AddOilPlaces() ([]models.OilPlaces, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddInstrumentType() ([]models.InstrumentType, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddProductiveHorizon() ([]models.ProductiveHorizon, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetAllOilPlaces() ([]models.OilPlaces, error) {
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
