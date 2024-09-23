package cache

import (
	"context"
	"fmt"
	"github.com/dostonshernazarov/mini-twitter/internal/pkg/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(cfg *config.Config) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		DB:   0,
	})

	duration, err := time.ParseDuration(cfg.CtxTimeout)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{
		client: client,
	}, nil
}

func (r *RedisStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisStorage) Get(ctx context.Context, key string) (interface{}, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisStorage) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
