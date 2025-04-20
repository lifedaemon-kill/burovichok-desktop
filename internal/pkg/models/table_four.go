package models

import "github.com/google/uuid"

// TableFour — Блок 4. Инклинометрия (MD, TVD, TVDSS)
type TableFour struct {
	ResearchID              uuid.UUID `db:"research_id"`
	MeasuredDepth           float64   `xlsx:"Глубина по стволу, м" db:"measure_depth"`                // MD
	TrueVerticalDepth       float64   `xlsx:"Глубина по вертикали, м" db:"true_vertical_depth"`       // TVD
	TrueVerticalDepthSubSea float64   `xlsx:"Абсолютная глубина, м" db:"true_vertical_depth_sub_sea"` // TVDSS
}

// TableName возвращает имя таблицы в БД для TableFour
func (TableFour) TableName() string {
	return "table_four"
}

// Columns возвращает список колонок в той же последовательности,
// в которой они используются при INSERT/SELECT
func (TableFour) Columns() []string {
	return []string{
		"research_id",
		"measure_depth",
		"true_vertical_depth",
		"true_vertical_depth_sub_sea",
	}
}

// Map конвертирует структуру TableFour в map[column]value,
// удобный для NamedExec или Squirrel
func (t TableFour) Map() map[string]interface{} {
	return map[string]interface{}{
		"research_id":                 t.ResearchID,
		"measure_depth":               t.MeasuredDepth,
		"true_vertical_depth":         t.TrueVerticalDepth,
		"true_vertical_depth_sub_sea": t.TrueVerticalDepthSubSea,
	}
}
