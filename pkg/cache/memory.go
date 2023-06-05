package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/begenov/courses-service/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key ...string) error
}

type MemoryCache struct {
	client *redis.Client
}

func NewMemoryCache(ctx context.Context, cfg config.RedisConfig) (*MemoryCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &MemoryCache{
		client: client,
	}, nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := m.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (m *MemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := m.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (m *MemoryCache) Delete(ctx context.Context, key ...string) error {
	err := m.client.Del(ctx, key...).Err()
	if err != nil {
		return err
	}

	return nil
}
