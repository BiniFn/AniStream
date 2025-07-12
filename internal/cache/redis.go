package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	r *redis.Client
}

func NewRedisClient(ctx context.Context, addr, pass string) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Printf("✅ Connected to Redis at %s", addr)

	return &RedisClient{r: rdb}, nil
}

func (c *RedisClient) Close() error {
	if c.r == nil {
		return nil
	}
	return c.r.Close()
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, expiresIn time.Duration) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.r.Set(ctx, key, valueBytes, expiresIn).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string, dest any) (bool, error) {
	raw, err := c.r.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// cache miss
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal([]byte(raw), dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.r.Del(ctx, key).Err()
}
