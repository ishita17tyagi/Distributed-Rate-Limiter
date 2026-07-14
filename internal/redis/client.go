package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *goredis.Client
}

func New(addr string) *Client {
	rdb := goredis.NewClient(&goredis.Options{
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
