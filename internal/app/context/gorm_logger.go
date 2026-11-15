package context

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	gormUtils "gorm.io/gorm/utils"
)

type zapWriter struct {
	Logger *zap.Logger
}

func (l *zapWriter) Printf(message string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf("[GORM TRACE] "+message, args...))
}

func (l *zapWriter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	rowsStr := "-"

	if rows != -1 {
		rowsStr = strconv.FormatInt(rows, 10)
	}

	l.Logger.Debug("[GORM TRACE]",
		zap.String("file", gormUtils.FileWithLineNum()),
		zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1e6),
		zap.String("rows", rowsStr),
		zap.String("sql", sql),
	)
}

func NewGormLogger() gormLogger.Interface {
	logger := gormLogger.New(&zapWriter{
		Logger: GetLogger().With(
			zap.String("system", "gorm"),
		),
	}, gormLogger.Config{
		SlowThreshold:             time.Second * 5,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      true,
		LogLevel:                  gormLogger.Info,
	})

	return logger
}
