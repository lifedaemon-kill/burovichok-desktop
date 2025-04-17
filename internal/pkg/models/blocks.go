package models

import "time"

// BlockOne — Блок 1. Загрузка забойного давления и температуры
// Поля соответствуют колонкам Excel через теги xlsx
type BlockOne struct {
	Timestamp        time.Time `xlsx:"Дата, время"`
	PressureDepth    float64   `xlsx:"Рзаб на глубине замера, кгс/см2"`
	TemperatureDepth float64   `xlsx:"Tзаб на глубине замера, °C"`
	PressureAtVDP    *float64  // расчётное поле
}

// BlockTwo — Блок 2. Замеры трубного, затрубного и линейного давления
type BlockTwo struct {
	TimestampTubing  time.Time `xlsx:"Дата трубного замера, Дата, время"`
	PressureTubing   float64   `xlsx:"Ртр, кгс/см2"`
	TimestampAnnulus time.Time `xlsx:"Дата затрубного замера, Дата, время"`
	PressureAnnulus  float64   `xlsx:"Рзтр, кгс/см2"`
	TimestampLinear  time.Time `xlsx:"Дата линейного замера, Дата, время"`
	PressureLinear   float64   `xlsx:"Рлин, кгс/см2"`
}

// BlockThree — Блок 3. Дебиты жидкости, воды, газа и расчётные поля
type BlockThree struct {
	Timestamp     time.Time `xlsx:"Дата, время"`
	FlowLiquid    float64   `xlsx:"Qж, м3/сут"`
	WaterCut      float64   `xlsx:"W, %"`
	FlowGas       float64   `xlsx:"Qг, тыс.м3/сут"`
	OilFlowRate   *float64  // Qн, расчётное поле
	WaterFlowRate *float64  // Qв, расчётное поле
	GasOilRatio   *float64  // ГФ, расчётное поле
}

// BlockFour — Блок 4. Инклинометрия (MD, TVD, TVDSS)
type BlockFour struct {
	MeasuredDepth           float64 `xlsx:"Глубина по стволу, м"`    // MD
	TrueVerticalDepth       float64 `xlsx:"Глубина по вертикали, м"` // TVD
	TrueVerticalDepthSubSea float64 `xlsx:"Абсолютная глубина, м"`   // TVDSS
}

// BlockFive — Блок 5. Общие сведения об исследовании
// Поля берутся из формы или справочников, поэтому xlsx-тегов нет
type BlockFive struct {
	Field                   OilField          // Месторождение
	WellNumber              int               // № скважины
	ClusterNumber           int               // № кустовой площадки
	Horizon                 ProductiveHorizon // Продуктивный горизонт
	StartTime               time.Time         // Дата начала исследования
	EndTime                 time.Time         // Дата окончания исследования
	InstrumentType          InstrumentType    // Тип прибора
	InstrumentNumber        int               // № прибора
	Inclinometry            BlockFour         // Данные инклинометрии
	VDPMeasuredDepth        float64           // MD ВДП
	VDPTrueVerticalDepth    *float64          // TVD ВДП, расчётное
	VDPTrueVerticalDepthSea *float64          // TVDSS ВДП, расчётное
	DiffInstrumentVDP       *float64          // Разница отметок, расчётное
	DensityOil              float64           // Плотность для дебита нефти, кг/м3
	DensityLiquidStopped    float64           // Плотность жидкости в простое, кг/м3
	DensityLiquidWorking    float64           // Плотность жидкости в работе, кг/м3
	PressureDiffStopped     *float64          // ΔP простоя, расчётное
	PressureDiffWorking     *float64          // ΔP работы, расчётное
}

// Справочники

type ProductiveHorizon string // Б1, Б2, Б3...
type OilField string          // Наименование месторождения
type InstrumentType string    // Тип прибора, например, ГС-АМТС, PPS 25, КАМА-2
