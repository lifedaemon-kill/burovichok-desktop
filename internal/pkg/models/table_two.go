package models

import "time"

// TableTwo — Блок 2. Замеры трубного, затрубного и линейного давления
type TableTwo struct {
	TimestampTubing  time.Time `xlsx:"Дата трубного замера, Дата, время"`
	PressureTubing   float64   `xlsx:"Ртр, кгс/см2"`
	TimestampAnnulus time.Time `xlsx:"Дата затрубного замера, Дата, время"`
	PressureAnnulus  float64   `xlsx:"Рзтр, кгс/см2"`
	TimestampLinear  time.Time `xlsx:"Дата линейного замера, Дата, время"`
	PressureLinear   float64   `xlsx:"Рлин, кгс/см2"`
}

// TableName возвращает имя таблицы в БД для TableTwo
func (TableTwo) TableName() string {
	return "table_two"
}

// Columns возвращает список колонок в той же последовательности,
// в которой они используются при INSERT/SELECT
func (TableTwo) Columns() []string {
	return []string{
		"timestamp_tubing",
		"pressure_tubing",
		"timestamp_annulus",
		"pressure_annulus",
		"timestamp_linear",
		"pressure_linear",
	}
}

// Map конвертирует структуру TableTwo в map[column]value,
// удобный для NamedExec или Squirrel
func (t TableTwo) Map() map[string]interface{} {
	return map[string]interface{}{
		"timestamp_tubing":  t.TimestampTubing,
		"pressure_tubing":   t.PressureTubing,
		"timestamp_annulus": t.TimestampAnnulus,
		"pressure_annulus":  t.PressureAnnulus,
		"timestamp_linear":  t.TimestampLinear,
		"pressure_linear":   t.PressureLinear,
	}
}
