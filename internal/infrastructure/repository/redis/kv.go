package cache

import (
	"context"
	"time"
)

type KV interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
}

var inst KV

func Init(store KV) {
	inst = store
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return inst.Set(ctx, key, value, expiration)
}

func Get(ctx context.Context, key string) (interface{}, error) {
	return inst.Get(ctx, key)
}

func Del(ctx context.Context, key string) error {
	return inst.Del(ctx, key)
}
