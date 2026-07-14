package api

import (
	"net/http"
)

func NewRouter(rateLimiter *RateLimitMiddleware) *http.ServeMux {

	mux := http.NewServeMux()

	mux.Handle(
		"/health",
		rateLimiter.Handler(
			http.HandlerFunc(HealthHandler),
		),
	)

	return mux
}
