package utils

import (
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger initializes the global logger based on environment
func InitLogger() *zap.Logger {
	once.Do(func() {
		env := strings.ToLower(GetDotEnv("ENVIRONMENT"))

		var config zap.Config

		if env == "production" {
			// Production configuration
			config = zap.NewProductionConfig()
			config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		} else {
			// Development configuration
			config = zap.NewDevelopmentConfig()
			config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// Customize output paths
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}

		// Build the logger
		var err error
		logger, err = config.Build()
		if err != nil {
			panic("Failed to initialize logger: " + err.Error())
		}
	})

	return logger
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if logger == nil {
		return InitLogger()
	}
	return logger
}

// LogInfo logs an info message
func LogInfo(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// LogError logs an error message
func LogError(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// LogWarn logs a warning message
func LogWarn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// LogDebug logs a debug message
func LogDebug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// LogFatal logs a fatal message and exits
func LogFatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
