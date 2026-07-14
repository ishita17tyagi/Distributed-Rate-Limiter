package limiter

import "time"

type Bucket struct {
	Tokens     float64
	LastRefill time.Time
}
