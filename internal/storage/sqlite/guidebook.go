package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func NewGuidebookStorage(db *sqlx.DB) (GuidebooksStorage, error) {
	return &sqlite{
		db: db,
	}, nil
}

func (s sqlite) GetAllOilField() ([]models.OilField, error) {
	var oilfields []models.OilField

	err := s.db.Select(&oilfields, "SELECT * FROM oilfield")
	if err != nil {
		return nil, err
	}
	return oilfields, nil
}

func (s sqlite) GetAllInstrumentType() ([]models.InstrumentType, error) {
	var instrumentTypes []models.InstrumentType
	err := s.db.Select(&instrumentTypes, "SELECT * FROM instrument_type")
	if err != nil {
		return nil, err
	}
	return instrumentTypes, nil
}

func (s sqlite) GetAllProductiveHorizon() ([]models.ProductiveHorizon, error) {
	var productiveHorizons []models.ProductiveHorizon
	err := s.db.Select(&productiveHorizons, "SELECT * FROM productive_horizon")
	if err != nil {
		return nil, err
	}
	return productiveHorizons, nil
}

func (s sqlite) AddOilField(fields []models.OilField) error {
	for _, field := range fields {
		_, err := s.db.NamedExec("INSERT INTO oilfield (name) VALUES (:name)", field.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s sqlite) AddInstrumentType(types []models.InstrumentType) error {
	for _, instrumentType := range types {
		_, err := s.db.NamedExec("INSERT INTO instrument_type (name) VALUES (:name)", instrumentType.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s sqlite) AddProductiveHorizon(horizons []models.ProductiveHorizon) error {
	for _, horizon := range horizons {
		_, err := s.db.NamedExec("INSERT INTO productive_horizon (name) VALUES (:name)", horizon.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
