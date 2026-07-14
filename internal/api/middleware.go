package api

import (
	"net/http"
	"strconv"

	"distributed-rate-limiter/internal/limiter"
)

type RateLimitMiddleware struct {
	limiter limiter.Limiter
}

func NewRateLimitMiddleware(l limiter.Limiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: l,
	}
}

func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// For now we'll identify everyone as the same user.
		// Later we'll use IP or Authentication.
		key := "global"

		result, err := m.limiter.Allow(r.Context(), key)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !result.Allowed {

			w.Header().Set(
				"Retry-After",
				strconv.FormatInt(result.RetryAfter, 10),
			)

			http.Error(
				w,
				"Rate limit exceeded",
				http.StatusTooManyRequests,
			)

			return
		}

		next.ServeHTTP(w, r)
	})
}
