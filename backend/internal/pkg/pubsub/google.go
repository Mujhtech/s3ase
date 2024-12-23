package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/mujhtech/s3ase/config"
)

type GooglePubsub struct {
	client *pubsub.Client
}

func NewGooglePubsub(cfg *config.Config, ctx context.Context) (*GooglePubsub, error) {

	client, err := pubsub.NewClient(ctx, cfg.Pubsub.Google.ProjectID)

	if err != nil {
		return nil, err
	}

	return &GooglePubsub{
		client: client,
	}, nil
}

func (g *GooglePubsub) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	return nil
}

func (g *GooglePubsub) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	return nil
}
