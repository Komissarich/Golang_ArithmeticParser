package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func SetupLogger() *zap.Logger {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
	}

	return logger
}
