package database

import "github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"

type exporter interface {
	SaveReport(
		tableOne []models.TableOne,
		tableTwo []models.TableTwo,
		tableThree []models.TableThree,
		tableFour []models.TableFour,
		tableFive []models.TableFive) error

	SaveInstrumentType([]models.InstrumentType) error
	SaveOilField([]models.OilField) error
	SaveProductiveHorizon([]models.ProductiveHorizon) error
}

func (d dbService) SaveReport(tableOne []models.TableOne, tableTwo []models.TableTwo, tableThree []models.TableThree, tableFour []models.TableFour, tableFive []models.TableFive) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) SaveInstrumentType(types []models.InstrumentType) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) SaveOilField(fields []models.OilField) error {
	//TODO implement me
	panic("implement me")
}

func (d dbService) SaveProductiveHorizon(horizons []models.ProductiveHorizon) error {
	//TODO implement me
	panic("implement me")
}
