package storage

import (
	"sync"

	"distributed-rate-limiter/internal/limiter"
	"distributed-rate-limiter/internal/logger"
)

type MemoryStore struct {
	mu      sync.RWMutex
	buckets map[string]*limiter.Bucket
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		buckets: make(map[string]*limiter.Bucket),
	}
}

func (m *MemoryStore) Get(key string) (*limiter.Bucket, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    bucket, exists := m.buckets[key]
    if !exists {
        logger.Log.Debug("cache miss: bucket created", "key", key)
        return nil, nil
    }
    logger.Log.Debug("cache hit: bucket retrieved", "key", key)
    return bucket, nil
}

func (m *MemoryStore) Save(key string, bucket *limiter.Bucket) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.buckets[key] = bucket

	return nil
}
