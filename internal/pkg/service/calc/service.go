package calc

import (
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// TODO Тип вывода Рзаб на ВПД соотвутствует типу Рзаб на глубине замера
// TODO Считает полностью все данные, не обрезая по времени
// единицы измерения должны быть идентичны выбранным ранее единицам измерения давления при импорте Рзаб на глубине
// замера. Т.е. формула для пересчета должна учитывать разные сценарии (пересчетные коэффициенты) в зависимости от ранее заданных при импорте значений Рзаб (атм,
// кгс/см2, бар)
// CalcTableOne

const g = 9.80665 // м/с²

// Service отвечает за логику расчетов данных в моделях.
type Service struct{}

// NewService создает новый экземпляр сервис импорта.
func NewService() *Service {
	return &Service{}
}

// CalcTableOne применяет гидростатику к одной записи, возвращая с заполненным PressureVPD.
func (s *Service) CalcTableOne(rec models.TableOne, cfg models.OperationConfig) models.TableOne {
	// 1) переводим измеренное давление в Па
	p0 := toPa(rec.PressureDepth, cfg.PressureUnit)

	// 2) выбираем плотность по времени
	var rho float64
	t := rec.Timestamp
	if !t.Before(cfg.WorkStart) && t.Before(cfg.WorkEnd) {
		rho = cfg.WorkDensity
	} else if !t.Before(cfg.IdleStart) && t.Before(cfg.IdleEnd) {
		rho = cfg.IdleDensity
	} else {
		// не попадает ни в один период
		return rec
	}

	// 3) гидростатическое приращение ΔP = ρ·g·Δh
	deltaPa := rho * g * cfg.DepthDiff

	// 4) итог в Па и обратно
	pVpd := p0 + deltaPa
	v := fromPa(pVpd, cfg.PressureUnit)
	rec.PressureVPD = &v
	return rec
}

func (s *Service) CalcBlockThree(table []models.TableThree) []models.TableThree {
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

// конвертация в Паскали и обратно
func toPa(p float64, unit string) float64 {
	switch unit {
	case "kgf/cm2":
		return p * 98066.5
	case "bar":
		return p * 1e5
	case "atm":
		return p * 101325
	default:
		return p
	}
}
func fromPa(pa float64, unit string) float64 {
	switch unit {
	case "kgf/cm2":
		return pa / 98066.5
	case "bar":
		return pa / 1e5
	case "atm":
		return pa / 101325
	default:
		return pa
	}
}
