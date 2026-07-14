package storage

import (
	"sync"

	"distributed-rate-limiter/internal/limiter"
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
		return nil, nil
	}

	return bucket, nil
}

func (m *MemoryStore) Save(key string, bucket *limiter.Bucket) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.buckets[key] = bucket

	return nil
}
