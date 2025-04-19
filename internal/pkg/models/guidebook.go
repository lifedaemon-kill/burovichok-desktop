package models

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

// ResearchType Тип исследования
type ResearchType struct {
	Name string `db:"name"`
}

// TableName возвращает имя таблицы для ProductiveHorizon
func (ProductiveHorizon) TableName() string {
	return "productive_horizon"
}

// Columns возвращает список колонок для ProductiveHorizon
func (ProductiveHorizon) Columns() []string {
	return []string{"name"}
}

// Map конвертирует ProductiveHorizon в map[column]value
func (p ProductiveHorizon) Map() map[string]interface{} {
	return map[string]interface{}{"name": p.Name}
}

// TableName возвращает имя таблицы для OilField
func (OilField) TableName() string {
	return "oilfield"
}

// Columns возвращает список колонок для OilField
func (OilField) Columns() []string {
	return []string{"name"}
}

// Map конвертирует OilField в map[column]value
func (o OilField) Map() map[string]interface{} {
	return map[string]interface{}{"name": o.Name}
}

// TableName возвращает имя таблицы для InstrumentType
func (InstrumentType) TableName() string {
	return "instrument_type"
}

// Columns возвращает список колонок для InstrumentType
func (InstrumentType) Columns() []string {
	return []string{"name"}
}

// Map конвертирует InstrumentType в map[column]value
func (i InstrumentType) Map() map[string]interface{} {
	return map[string]interface{}{"name": i.Name}
}

// TableName возвращает имя таблицы для ResearchType
func (ResearchType) TableName() string {
	return "research_type"
}

// Columns возвращает список колонок для ResearchType
func (ResearchType) Columns() []string {
	return []string{"name"}
}

// Map конвертирует ResearchType в map[column]value
func (r ResearchType) Map() map[string]interface{} {
	return map[string]interface{}{"name": r.Name}
}
