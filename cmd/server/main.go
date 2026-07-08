package main

import (
	"fmt"
	"log"

	"distributed-rate-limiter/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	fmt.Printf("server_port=%s\n", cfg.ServerPort)
	fmt.Printf("redis_addr=%s\n", cfg.RedisAddr)
	fmt.Printf("default_rate_limit=%d\n", cfg.DefaultRateLimit)
	fmt.Printf("rate_limit_window=%s\n", cfg.RateLimitWindow)
}
