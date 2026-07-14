package limiter

type Result struct {
	Allowed    bool
	Remaining  int
	RetryAfter int64
}
