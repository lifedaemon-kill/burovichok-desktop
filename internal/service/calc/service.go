package calc

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/models"
)

// TODO Тип вывода Рзаб на ВПД соотвутствует типу Рзаб на глубине замера
// TODO Считает полностью все данные, не обрезая по времени
// единицы измерения должны быть идентичны выбранным ранее единицам измерения давления при импорте Рзаб на глубине
// замера. Т.е. формула для пересчета должна учитывать разные сценарии (пересчетные коэффициенты) в зависимости от ранее заданных при импорте значений Рзаб (атм,
// кгс/см2, бар)
// CalcBlockOne

// Service отвечает за логику расчетов данных в моделях.
type Service struct{}

// NewService создает новый экземпляр сервис импорта.
func NewService() *Service {
	return &Service{}
}

func (s *Service) CalcBlockOne(table []models.BlockOne) []models.BlockOne {
	// for i, value := range table {
	// 	//TODO
	// 	table[i].PressureVPD = 0
	// }
	return nil
}

func (s *Service) CalcBlockThree(table []models.BlockThree) []models.BlockThree {
	// for i, value := range table {
	// 	//TODO
	//  table[i].OilFlowRate = 0
	// 	table[i].WaterFlowRate = 0
	// 	table[i].GasToOilRatio = 0
	// }
	return nil
}

func (s *Service) CalcBlockFive(table []models.GeneralInformation) []models.GeneralInformation {
	// for i, value := range table {
	// 	//TODO
	// 	table[i].TrueVerticalDepth = 0
	// 	table[i].TrueVerticalDepthSubSea = 0

	// 	table[i].VDPTrueVerticalDepth = 0
	// 	table[i].VDPTrueVerticalDepthSubSea = 0
	// 	table[i].DifferenceInstrumentAndVDP = 0
	// }
	return nil
}
