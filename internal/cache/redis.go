package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	r Redis
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

func (c *RedisClient) Pipeline() *RedisClient {
	return &RedisClient{r: c.r.Pipeline()}
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

	ok, getErr := c.Get(ctx, key, dest)
	if getErr == nil && ok {
		log.Printf("cache hit for key %s", key)
		return true, nil
	}
	if getErr != nil {
		log.Printf("cache get %s error: %v (fetching anyway)", key, getErr)
	}

	val, err := fetch(ctx)
	if err != nil {
		return false, err
	}

	if setErr := c.Set(ctx, key, val, ttl); setErr != nil {
		log.Printf("cache set %s error: %v", key, setErr)
	}

	if dest != nil {
		dv := reflect.ValueOf(dest).Elem()
		vv := reflect.ValueOf(val)
		if vv.Type().AssignableTo(dv.Type()) {
			dv.Set(vv)
		} else {
			if b, marshalErr := json.Marshal(val); marshalErr == nil {
				if unmarshalErr := json.Unmarshal(b, dest); unmarshalErr != nil {
					log.Printf("cache copy %s unmarshal error: %v", key, unmarshalErr)
				}
			} else {
				log.Printf("cache copy %s marshal error: %v", key, marshalErr)
			}
		}
	}

	return false, nil
}
