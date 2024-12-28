package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/internal/redis"
)

type RedisCache struct {
	cache *cache.Cache
}

func NewRedisCache(cfg *config.Config, redis *redis.Redis) (*RedisCache, error) {
	client := redis.Client()

	return &RedisCache{
		cache: cache.New(&cache.Options{
			Redis: client,
		}),
	}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string, value interface{}) error {

	err := r.cache.Get(ctx, key, &value)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}
