package models

import "time"

// TableOne — Блок 1. Загрузка забойного давления и температуры
// Поля соответствуют колонкам Excel через теги xlsx
type TableOne struct {
	Timestamp        time.Time `xlsx:"Дата, время"`
	PressureDepth    float64   `xlsx:"Рзаб на глубине замера, кгс/см2"`
	TemperatureDepth float64   `xlsx:"Tзаб на глубине замера, °C"`
	PressureAtVDP    *float64  // расчётное поле
}

// TableTwo — Блок 2. Замеры трубного, затрубного и линейного давления
type TableTwo struct {
	TimestampTubing  time.Time `xlsx:"Дата трубного замера, Дата, время"`
	PressureTubing   float64   `xlsx:"Ртр, кгс/см2"`
	TimestampAnnulus time.Time `xlsx:"Дата затрубного замера, Дата, время"`
	PressureAnnulus  float64   `xlsx:"Рзтр, кгс/см2"`
	TimestampLinear  time.Time `xlsx:"Дата линейного замера, Дата, время"`
	PressureLinear   float64   `xlsx:"Рлин, кгс/см2"`
}

// TableThree — Блок 3. Дебиты жидкости, воды, газа и расчётные поля
type TableThree struct {
	Timestamp     time.Time `xlsx:"Дата, время"`
	FlowLiquid    float64   `xlsx:"Qж, м3/сут"`
	WaterCut      float64   `xlsx:"W, %"`
	FlowGas       float64   `xlsx:"Qг, тыс.м3/сут"`
	OilFlowRate   *float64  // Qн, расчётное поле
	WaterFlowRate *float64  // Qв, расчётное поле
	GasOilRatio   *float64  // ГФ, расчётное поле
}

// TableFour — Блок 4. Инклинометрия (MD, TVD, TVDSS)
type TableFour struct {
	MeasuredDepth           float64 `xlsx:"Глубина по стволу, м" db:"measure_depth"`                // MD
	TrueVerticalDepth       float64 `xlsx:"Глубина по вертикали, м" db:"true_vertical_depth"`       // TVD
	TrueVerticalDepthSubSea float64 `xlsx:"Абсолютная глубина, м" db:"true_vertical_depth_sub_sea"` // TVDSS
}

// TableFive — Блок 5. Общие сведения об исследовании
// Поля берутся из формы или справочников, поэтому xlsx-тегов нет
type TableFive struct {
	FieldName               string    `db:"field_name"`                  // Месторождение
	FieldNumber             int       `db:"field_number"`                // № скважины
	ClusterNumber           *int      `db:"cluster_number"`              // № кустовой площадки (может быть не у всех)
	Horizon                 string    `db:"horizon"`                     // Продуктивный горизонт
	StartTime               time.Time `db:"start_time"`                  // Дата начала исследования
	EndTime                 time.Time `db:"end_time"`                    // Дата окончания исследования
	InstrumentType          string    `db:"instrument_type"`             // Тип прибора
	InstrumentNumber        *int      `db:"instrument_number"`           // № прибора (может быть не у всех)
	MeasuredDepth           float64   `db:"measure_depth"`               // MD
	TrueVerticalDepth       float64   `db:"true_vertical_depth"`         // TVD
	TrueVerticalDepthSubSea float64   `db:"true_vertical_depth_sub_sea"` // TVDSS // Данные инклинометрии
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

// Справочники

// ProductiveHorizon Б1, Б2, Б3...
type ProductiveHorizon struct {
	Name string `db:"name"`
}

// OilField Наименование месторождения
type OilField struct {
	Name string `db:"name"`
}

// InstrumentType Тип прибора, например, ГС-АМТС, PPS 25, КАМА-2
type InstrumentType struct {
	Name string `db:"name"`
}
