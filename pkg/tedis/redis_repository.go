package tedis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redis_repo struct {
	cli *redis.Client
}

// https://redis.uptrace.dev/guide/get-all-keys.html
func (r *redis_repo) GetKeys(ctx context.Context, cursor uint64, match string, count int64) ([]string, error) {
	var keys []string

	iter := r.cli.Scan(ctx, cursor, match, count).Iterator()

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return keys, err
	}

	return keys, nil
}
