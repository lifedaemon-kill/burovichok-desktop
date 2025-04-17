package sqlite

import (
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	DB *sqlx.DB
}

func (s sqlite) GetBlockOne() ([]models.BlockOne, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockTwo() ([]models.BlockTwo, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockThree() ([]models.BlockThree, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockFour() ([]models.Inclinometry, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) GetBlockFive() ([]models.GeneralInformation, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockOne() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockTwo() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockThree() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFour() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (s sqlite) AddBlockFive() (int32, error) {
	//TODO implement me
	panic("implement me")
}

func NewDB(conf config.DBConf) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conf.DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlite.NewDB")
	}
	return db, err
}

type GuidebooksStorage interface {
	AddOilPlaces() ([]models.OilPlaces, error)
	AddInstrumentType() ([]models.InstrumentType, error)
	AddProductiveHorizon() ([]models.ProductiveHorizon, error)

	GetAllOilPlaces() ([]models.OilPlaces, error)
	GetAllInstrumentType() ([]models.InstrumentType, error)
	GetAllProductiveHorizon() ([]models.ProductiveHorizon, error)
}

type BlocksStorage interface {
	GetBlockOne() ([]models.BlockOne, error)
	GetBlockTwo() ([]models.BlockTwo, error)
	GetBlockThree() ([]models.BlockThree, error)
	GetBlockFour() ([]models.Inclinometry, error)
	GetBlockFive() ([]models.GeneralInformation, error)

	AddBlockOne() (int32, error)
	AddBlockTwo() (int32, error)
	AddBlockThree() (int32, error)
	AddBlockFour() (int32, error)
	AddBlockFive() (int32, error)
}
