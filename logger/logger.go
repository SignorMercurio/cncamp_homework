package logger

import (
	"os"

	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "httpserver.log" // for testing
	}
	cfg.OutputPaths = []string{
		"stdout",
		logFile,
	}

	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		cfg.Level.SetLevel(zap.DebugLevel)
	case "WARNING":
		cfg.Level.SetLevel(zap.WarnLevel)
	case "ERROR":
		cfg.Level.SetLevel(zap.ErrorLevel)
	default:
		cfg.Level.SetLevel(zap.InfoLevel)
	}
	return cfg.Build()
}
