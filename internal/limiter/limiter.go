package limiter

import "context"

type Limiter interface {
    Allow(ctx context.Context, key string) (*Result, error)
}
