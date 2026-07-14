package limiter

import (
	"context"
	"sync"
	"time"
)

type TokenBucket struct {
	mu sync.Mutex

	buckets map[string]*Bucket

	capacity   float64
	refillRate float64
}

func NewTokenBucket(capacity int, window time.Duration) *TokenBucket {
	refillRate := float64(capacity) / window.Seconds()

	return &TokenBucket{
		buckets: make(map[string]*Bucket),

		capacity:   float64(capacity),
		refillRate: refillRate,
	}
}

func (tb *TokenBucket) getBucket(key string) *Bucket {
	bucket, exists := tb.buckets[key]
	if !exists {
		bucket = &Bucket{
			Tokens:     tb.capacity,
			LastRefill: time.Now(),
		}
		tb.buckets[key] = bucket
	}

	return bucket
}

func (tb *TokenBucket) refill(bucket *Bucket) {
	now := time.Now()

	elapsed := now.Sub(bucket.LastRefill).Seconds()

	bucket.Tokens += elapsed * tb.refillRate

	if bucket.Tokens > tb.capacity {
		bucket.Tokens = tb.capacity
	}

	bucket.LastRefill = now
}

func (tb *TokenBucket) consume(bucket *Bucket) bool {
	if bucket.Tokens >= 1 {
		bucket.Tokens--
		return true
	}

	return false
}

func (tb *TokenBucket) Allow(ctx context.Context, key string) (*Result, error) {
	// ctx will be used later in the Redis implementation.
	_ = ctx

	tb.mu.Lock()
	defer tb.mu.Unlock()

	bucket := tb.getBucket(key)

	tb.refill(bucket)

	allowed := tb.consume(bucket)

	result := &Result{
		Allowed:    allowed,
		Remaining:  int(bucket.Tokens),
		RetryAfter: 0,
	}

	if !allowed {
		// Number of seconds until one token becomes available.
		missingTokens := 1 - bucket.Tokens
		result.RetryAfter = int64(missingTokens / tb.refillRate)
	}

	return result, nil
}
