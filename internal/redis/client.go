package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	RDB *goredis.Client // <-- Exported so other packages can use it
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
		RDB: rdb,
	}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.RDB.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.RDB.Close()
}
