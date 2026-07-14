package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"distributed-rate-limiter/internal/limiter"
	"distributed-rate-limiter/internal/logger"
	redisclient "distributed-rate-limiter/internal/redis"
)

const (
	keyPrefix = "rate_limit:"
	bucketTTL = 2 * time.Minute
)

type RedisStore struct {
	client *redisclient.Client
}

func NewRedisStore(client *redisclient.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) Get(key string) (*limiter.Bucket, error) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, keyPrefix+key)
	if err != nil {
		return nil, err
	}

	if data == "" {
		return nil, nil
	}

	var bucket limiter.Bucket

	if err := json.Unmarshal([]byte(data), &bucket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bucket: %w", err)
	}

	logger.Log.Debug(
		"bucket loaded",
		"key", key,
	)

	return &bucket, nil
}

func (r *RedisStore) Save(key string, bucket *limiter.Bucket) error {
	ctx := context.Background()

	data, err := json.Marshal(bucket)
	if err != nil {
		return err
	}

	logger.Log.Debug(
		"bucket saved",
		"key", key,
	)

	return r.client.Set(
		ctx,
		keyPrefix+key,
		data,
		bucketTTL,
	)
}