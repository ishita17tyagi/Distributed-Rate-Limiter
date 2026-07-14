package limiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, key string, capacity int, window time.Duration) (*Result, error)
}
