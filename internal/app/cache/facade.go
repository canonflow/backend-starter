package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/canonflow/backend-starter/internal/app/cache/repository"
	appContext "github.com/canonflow/backend-starter/internal/app/context"
	"go.uber.org/zap"
)

func GetFromCache[T any](
	ctx context.Context,
	cacheRepo repository.ICache,
	prefix, key string,
	ttl time.Duration,
	fn func(ctx context.Context) (T, error),
) (T, error) {
	var zero T
	cacheKey := buildCacheKey(prefix, key)
	log := appContext.GetLogger().
		With(zap.String("__system__", "cache-facade"))

	// 1. Try the cache
	raw, err := cacheRepo.Get(ctx, cacheKey)

	switch {
	case err == nil:
		var cached T
		if uerr := sonic.Unmarshal([]byte(raw), &cached); uerr == nil {
			return cached, nil // cache hit
		} else {
			// corrupt / outdated shape: treat as a miss and overwrite below
			log.Warn(fmt.Sprintf("cache: unmarshal %q failed %v", cacheKey, uerr))
		}
	case errors.Is(err, repository.ErrCacheMiss):
		// normal miss
	default:
		log.Warn(fmt.Sprintf("cache: get %q failed %v", cacheKey, err))
	}

	// 2. Miss: run the real function
	val, err := fn(ctx)
	if err != nil {
		return zero, err // never cache error
	}

	// 3. Store the result (best effort)
	if payload, merr := sonic.Marshal(val); merr != nil {
		log.Warn(fmt.Sprintf("cache: marshal %q failed %v", cacheKey, merr))
	} else if serr := cacheRepo.Set(ctx, cacheKey, string(payload), ttl); serr != nil {
		log.Warn(fmt.Sprintf("cache: marshal %q failed %v", cacheKey, merr))
	}

	return val, nil
}

func buildCacheKey(prefix, key string) string {
	prefix = strings.Trim(prefix, ":")

	return fmt.Sprintf("%s:%s", prefix, key)
}
