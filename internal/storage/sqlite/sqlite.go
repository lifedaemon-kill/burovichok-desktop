package sqlite

import (
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/config"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	db *sqlx.DB
}

func NewDB(conf config.DBConf) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", conf.DSN)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlite.NewDB")
	}
	err = db.Ping()

	return db, err
}

type GuidebooksStorage interface {
	AddOilField([]models.OilField) error
	AddInstrumentType([]models.InstrumentType) error
	AddProductiveHorizon([]models.ProductiveHorizon) error

	GetAllOilField() ([]models.OilField, error)
	GetAllInstrumentType() ([]models.InstrumentType, error)
	GetAllProductiveHorizon() ([]models.ProductiveHorizon, error)
}

type BlocksStorage interface {
	GetAllTableFive() ([]models.TableFive, error)
	AddBlockFive(data models.TableFive) (reportID int64, err error)
}
