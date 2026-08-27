package pkg

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// DeleteByPattern removes all keys matching the given pattern (e.g. "role:list:page:*").
// It uses SCAN to avoid blocking Redis and deletes keys in batches.
func DeleteByPattern(ctx context.Context, rdb *redis.Client, pattern string) error {
	var (
		cursor uint64
		batch  []string
	)

	for {
		keys, next, err := rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		batch = append(batch, keys...)

		// Flush the batch periodically to keep memory bounded.
		if len(batch) >= 500 {
			if err := rdb.Del(ctx, batch...).Err(); err != nil {
				return err
			}
			batch = batch[:0]
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	if len(batch) > 0 {
		if err := rdb.Del(ctx, batch...).Err(); err != nil {
			return err
		}
	}

	return nil
}
