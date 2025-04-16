package logger

import (
	"go.uber.org/zap"
)

// Logger является интерфейсом для логирования.
type Logger interface {
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

// Wrapper оборачивает *zap.SugaredLogger.
type Wrapper struct {
	*zap.SugaredLogger
}

// NewLogger создает новый экземпляр логгера в зависимости от окружения.
func NewLogger(env string) (Logger, error) {
	var baseLogger *zap.Logger
	var err error
	switch env {
	case "prod":
		baseLogger, err = zap.NewProduction()
	default:
		baseLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}
	return &Wrapper{baseLogger.Sugar()}, nil
}
