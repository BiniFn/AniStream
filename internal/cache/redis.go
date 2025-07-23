package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
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

func (c *RedisClient) GetOrFill(
	ctx context.Context,
	key string,
	dest any,
	ttl time.Duration,
	fetch func(context.Context) (any, error),
) (hit bool, err error) {
	if dest != nil && reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return false, fmt.Errorf("dest for key %s must be a pointer", key)
	}

	ok, err := c.Get(ctx, key, dest)
	if err == nil && ok && c.appEnv != "development" {
		c.log.Debug("cache hit", "key", key)
		return true, nil
	} else if err != nil {
		c.log.Warn("cache get failed, fetching", "key", key, "err", err)
	}

	val, err := fetch(ctx)
	if err != nil {
		return false, err
	}

	if setErr := c.Set(ctx, key, val, ttl); setErr != nil {
		c.log.Warn("cache set failed", "key", key, "err", err)
	}

	if dest != nil {
		dv := reflect.ValueOf(dest).Elem()
		vv := reflect.ValueOf(val)
		if vv.Type().AssignableTo(dv.Type()) {
			dv.Set(vv)
		} else if b, err := json.Marshal(val); err == nil {
			_ = json.Unmarshal(b, dest)
		}
	}

	if c.appEnv == "development" {
		c.log.Debug("cache disabled in dev mode", "key", key)
	}

	return false, nil
}
