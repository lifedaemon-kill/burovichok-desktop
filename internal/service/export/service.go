package export

import (
	"github.com/jmoiron/sqlx"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/storage/sqlite"
)

type Service interface {
	SaveReport(
		tableOne []models.TableOne,
		tableTwo []models.TableTwo,
		tableThree []models.TableThree,
		tableFour []models.TableFour,
		tableFive []models.TableFive) error

	AddInstrumentType([]models.InstrumentType) error
	AddOilField([]models.OilField) error
	AddProductiveHorizon([]models.ProductiveHorizon) error
}
type dbService struct {
	g sqlite.GuidebooksStorage
	b sqlite.BlocksStorage
}

func NewExporter(db *sqlx.DB, zLog logger.Logger) (Service, error) {
	//Репозитории для работы с данными
	guidebooksRepository, err := sqlite.NewGuidebookStorage(db)
	if err != nil {
		zLog.Errorw("sqlite.NewGuidebookStorage", "error", err)
		return nil, err
	}
	blocksRepository, err := sqlite.NewBlockStorage(db)
	if err != nil {
		zLog.Errorw("sqlite.NewBlockStorage", "error", err)
		return nil, err
	}
	zLog.Infow("Repositories initialized successfully")

	return &dbService{
		g: guidebooksRepository,
		b: blocksRepository,
	}, nil
}

func (d dbService) SaveReport(tableOne []models.TableOne, tableTwo []models.TableTwo, tableThree []models.TableThree, tableFour []models.TableFour, tableFive []models.TableFive) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) AddInstrumentType(types []models.InstrumentType) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) AddOilField(fields []models.OilField) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) AddProductiveHorizon(horizons []models.ProductiveHorizon) error {
	//TODO implement me
	panic("implement me")
}
