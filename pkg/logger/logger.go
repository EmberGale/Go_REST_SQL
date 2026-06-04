package logger

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New создаёт новый zap.SugaredLogger с указанным уровнем логирования
func New(level string) (*zap.SugaredLogger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		// Если уровень не распознан, используем info
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger.Sugar(), nil
}

// -----

func InitLogger(level string) *zap.SugaredLogger {
	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil
	logConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	logConfig.DisableStacktrace = true

	logConfig.Level = zap.NewAtomicLevelAt(DetermineLogLevel(level))

	logger, err := logConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	return logger.Sugar()
}

func DetermineLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
