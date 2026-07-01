package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, key ...string) *redis.IntCmd
}

// FakeCache is a fake implementation of the Cache interface for testing
type FakeCache struct {
	GetFunc func(ctx context.Context, key string) *redis.StringCmd
	SetFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	DelFunc func(ctx context.Context, key ...string) *redis.IntCmd
}

func (f *FakeCache) Get(ctx context.Context, key string) *redis.StringCmd {
	return f.GetFunc(ctx, key)
}

func (f *FakeCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return f.SetFunc(ctx, key, value, expiration)
}

func (f *FakeCache) Del(ctx context.Context, key ...string) *redis.IntCmd {
	return f.DelFunc(ctx, key...)
}
