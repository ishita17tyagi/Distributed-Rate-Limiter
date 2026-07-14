package storage

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"distributed-rate-limiter/internal/limiter"
	"distributed-rate-limiter/internal/redis"
	goredis "github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (s *RedisStore) Get(key string) (*limiter.Bucket, error) {
	ctx := context.Background()
	val, err := s.client.RDB.Get(ctx, key).Result()

	if errors.Is(err, goredis.Nil) {
		return nil, nil // Bucket doesn't exist yet
	} else if err != nil {
		return nil, err
	}

	var bucket limiter.Bucket
	if err := json.Unmarshal([]byte(val), &bucket); err != nil {
		return nil, err
	}

	return &bucket, nil
}

func (s *RedisStore) Save(key string, bucket *limiter.Bucket) error {
	ctx := context.Background()

	data, err := json.Marshal(bucket)
	if err != nil {
		return err
	}

	// Save to Redis with a TTL (e.g., 1 hour so unused buckets expire)
	return s.client.RDB.Set(ctx, key, data, 1*time.Hour).Err()
}
