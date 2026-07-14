package limiter

import (
	"context"
	"time"

	"distributed-rate-limiter/internal/logger"
)

type Store interface {
	Get(key string) (*Bucket, error)
	Save(key string, bucket *Bucket) error
}

type MemoryTokenBucket struct {
	// Global mutex REMOVED. We now rely on individual Bucket mutexes.
	store Store
}

func NewMemoryTokenBucket(store Store) *MemoryTokenBucket {
	return &MemoryTokenBucket{
		store: store,
	}
}

func (tb *MemoryTokenBucket) getBucket(key string, capacity float64) (*Bucket, error) {
	bucket, err := tb.store.Get(key)
	if err != nil {
		return nil, err
	}

	if bucket == nil {
		bucket = &Bucket{
			Tokens:     capacity,
			LastRefill: time.Now(),
		}
		if err := tb.store.Save(key, bucket); err != nil {
			return nil, err
		}
	}
	return bucket, nil
}

func (tb *MemoryTokenBucket) refill(bucket *Bucket, capacity float64, refillRate float64) {
	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()
	
	if elapsed < 0 {
		elapsed = 0 // Clock skew safety
	}

	bucket.Tokens += elapsed * refillRate

	if bucket.Tokens > capacity {
		bucket.Tokens = capacity
	}
	bucket.LastRefill = now
}

func (tb *MemoryTokenBucket) consume(bucket *Bucket) bool {
	if bucket.Tokens >= 1 {
		bucket.Tokens--
		return true
	}
	return false
}

func (tb *MemoryTokenBucket) Allow(ctx context.Context, key string, capacity int, window time.Duration) (*Result, error) {
	capFloat := float64(capacity)
	refillRate := capFloat / window.Seconds()

	bucket, err := tb.getBucket(key, capFloat)
	if err != nil {
		return nil, err
	}

	// Lock ONLY this specific user's bucket
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	tb.refill(bucket, capFloat, refillRate)
	allowed := tb.consume(bucket)

	logger.Log.Debug( // Changed to Debug to prevent log spam during failover
		"memory rate limit check (fallback)",
		"user", key,
		"allowed", allowed,
		"remaining", int(bucket.Tokens),
	)

	result := &Result{
		Allowed:    allowed,
		Remaining:  int(bucket.Tokens),
		RetryAfter: 0,
	}

	if !allowed {
		missingTokens := 1 - bucket.Tokens
		result.RetryAfter = int64(missingTokens / refillRate)
	}

	return result, nil
}