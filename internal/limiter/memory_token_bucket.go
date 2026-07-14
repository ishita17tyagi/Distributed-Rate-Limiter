package limiter

import (
	"context"
	"sync"
	"time"
)

// Define the interface here to avoid import cycles!
type Store interface {
	Get(key string) (*Bucket, error)
	Save(key string, bucket *Bucket) error
}

type MemoryTokenBucket struct {
	mu sync.Mutex

	store Store

	capacity   float64
	refillRate float64
}

func NewMemoryTokenBucket(
	store Store, // <-- Fixed: Changed from storage.BucketStore to Store
	capacity int,
	window time.Duration,
) *MemoryTokenBucket {
	refillRate := float64(capacity) / window.Seconds()

	return &MemoryTokenBucket{
		store:      store,
		capacity:   float64(capacity),
		refillRate: refillRate,
	}
}

func (tb *MemoryTokenBucket) getBucket(key string) (*Bucket, error) {
	bucket, err := tb.store.Get(key)
	if err != nil {
		return nil, err
	}

	if bucket == nil {
		bucket = &Bucket{
			Tokens:     tb.capacity,
			LastRefill: time.Now(),
		}

		if err := tb.store.Save(key, bucket); err != nil {
			return nil, err
		}
	}

	return bucket, nil
}

func (tb *MemoryTokenBucket) refill(bucket *Bucket) {
	now := time.Now()

	elapsed := now.Sub(bucket.LastRefill).Seconds()

	bucket.Tokens += elapsed * tb.refillRate

	if bucket.Tokens > tb.capacity {
		bucket.Tokens = tb.capacity
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

func (tb *MemoryTokenBucket) Allow(ctx context.Context, key string) (*Result, error) {
	_ = ctx // Used later

	tb.mu.Lock()
	defer tb.mu.Unlock()

	bucket, err := tb.getBucket(key)
	if err != nil {
		return nil, err
	}

	tb.refill(bucket)

	allowed := tb.consume(bucket)

	if err := tb.store.Save(key, bucket); err != nil {
		return nil, err
	}

	result := &Result{
		Allowed:    allowed,
		Remaining:  int(bucket.Tokens),
		RetryAfter: 0,
	}

	if !allowed {
		missingTokens := 1 - bucket.Tokens
		result.RetryAfter = int64(missingTokens / tb.refillRate)
	}

	return result, nil
}
