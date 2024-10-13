package logger_test

import (
	"github.com/AxMdv/go-url-shortener/pkg/logger"
	"go.uber.org/zap"
)

func ExampleInitLogger() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}
	logger.Log.Info("Example message:",
		zap.String("Key 1", "Value 1"),
		zap.String("Key 2", "Value 2"),
	)
}
