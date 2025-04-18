package database

import "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"

type importer interface {
	GetAllInstrumentsTypes() ([]models.InstrumentType, error)
	GetAllProductiveHorizons() ([]models.ProductiveHorizon, error)
	GetAllOilFields() ([]models.OilField, error)
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
