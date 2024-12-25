package pubsub

import (
	"context"

	"github.com/mujhtech/s3ase/config"
)

type InMemory struct {
}

func NewInMemory(cfg *config.Config) (*InMemory, error) {
	return &InMemory{}, nil
}

func (i *InMemory) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	return nil
}

func (i *InMemory) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	return nil
}
