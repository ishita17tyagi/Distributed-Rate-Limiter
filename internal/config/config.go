package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ClientLimit struct {
	Capacity int
	Window   time.Duration // Changed from string to parsed Duration
}

type Config struct {
	ServerPort       string
	RedisAddr        string
	RedisFailureMode string
	DefaultRateLimit int
	RateLimitWindow  time.Duration
	ClientLimits     map[string]ClientLimit
}

func Load() (*Config, error) {
	_ = loadDotEnv(".env")

	defaultLimit, _ := getEnvInt("DEFAULT_RATE_LIMIT", 10)
	window, _ := getEnvDuration("RATE_LIMIT_WINDOW", time.Minute)

	limits := make(map[string]ClientLimit)
	data, err := os.ReadFile("client_limits.json")
	if err == nil {
		// Temporary struct to decode JSON strings
		type rawLimit struct {
			Capacity int    `json:"capacity"`
			Window   string `json:"window"`
		}
		rawLimits := make(map[string]rawLimit)
		if err := json.Unmarshal(data, &rawLimits); err == nil {
			for key, val := range rawLimits {
				parsedWindow, err := time.ParseDuration(val.Window)
				if err == nil {
					limits[key] = ClientLimit{
						Capacity: val.Capacity,
						Window:   parsedWindow,
					}
				}
			}
		}
	}

	return &Config{
		ServerPort:       getEnv("PORT", "8080"),
		RedisAddr:        getEnv("REDIS_ADDR", "localhost:6379"),
		RedisFailureMode: getEnv("REDIS_FAILURE_MODE", "memory"),
		DefaultRateLimit: defaultLimit,
		RateLimitWindow:  window,
		ClientLimits:     limits,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return value, nil
}

func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback, nil
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return value, nil
}

func loadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
	return nil
}