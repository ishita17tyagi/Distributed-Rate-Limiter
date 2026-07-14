package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-rate-limiter/internal/config"
	"distributed-rate-limiter/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(cfg)

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("🚀 Server listening on :%s", cfg.ServerPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Create a channel to receive OS signals
	stop := make(chan os.Signal, 1)

	// Notify the channel when SIGINT or SIGTERM is received
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait until a signal arrives
	<-stop

	log.Println("Shutdown signal received...")

	// Give existing requests time to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully stop the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped gracefully.")
}