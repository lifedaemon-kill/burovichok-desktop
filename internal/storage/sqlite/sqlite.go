package sqlite

import (
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	DB *sqlx.DB
}

func NewDB(conf config.DBConf) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conf.DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlite.NewDB")
	}
	return db, err
}

type GuidebooksStorage interface {
	AddOilPlaces() ([]models.OilField, error)
	AddInstrumentType() ([]models.InstrumentType, error)
	AddProductiveHorizon() ([]models.ProductiveHorizon, error)

	GetAllOilPlaces() ([]models.OilField, error)
	GetAllInstrumentType() ([]models.InstrumentType, error)
	GetAllProductiveHorizon() ([]models.ProductiveHorizon, error)
}

type BlocksStorage interface {
	GetBlockOne() ([]models.TableOne, error)
	GetBlockTwo() ([]models.TableTwo, error)
	GetBlockThree() ([]models.TableThree, error)
	GetBlockFour() ([]models.TableFour, error)
	GetBlockFive() ([]models.TableFive, error)

	AddBlockOne() (int32, error)
	AddBlockTwo() (int32, error)
	AddBlockThree() (int32, error)
	AddBlockFour() (int32, error)
	AddBlockFive() (int32, error)
}
