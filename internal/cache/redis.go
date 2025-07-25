package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, key ...string) *redis.IntCmd
	Pipeline() redis.Pipeliner
}

type RedisClient struct {
	r      Redis
	appEnv string
	log    *slog.Logger
}

func NewRedisClient(
	ctx context.Context,
	appEnv, addr, pass string,
	log *slog.Logger,
) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	log.Info("connected to redis", "addr", addr)

	return &RedisClient{r: rdb, appEnv: appEnv, log: log}, nil
}

func (c *RedisClient) Close() error {
	if client, ok := c.r.(*redis.Client); ok {
		return client.Close()
	}
	return fmt.Errorf("unsupported Redis client type")
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, expiresIn time.Duration) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.r.Set(ctx, key, valueBytes, expiresIn).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string, dest any) (bool, error) {
	raw, err := c.r.Get(ctx, key).Bytes()
	switch err {
	case nil:
	case redis.Nil:
		return false, nil // cache miss
	default:
		return false, err
	}

	return true, json.Unmarshal(raw, dest)
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.r.Del(ctx, key).Err()
}

func (c *RedisClient) Pipeline() *RedisClient {
	return &RedisClient{r: c.r.Pipeline(), appEnv: c.appEnv, log: c.log}
}

func (c *RedisClient) Exec(ctx context.Context) ([]redis.Cmder, error) {
	if pipe, ok := c.r.(redis.Pipeliner); ok {
		return pipe.Exec(ctx)
	}
	return nil, fmt.Errorf("unsupported Redis client type for pipeline execution")
}

func GetOrFill[T any](
	ctx context.Context,
	rc *RedisClient,
	key string,
	ttl time.Duration,
	fetch func(context.Context) (T, error),
) (val T, err error) {
	var zero T

	var tmp T
	if ok, err := rc.Get(ctx, key, &tmp); err == nil && ok && rc.appEnv != "development" {
		rc.log.Debug("cache hit", "key", key)
		return tmp, nil
	} else if err != nil {
		rc.log.Warn("cache get failed, fetching", "key", key, "err", err)
	}

	tmp, err = fetch(ctx)
	if err != nil {
		return zero, err
	}

	if err = rc.Set(ctx, key, tmp, ttl); err != nil {
		rc.log.Warn("cache set failed", "key", key, "err", err)
	}

	return tmp, nil
}
