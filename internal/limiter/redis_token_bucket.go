package limiter

import "context"

type RedisTokenBucket struct {
	// We'll inject the Redis store in Phase 8B.
}

func NewRedisTokenBucket() *RedisTokenBucket {
	return &RedisTokenBucket{}
}

func (r *RedisTokenBucket) Allow(ctx context.Context, key string) (*Result, error) {
	panic("not implemented")
}
