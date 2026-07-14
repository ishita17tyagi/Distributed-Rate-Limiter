package main

import (
	"log"

	"distributed-rate-limiter/internal/config"
	"distributed-rate-limiter/internal/server"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(cfg)

	log.Printf("Server listening on :%s", cfg.ServerPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
