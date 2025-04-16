package models

import "time"

// BlockOne модель для файла  Рзаб и Тзаб.xlsx.
type BlockOne struct {
	Timestamp   time.Time `xlsx:"Дата, время"`                     // колонка A
	Pressure    float64   `xlsx:"Рзаб на глубине замера, кгс/см2"` // колонка B
	Temperature float64   `xlsx:"Tзаб на глубине замера, оС"`      // колонка C
}

// BlockTwo соответствует блокам «Загрузка Ртр, Рзтр, Рлин».
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

type BlockThree struct {
}

type BlockFour struct {
}

type BlockFive struct {
}

type BlockSix struct {
}

type BlockSeven struct {
}
