package export

import (
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
