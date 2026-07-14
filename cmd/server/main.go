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
	"distributed-rate-limiter/internal/logger"
	"distributed-rate-limiter/internal/redis"
	"distributed-rate-limiter/internal/server"
	"distributed-rate-limiter/internal/storage"
)

func main() {
	
	logger.Init()
	
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger.Log.Info("configuration loaded", "port", cfg.ServerPort, "failure_mode", cfg.RedisFailureMode)

	redisClient := redis.New(cfg.RedisAddr)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx); err != nil {
		logger.Log.Warn("Redis connection failed on startup, fallback modes will apply", "error", err)
	} else {
		logger.Log.Info("✅ Connected to Redis")
	}
	defer redisClient.Close()

	// RedisStore now implements Limiter Interface directly using Lua execution
	redisLimiter := storage.NewRedisStore(redisClient, cfg)
	rateLimiterMiddleware := api.NewRateLimitMiddleware(redisLimiter, cfg)

	logger.Log.Info("rate limiter middleware initialized", "default_limit", cfg.DefaultRateLimit)

	srv := server.New(cfg, rateLimiterMiddleware)

	go func() {
		logger.Log.Info("🚀 Server listening", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Log.Info("Shutdown signal received...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Shutdown failed", "error", err)
	}
	logger.Log.Info("Server stopped gracefully.")
}
