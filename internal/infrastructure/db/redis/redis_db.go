package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Conn
}

func NewWithConn(url string) (*Redis, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("NewWithConn: failed parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	return &Redis{
		client: client.Conn(),
	}, nil
}

func (r *Redis) CloseConn() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("CloseConn: failed close redis conn: %w", err)
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, key string, out interface{}) (bool, error) {
	receivedData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, fmt.Errorf("Get: failed get data from redis database: %w", err)
	}

	if err := json.Unmarshal([]byte(receivedData), &out); err != nil {
		return false, fmt.Errorf("Get: failed unmarshal cache data to out interface: %w", err)
	}
	return true, nil
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Set: failed marshalling data: %w", err)
	}

	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("Set: failed set new data to redis database: %w", err)
	}

	return nil
}

func (r *Redis) Delete(ctx context.Context, pattern string) (bool, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return false, fmt.Errorf("Delete: failed get keys by pattern: %w", err)
	}

	for _, key := range keys {
		if err := r.client.Del(ctx, key); err != nil {
			return false, fmt.Errorf("Delete: failed delete data from redis: %w", err.Err())
		}
	}
	return true, nil
}
