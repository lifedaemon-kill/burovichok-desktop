package models

import "time"

// TableFive — Блок 5. Общие сведения об исследовании
// Поля берутся из формы или справочников, поэтому xlsx-тегов нет
type TableFive struct {
	ID                      int       `db:"id"`
	FieldName               string    `db:"field_name"`                  // Месторождение
	FieldNumber             int       `db:"field_number"`                // № скважины
	ClusterNumber           int       `db:"cluster_number"`              // № кустовой площадки (может быть не у всех)
	Horizon                 string    `db:"horizon"`                     // Продуктивный горизонт
	StartTime               time.Time `db:"start_time"`                  // Дата начала исследования
	EndTime                 time.Time `db:"end_time"`                    // Дата окончания исследования
	InstrumentType          string    `db:"instrument_type"`             // Тип прибора
	InstrumentNumber        int       `db:"instrument_number"`           // № прибора (может быть не у всех)
	MeasuredDepth           float64   `db:"measure_depth"`               // MD
	TrueVerticalDepth       *float64  `db:"true_vertical_depth"`         // TVD
	TrueVerticalDepthSubSea *float64  `db:"true_vertical_depth_sub_sea"` // TVDSS // Данные инклинометрии
	VDPMeasuredDepth        float64   `db:"vdp_measured_depth"`          // MD ВДП
	VDPTrueVerticalDepth    *float64  `db:"vdp_true_vertical_depth"`     // TVD ВДП, расчётное
	VDPTrueVerticalDepthSea *float64  `db:"vdp_true_vertical_depth_sea"` // TVDSS ВДП, расчётное
	DiffInstrumentVDP       *float64  `db:"diff_instrument_vdp"`         // Разница отметок, расчётное
	DensityOil              float64   `db:"density_oil"`                 // Плотность для дебита нефти, кг/м3
	DensityLiquidStopped    float64   `db:"density_liquid_stopped"`      // Плотность жидкости в простое, кг/м3
	DensityLiquidWorking    float64   `db:"density_liquid_working"`      // Плотность жидкости в работе, кг/м3
	PressureDiffStopped     *float64  `db:"pressure_diff_stopped"`       // ΔP простоя, расчётное
	PressureDiffWorking     *float64  `db:"pressure_diff_working"`       // ΔP работы, расчётное
}

// TableName имя таблицы.
func (TableFive) TableName() string {
	return "reports"
}

// Columns возвращает список имён колонок таблицы reports в том порядке,
// в котором они обычно используются для SELECT.
func (TableFive) Columns() []string {
	return []string{
		"id",
		"field_name",
		"field_number",
		"cluster_number",
		"horizon",
		"start_time",
		"end_time",
		"instrument_type",
		"instrument_number",
		"measure_depth",
		"true_vertical_depth",
		"true_vertical_depth_sub_sea",
		"vdp_measured_depth",
		"vdp_true_vertical_depth",
		"vdp_true_vertical_depth_sea",
		"diff_instrument_vdp",
		"density_oil",
		"density_liquid_stopped",
		"density_liquid_working",
		"pressure_diff_stopped",
		"pressure_diff_working",
	}
}

// Map конвертация TableFive в map[string]interface{}.
func (t TableFive) Map() map[string]any {
	return map[string]any{
		"field_name":                  t.FieldName,
		"field_number":                t.FieldNumber,
		"cluster_number":              t.ClusterNumber,
		"horizon":                     t.Horizon,
		"start_time":                  t.StartTime,
		"end_time":                    t.EndTime,
		"instrument_type":             t.InstrumentType,
		"instrument_number":           t.InstrumentNumber,
		"measure_depth":               t.MeasuredDepth,
		"true_vertical_depth":         t.TrueVerticalDepth,
		"true_vertical_depth_sub_sea": t.TrueVerticalDepthSubSea,
		"vdp_measured_depth":          t.VDPMeasuredDepth,
		"vdp_true_vertical_depth":     t.VDPTrueVerticalDepth,
		"vdp_true_vertical_depth_sea": t.VDPTrueVerticalDepthSea,
		"diff_instrument_vdp":         t.DiffInstrumentVDP,
		"density_oil":                 t.DensityOil,
		"density_liquid_stopped":      t.DensityLiquidStopped,
		"density_liquid_working":      t.DensityLiquidWorking,
		"pressure_diff_stopped":       t.PressureDiffStopped,
		"pressure_diff_working":       t.PressureDiffWorking,
	}
}
