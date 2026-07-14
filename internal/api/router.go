package api

import (
	"net/http"
)

func NewRouter(rateLimiter *RateLimitMiddleware) *http.ServeMux {
	mux := http.NewServeMux()

	handler := LoggingMiddleware(
		rateLimiter.Handler(
			http.HandlerFunc(HealthHandler),
		),
	)

	mux.Handle("/health", handler)

	return mux
}
