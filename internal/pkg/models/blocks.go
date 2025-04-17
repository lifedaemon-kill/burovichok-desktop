package models

import "time"

//ФУНКЦИОНАЛЬНЫЕ БЛОКИ

// TableOne соответствует блоку 1: "Загрузка Рзаб и Тзаб".
type TableOne struct {
	Timestamp     time.Time `xlsx:"Дата, время"`                     // колонка A
	PressureDepth float64   `xlsx:"Рзаб на глубине замера, кгс/см2"` // колонка B
	Temperature   float64   `xlsx:"Tзаб на глубине замера, оС"`      // колонка C
	PressureVPD   *float64  //
}

// TableTwo соответствует блоку 2: "Загрузка Ртр, Рзтр, Рлин".
// Каждая запись содержит три замера давления с их временными метками.
type TableTwo struct {
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

// TableThree соответствует блоку 3: "Данные по дебитам".
type TableThree struct {
	Timestamp     time.Time `json:"timestamp"`
	FlowLiquid    float64   `json:"flow_liquid"` // Qж, м3/сут
	WaterCut      float64   `json:"water_cut"`   // W, %
	FlowGas       float64   `json:"flow_gas"`    // Qг, тыс. м3/сут
	OilFlowRate   *float64  // Qн, m/сут Дебит нефти
	WaterFlowRate *float64  // Qв, m/сут Дебит воды
	GasToOilRatio *float64  // ГФ, м3/m Газовый фактор

}

// Inclinometry BlockFour соответствует блоку 4: Инклинометрия.
type Inclinometry struct {
	MeasuredDepth           float64  // Метры, Глубина спуска прибора по стволу (MD)
	TrueVerticalDepth       *float64 // Метры, Глубина спуска прибора по вертикали (TVD)
	TrueVerticalDepthSubSea *float64 // Метры, Абсолютная отметка (TVDSS)
}

//ИНФОРМАЦИОННЫЙ БЛОК

// GeneralInformation BlockFive Общие сведения об исследовании.
type GeneralInformation struct {
	OilPlaces                  OilPlaces         // Месторождение
	WellNumber                 int               // № скважины
	ClusterSiteNumber          int               // № кустовой площадки
	ProductiveHorizon          ProductiveHorizon // Продуктивный горизонт, пласт
	StartDate                  time.Time         // Дата начала исследования
	EndDate                    time.Time         // Дата окончания исследования
	InstrumentType             InstrumentType    // Тип прибора
	InstrumentNumber           int               // № прибора
	Inclinometry               Inclinometry      // Инклинометрия
	VDPMeasuredDepth           float64           // Метры, Верхние дыры перфорации по стволу (MD)
	VDPTrueVerticalDepth       *float64          // Метры, Верхние дыры перфорации по вертикали (TVD)
	VDPTrueVerticalDepthSubSea *float64          // Метры, Верхние дыры перфорации (ВДП) абсолютная отметка (TVDSS)
	DifferenceInstrumentAndVDP *float64          // Метры, Разница между прибором и ВДП по абсолютным отметкам, м
	DensityOil                 float64           // кгм/м3, Плотность для пересчета дебита нефти
	DensityLiquidStopped       float64           // кгм/м3, Плотность жидкости для пересчета давления на ВДП в остановленной скважине
	DensityLiquidWorking       float64           // кгм/м3, Плотность жидкости для пересчета давления на ВДП в работающей скважине
	PressureDifferenceStopped  *float64          // Единицы, выбранные при импорте Рзаб, Разница между давлением на глубине замера и ВДП в остановленной скважине
	PressureDifferenceWorking  *float64          // Единицы, выбранные при импорте Рзаб, Разница между давлением на глубине замера и ВДП в работающей скважине
}

//СПРАВОЧНИКИ

// ProductiveHorizon BlockSix Продуктивный горизонт, пласт.
type ProductiveHorizon string // Б1 Б2 Б3...

// OilPlaces BlockSeven Наименование месторождений
type OilPlaces string // Куюмбинское Юрубчено-Тохомское Ванкорское.

// InstrumentType BlockEight Тип приборов для замеров давления и температуры.
type InstrumentType string // ГС-АМТС PPS 25 КАМА-2
