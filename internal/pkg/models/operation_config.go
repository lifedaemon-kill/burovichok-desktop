package models

import "time"

// OperationConfig хранит параметры гидростатики и контекст парсинга блока 1.
type OperationConfig struct {
	PressureUnit string  // "kgf/cm2", "bar" или "atm"
	DepthDiff    float64 // Δh между замером и ВДП (в метрах)

	WorkStart   time.Time // начало рабочего периода
	WorkEnd     time.Time // конец рабочего периода
	WorkDensity float64   // плотность в рабочем периоде

	IdleStart   time.Time // начало простоя
	IdleEnd     time.Time // конец простоя
	IdleDensity float64   // плотность при простое
}
