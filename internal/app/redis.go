package app

import (
	"context"
	"fmt"

	"github.com/canonflow/backend-starter/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func NewRedis(log *zap.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Get[string](config.RedisHost), config.Get[string](config.RedisPort)),
		Password: config.Get[string](config.RedisPassword),
		DB:       0,
		Protocol: 2,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:",
			zap.String("error", err.Error()),
		)
	}

	return rdb
}
