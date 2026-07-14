package server

import (
	"net/http"
	"time"

	"distributed-rate-limiter/internal/api"
	"distributed-rate-limiter/internal/config"
)

func New(cfg *config.Config) *http.Server {

	return &http.Server{
		Addr: ":" + cfg.ServerPort,

		Handler: api.NewRouter(),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}
