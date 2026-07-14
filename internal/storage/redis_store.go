package storage

import (
	"context"
	"fmt"
	"time"

	"distributed-rate-limiter/internal/config"
	"distributed-rate-limiter/internal/limiter"
	"distributed-rate-limiter/internal/logger"
	redisclient "distributed-rate-limiter/internal/redis"
	"distributed-rate-limiter/internal/script"

	"github.com/redis/go-redis/v9"
)

var luaScript = script.TokenBucketLua

type RedisStore struct {
	client      *redisclient.Client
	script      *redis.Script
	fallback    *limiter.MemoryTokenBucket
	failureMode string
}

func NewRedisStore(client *redisclient.Client, cfg *config.Config) *RedisStore {
	// go-redis handles SHA caching and NOSCRIPT fallbacks natively!
	script := redis.NewScript(luaScript)

	return &RedisStore{
		client:      client,
		script:      script,
		fallback:    limiter.NewMemoryTokenBucket(NewMemoryStore()),
		failureMode: cfg.RedisFailureMode,
	}
}

func (r *RedisStore) Allow(ctx context.Context, key string, capacity int, window time.Duration) (*limiter.Result, error) {
	refillRate := float64(capacity) / window.Seconds()
	nowMs := time.Now().UnixMilli()

	// Use Run.Result() to get the interface{} return from Lua
	res, err := r.script.Run(ctx, r.client.RDB, []string{"rate_limit:" + key}, capacity, refillRate, nowMs).Result()

	if err != nil {
		// Log the actual error from Redis
		logger.Log.Error("redis operation failed", "error", err, "mode", r.failureMode, "key", key)

		if r.failureMode == "memory" {
			logger.Log.Warn("falling back to memory store")
			return r.fallback.Allow(ctx, key, capacity, window)
		}
		return &limiter.Result{Allowed: false}, err
	}

	// Safely assert the slice
	vals, ok := res.([]interface{})
	if !ok || len(vals) < 3 {
		return nil, fmt.Errorf("invalid lua response format")
	}

	allowedInt := vals[0].(int64)
	remaining := vals[1].(int64)
	retryAfter := vals[2].(int64)

	return &limiter.Result{
		Allowed:    allowedInt == 1,
		Remaining:  int(remaining),
		RetryAfter: retryAfter,
	}, nil
}