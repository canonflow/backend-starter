package app

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger = createLogger()

	defer logger.Sync()
}
