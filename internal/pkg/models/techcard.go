package models

const (
	FileExtension_XSLX = iota
	FileExtension_CSV
)

const (
	ContentType_ZABOY_PRESSURE_TEMPERATURE = iota
	ContentType_PIPE_CASING_LINEAR_PRESSURE
	ContentType_DEBIT
	ContentType_INCLINOMETRY
)

const (
	ResponseStatus_OK = iota
	ResponseStatus_BAD_REQUSET
	ResponseStatus_INTERNAL_ERROR
)

// TODO Дописать единицы измерения