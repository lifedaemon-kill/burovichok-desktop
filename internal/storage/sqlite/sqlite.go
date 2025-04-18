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
	err = db.Ping()

	return db, err
}

type GuidebooksStorage interface {
	AddOilPlaces([]models.OilField) error
	AddInstrumentType([]models.InstrumentType) error
	AddProductiveHorizon([]models.ProductiveHorizon) error

	GetAllOilPlaces() ([]models.OilField, error)
	GetAllInstrumentType() ([]models.InstrumentType, error)
	GetAllProductiveHorizon() ([]models.ProductiveHorizon, error)
}

type BlocksStorage interface {
	GetTableOneByReportID() ([]models.TableOne, error)
	GetTableTwoByReportID() ([]models.TableTwo, error)
	GetTableThreeByReportID() ([]models.TableThree, error)
	GetTableFourByReportID() ([]models.TableFour, error)
	GetTableFive() ([]models.TableFive, error)

	AddBlockOne(data []models.TableOne, reportID int32) error
	AddBlockTwo(data []models.TableTwo, reportID int32) error
	AddBlockThree(data []models.TableThree, reportID int32) error
	AddBlockFour(data []models.TableFour, reportID int32) error
	AddBlockFive(data []models.TableFive) (reportID int32, err error)
}
