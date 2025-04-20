package calc

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/logger"
	"github.com/lifedaemon-kill/burovichok-desktop/internal/pkg/models"
)

// TODO Тип вывода Рзаб на ВПД соотвутствует типу Рзаб на глубине замера
// TODO Считает полностью все данные, не обрезая по времени
// единицы измерения должны быть идентичны выбранным ранее единицам измерения давления при импорте Рзаб на глубине
// замера. Т.е. формула для пересчета должна учитывать разные сценарии (пересчетные коэффициенты) в зависимости от ранее заданных при импорте значений Рзаб (атм,
// кгс/см2, бар)
// CalcTableOne

const g = 9.80665 // м/с²

type DataLoader interface {
	GetBlockFourByResearchID(ctx context.Context, researchID uuid.UUID) ([]models.TableFour, error)
}

// Service отвечает за логику расчетов данных в моделях.
type Service struct {
	dataLoader DataLoader
	logger     logger.Logger
}

// NewService создает новый экземпляр сервис импорта.
func NewService(dataLoader DataLoader, logger logger.Logger) *Service {
	return &Service{
		dataLoader: dataLoader,
		logger:     logger,
	}
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
	rec.PressureAtVDP = &v
	return rec
}

// CalcBlockThree рассчитывает дебиты нефти (Qн), воды (Qв) и газовый фактор (ГФ)
// на основании входных данных TableThree:
//
//	FlowLiquid — общий дебит жидкости Qж, м³/сут
//	WaterCut    — обводненность W, %
//	FlowGas     — дебит газа Qг, тыс. м³/сут
func (s *Service) CalcBlockThree(tbl models.TableThree) models.TableThree {
	// 1) вычисляем дебит воды Qв = Qж * W/100
	waterRate := tbl.FlowLiquid * tbl.WaterCut / 100.0

	// 2) вычисляем дебит нефти Qн = Qж – Qв
	oilRate := tbl.FlowLiquid - waterRate

	// 3) вычисляем газовый фактор ГФ = (Qг * 1000) / Qн
	//    (умножаем Qг на 1000, чтобы перевести из тыс. м³/сут в м³/сут)
	var gfPtr *float64
	if oilRate > 0 {
		gf := (tbl.FlowGas * 1000.0) / oilRate
		gfPtr = &gf
	}

	// 4) сохраняем результаты в исходную структуру
	tbl.WaterFlowRate = &waterRate
	tbl.OilFlowRate = &oilRate
	tbl.GasOilRatio = gfPtr

	return tbl
}

// CalcBlockFive calculates automatic fields for TableFive using Block 4 survey data.
func (s *Service) CalcBlockFive(ctx context.Context, tbl models.TableFive, researchID uuid.UUID) models.TableFive {
	// 1. Получаем данные инклинометрии (Block 4) для скважины
	survey, err := s.dataLoader.GetBlockFourByResearchID(ctx, researchID)
	if err != nil {
		s.logger.Errorw("CalcBlockFive: failed to get block four data", "error", err)
		return tbl
	}

	// 2. Вычисляем TVD и TVDSS для прибора по MD
	tvd, tvdss := interpolateTVD(survey, tbl.MeasuredDepth)
	tbl.TrueVerticalDepth = lo.ToPtr(tvd)
	tbl.TrueVerticalDepthSubSea = lo.ToPtr(tvdss)

	// 3. Если задана MD перфорации (VDP), рассчитываем ее отметки
	if tbl.VDPMeasuredDepth > 0 {
		vdpTVD, vdpTVDSS := interpolateTVD(survey, tbl.VDPMeasuredDepth)
		tbl.VDPTrueVerticalDepth = &vdpTVD
		tbl.VDPTrueVerticalDepthSea = &vdpTVDSS

		// 4. Разница между прибором и ВДП по абсолютным отметкам (TVDSS)
		if tbl.TrueVerticalDepthSubSea != nil {
			delta := *tbl.TrueVerticalDepthSubSea - vdpTVDSS
			tbl.DiffInstrumentVDP = &delta
		}

		// 5. Гидростатическое давление (ΔP = ρ * g * Δh)
		//    g = 9.81 m/s²
		const g = 9.81
		var heightDiff float64
		if tbl.TrueVerticalDepthSubSea != nil {
			heightDiff = *tbl.TrueVerticalDepthSubSea - vdpTVDSS
		}
		pStopped := tbl.DensityLiquidStopped * g * heightDiff
		pWorking := tbl.DensityLiquidWorking * g * heightDiff
		tbl.PressureDiffStopped = &pStopped
		tbl.PressureDiffWorking = &pWorking
	}

	return tbl
}

// interpolateTVD performs linear interpolation on survey data to find TVD and TVDSS at a given MD.
func interpolateTVD(survey []models.TableFour, md float64) (tvd, tvdss float64) {
	if len(survey) == 0 {
		return 0, 0
	}

	// Если MD меньше минимального, возвращаем первую точку
	if md <= survey[0].MeasuredDepth {
		return survey[0].TrueVerticalDepth, survey[0].TrueVerticalDepthSubSea
	}

	// Ищем два соседних замера для интерполяции
	for i := 1; i < len(survey); i++ {
		prev := survey[i-1]
		curr := survey[i]
		if md <= curr.MeasuredDepth {
			ratio := (md - prev.MeasuredDepth) / (curr.MeasuredDepth - prev.MeasuredDepth)
			tvd = prev.TrueVerticalDepth + ratio*(curr.TrueVerticalDepth-prev.TrueVerticalDepth)
			tvdss = prev.TrueVerticalDepthSubSea + ratio*(curr.TrueVerticalDepthSubSea-prev.TrueVerticalDepthSubSea)
			return tvd, tvdss
		}
	}

	// Если MD больше максимального, возвращаем последнюю точку
	last := survey[len(survey)-1]
	return last.TrueVerticalDepth, last.TrueVerticalDepthSubSea
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
