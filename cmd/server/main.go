package main

import (
	"log"
	"net/http"
	"time"

	"distributed-rate-limiter/internal/api"
	"distributed-rate-limiter/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", api.HealthCheck)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      api.RequestLogger(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
