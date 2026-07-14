package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client // Unexported, enforcing encapsulation
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     10,
		MinIdleConns: 2,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	return &Client{
		rdb: rdb,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result, err := c.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Client) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return c.rdb.Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
}
