package limiter

import (
	"sync"
	"time"
)

type Bucket struct {
	mu         sync.Mutex
	Tokens     float64
	LastRefill time.Time
}