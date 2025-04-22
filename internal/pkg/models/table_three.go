package models

import "time"

// TableThree — Блок 3. Дебиты жидкости, воды, газа и расчётные поля
type TableThree struct {
	Timestamp      time.Time `xlsx:"Дата, время"`
	LiquidFlowRate float64   `xlsx:"Qж, м3/сут"`
	WaterCut       float64   `xlsx:"W, %"`
	GasFlowRate    float64   `xlsx:"Qг, тыс.м3/сут"`
	OilFlowRate    *float64  // Qн, расчётное поле
	WaterFlowRate  *float64  // Qв, расчётное поле
	GasFactor      *float64  // ГФ, расчётное поле
}

// TableName возвращает имя таблицы в БД для TableThree
func (TableThree) TableName() string {
	return "table_three"
}

// Columns возвращает список колонок в той же последовательности,
// в которой они используются при INSERT/SELECT
func (TableThree) Columns() []string {
	return []string{
		"timestamp",
		"flow_liquid",
		"water_cut",
		"flow_gas",
		"oil_flow_rate",
		"water_flow_rate",
		"gas_oil_ratio",
	}
}

// Map конвертирует структуру TableThree в map[column]value,
// удобный для NamedExec или Squirrel
func (t TableThree) Map() map[string]interface{} {
	return map[string]interface{}{
		"timestamp":       t.Timestamp,
		"flow_liquid":     t.LiquidFlowRate,
		"water_cut":       t.WaterCut,
		"flow_gas":        t.GasFlowRate,
		"oil_flow_rate":   t.OilFlowRate,
		"water_flow_rate": t.WaterFlowRate,
		"gas_oil_ratio":   t.GasFactor,
	}
}
