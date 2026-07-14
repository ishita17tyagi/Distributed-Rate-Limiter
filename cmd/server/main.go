package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-rate-limiter/internal/api"
	"distributed-rate-limiter/internal/config"
	"distributed-rate-limiter/internal/limiter"
	"distributed-rate-limiter/internal/logger"
	"distributed-rate-limiter/internal/redis"
	"distributed-rate-limiter/internal/server"
	"distributed-rate-limiter/internal/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize structured logger
	logger.Init()

	redisClient := redis.New(cfg.RedisAddr)

	// Startup context for Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx); err != nil {
		logger.Log.Error("Redis connection failed", "error", err)
		os.Exit(1)
	}

	logger.Log.Info("✅ Connected to Redis")
	defer redisClient.Close()

	// Initialize the Redis Store using our connected client
	redisStore := storage.NewRedisStore(redisClient)

	// Initialize the Token Bucket with the Redis Store injected
	tokenBucket := limiter.NewMemoryTokenBucket(
		redisStore,
		cfg.DefaultRateLimit,
		cfg.RateLimitWindow,
	)

	rateLimiter := api.NewRateLimitMiddleware(tokenBucket)

	// Pass the middleware to the server
	srv := server.New(
		cfg,
		rateLimiter,
	)

	// Start HTTP server in a goroutine
	go func() {
		logger.Log.Info(
			"server starting",
			"port", cfg.ServerPort,
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Create a channel to receive OS signals
	stop := make(chan os.Signal, 1)

	// Notify the channel when SIGINT or SIGTERM is received
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait until a signal arrives
	<-stop

	logger.Log.Info("Shutdown signal received...")

	// Give existing requests time to finish using a new context
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Gracefully stop the server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Log.Info("Server stopped gracefully.")
}