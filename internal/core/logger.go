package core

import (
	"os"
	"strings"

	"github.com/canonflow/backend-starter/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Used inside init function for app package
func CreateLogger() {
	if logger == nil {
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		config := zap.Config{
			Level:             zap.NewAtomicLevelAt(getLogLevelFromEnv()),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths: []string{
				"stdout",
			},
			ErrorOutputPaths: []string{
				"stdout",
			},
			InitialFields: map[string]interface{}{
				"pid": os.Getpid(),
			},
		}

		log := zap.Must(config.Build())

		defer log.Sync()

		logger = log
	}
}

func getLogLevelFromEnv() zapcore.Level {
	levelStr := strings.ToLower(config.GetOrDefault(config.LogLevel, "info"))
	switch levelStr {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func GetLogger() *zap.Logger {
	return logger
}
