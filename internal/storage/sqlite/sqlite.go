package sqlite

import (
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	DB *sqlx.DB
}

func New(conf config.DBConf) (storage.GuidebooksStorage, error) {
	db, err := sqlx.Connect("sqlite3", conf.DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "не удалось открыть бд")
	}

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
