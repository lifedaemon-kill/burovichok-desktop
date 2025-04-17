package models

import "time"

//ФУНКЦИОНАЛЬНЫЕ БЛОКИ

// BlockOne модель для файла "Рзаб и Тзаб".
type BlockOne struct {
	Timestamp   time.Time `xlsx:"Дата, время"`                     // колонка A
	Pressure    float64   `xlsx:"Рзаб на глубине замера, кгс/см2"` // колонка B
	Temperature float64   `xlsx:"Tзаб на глубине замера, оС"`      // колонка C
}

// BlockTwo соответствует блокам "Загрузка Ртр, Рзтр, Рлин".
// Каждая запись содержит три замера давления с их временными метками.
type BlockTwo struct {
	// Замер давления в трубном пространстве
	TimestampTubing time.Time `json:"timestamp_tubing"` // из колонки A: «Дата, время»
	PressureTubing  float64   `json:"pressure_tubing"`  // из колонки B: «Ртр, кгс/см2»

	// Замер давления в затрубном пространстве
	TimestampAnnulus time.Time `json:"timestamp_annulus"` // из колонки C: «Дата, время»
	PressureAnnulus  float64   `json:"pressure_annulus"`  // из колонки D: «Рзтр, кгс/см2»

	// Замер линейного давления
	TimestampLinear time.Time `json:"timestamp_linear"` // из колонки E: «Дата, время»
	PressureLinear  float64   `json:"pressure_linear"`  // из колонки F: «Рлин, кгс/см2»
}

// BlockThree соответствует блоку 3: "Данные по дебитам".
type BlockThree struct {
	Timestamp  time.Time `json:"timestamp"`
	FlowLiquid float64   `json:"flow_liquid"` // Qж, м3/сут
	WaterCut   float64   `json:"water_cut"`   // W, %
	FlowGas    float64   `json:"flow_gas"`    // Qг, тыс. м3/сут
}

// BlockFour Инклинометрия
type BlockFour struct {
	DepthMeasured float64 // Метры, Глубина спуска прибора по стволу (MD)
	DepthVertical float64 // Метры, Глубина спуска прибора по вертикали (TVD)
	AbsoluteDepth float64 // Метры, Абсолютная отметка (TVDSS)
}

//ИНФОРМАЦИОННЫЙ БЛОК

// BlockFive Общие сведения об исследовании
type BlockFive struct {
	FieldName                           // Месторождение
	WellNumber                int       // № скважины
	ClusterSiteNumber         int       // № кустовой площадки
	ProductiveHorizon                   // Продуктивный горизонт, пласт
	StartDate                 time.Time // Дата начала исследования
	EndDate                   time.Time // Дата окончания исследования
	InstrumentType                      // Тип прибора
	InstrumentNumber          int       // № прибора
	BlockFour                           // Инклинометрия
	PerforationDepthMeasured  float64   // Метры, Верхние дыры перфорации по стволу (MD)
	PerforationDepthVertical  float64   // Метры, Верхние дыры перфорации по вертикали (TVD)
	DepthDifference           float64   // Метры, Разница между прибором и ВДП по абсолютным отметкам
	DensityOil                float64   // кгм/м3, Плотность для пересчета дебита нефти
	DensityLiquidStopped      float64   // кгм/м3, Плотность жидкости для пересчета давления на ВДП в остановленной скважине
	DensityLiquidWorking      float64   // кгм/м3, Плотность жидкости для пересчета давления на ВДП в работающей скважине
	PressureDifferenceStopped float64   // Единицы, выбранные при импорте Рзаб, Разница между давлением на глубине замера и ВДП в остановленной скважине
	PressureDifferenceWorking float64   // Единицы, выбранные при импорте Рзаб, Разница между давлением на глубине замера и ВДП в работающей скважине
}

//СПРАВОЧНИКИ

// ProductiveHorizon BlockSix Продуктивный горизонт, пласт
// Б1 Б2 Б3...
type ProductiveHorizon string

// FieldName BlockSeven Наименование месторождений
// Куюмбинское Юрубчено-Тохомское Ванкорское
type FieldName string

// InstrumentType BlockEight Тип приборов для замеров давления и температуры
// ГС-АМТС PPS 25 КАМА-2
type InstrumentType string
