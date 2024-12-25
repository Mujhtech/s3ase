package pubsub

import (
	"context"

	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/internal/redis"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error
}

type Pubsub interface {
	Publisher

	Subscribe(ctx context.Context, topic string,
		handler func(payload []byte) error, opts ...SubscribeOption) Consumer
}

type Consumer interface {
	Subscribe(ctx context.Context, topics ...string) error
	Unsubscribe(ctx context.Context, topics ...string) error
	Close() error
}

func NewPubsub(cfg *config.Config, ctx context.Context, redis *redis.Redis) (Pubsub, error) {

	switch cfg.Pubsub.Provider {
	case config.PubsubProviderAwsSqs:
		return NewAwsSqs(cfg, ctx)
	case config.PubsubProviderRedis:
		return NewRedis(cfg, redis)
	case config.PubsubProviderGoogle:
		return NewGooglePubsub(cfg, ctx)
	case config.PubsubProviderKafka:
		return NewApacheKafa(cfg, ctx)
	default:
		return NewInMemory(cfg)
	}
}
