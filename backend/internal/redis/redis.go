package redis

import (
	"fmt"

	"github.com/mujhtech/s3ase/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cfg    config.Redis
	client redis.UniversalClient
}

func NewRedis(cfg *config.Config) (*Redis, error) {

	var options *redis.Options

	if cfg.Redis.Password != "" {

		var err error

		options, err = redis.ParseURL(cfg.Redis.BuildDsn())

		if err != nil {
			return nil, err
		}

		options.MaxRetries = cfg.Redis.MaxRetries
		options.MinIdleConns = cfg.Redis.MinIdleConnections
		options.DB = cfg.Redis.DB

	} else {
		options = &redis.Options{
			Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			MaxRetries:   cfg.Redis.MaxRetries,
			MinIdleConns: cfg.Redis.MinIdleConnections,
			DB:           cfg.Redis.DB,
		}
	}

	return &Redis{
		client: redis.NewClient(options),
		cfg:    cfg.Redis,
	}, nil
}

func (r *Redis) Client() redis.UniversalClient {
	return r.client
}

func (r *Redis) MakeRedisClient() interface{} {
	return r.client
}
