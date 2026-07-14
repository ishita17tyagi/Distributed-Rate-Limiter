package api

import (
	"net/http"
	"strconv"
	"time"

	"distributed-rate-limiter/internal/config"
	"distributed-rate-limiter/internal/limiter"
	"distributed-rate-limiter/internal/logger"
)

type RateLimitMiddleware struct {
	limiter limiter.Limiter
	cfg     *config.Config
}

func NewRateLimitMiddleware(l limiter.Limiter, cfg *config.Config) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: l, cfg: cfg}
}

func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("evaluating rate limit", "method", r.Method, "path", r.URL.Path)
		key := r.Header.Get("X-API-Key")
		if key == "" {
			key = "global"
		}

		capacity := m.cfg.DefaultRateLimit
		window := m.cfg.RateLimitWindow

		// Now using pre-parsed values directly from memory - ZERO allocations!
		if limit, exists := m.cfg.ClientLimits[key]; exists {
			capacity = limit.Capacity
			window = limit.Window
		}

		start := time.Now()
		result, err := m.limiter.Allow(r.Context(), key, capacity, window)
		latency := time.Since(start)
		logger.Log.Info("rate limit decision",
            "client", key,
            "allowed", result.Allowed,
            "remaining", result.Remaining,
            "latency", latency.String(),
        )

		if err != nil {
			logger.Log.Error("rate limit error", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		logger.Log.Info("rate limit check",
			"client", key,
			"allowed", result.Allowed,
			"remaining", result.Remaining,
			"latency", latency.String(),
		)

		if !result.Allowed {
			w.Header().Set("Retry-After", strconv.FormatInt(result.RetryAfter, 10))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}