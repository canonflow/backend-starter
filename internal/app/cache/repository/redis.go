package repository

import (
	"context"
	"fmt"
	"time"

	appContext "github.com/canonflow/backend-starter/internal/app/context"
	"github.com/canonflow/backend-starter/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type redisCache struct {
	rdb *redis.Client
}

func NewRedis() ICache {
	log := appContext.GetLogger().
		With(zap.String("__system__", "redis"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Get[string](config.RedisHost), config.Get[string](config.RedisPort)),
		Password: config.GetOrDefault(config.RedisPassword, ""),
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

	return &redisCache{
		rdb: rdb,
	}
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()

	return val, err
}

func (r *redisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	r.rdb.Set(ctx, key, value, ttl)

	return nil
}

func (r *redisCache) Del(ctx context.Context, key string) {
	r.rdb.Del(ctx, key)
}

func (r *redisCache) DelMany(ctx context.Context, keys ...string) {
	r.rdb.Del(ctx, keys...)
}
