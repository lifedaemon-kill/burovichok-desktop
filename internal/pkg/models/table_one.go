package models

import "time"

// TableOne — Блок 1. Загрузка забойного давления и температуры
// Поля соответствуют колонкам Excel через теги xlsx
type TableOne struct {
	Timestamp        time.Time `xlsx:"Дата, время"`                     // метка времени
	PressureDepth    float64   `xlsx:"Рзаб на глубине замера, кгс/см2"` // забойное давление
	TemperatureDepth float64   `xlsx:"Tзаб на глубине замера, °C"`      // забойная температура
	PressureAtVDP    float64   // расчётное поле
}

// TableName возвращает имя таблицы в БД для TableOne
func (TableOne) TableName() string {
	return "table_one"
}

// Columns возвращает список колонок в той же последовательности,
// в которой они должны использоваться при INSERT/SELECT
func (TableOne) Columns() []string {
	return []string{
		"timestamp",
		"pressure_depth",
		"temperature_depth",
		"pressure_at_vdp",
	}
}

// Map конвертирует TableOne в map[column]value, пригодный для NamedExec или Squirrel
func (t TableOne) Map() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":         t.Timestamp,
		"pressure_depth":    t.PressureDepth,
		"temperature_depth": t.TemperatureDepth,
		"pressure_at_vdp":   t.PressureAtVDP,
	}
}
