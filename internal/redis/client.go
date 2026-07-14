package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *redis.Client
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     10,
		MinIdleConns: 2,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
	return &Client{RDB: rdb}
}

func (c *Client) Ping(ctx context.Context) error { return c.RDB.Ping(ctx).Err() }
func (c *Client) Close() error                   { return c.RDB.Close() }