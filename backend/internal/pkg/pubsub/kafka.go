package pubsub

import (
	"context"

	"github.com/mujhtech/s3ase/config"
	"github.com/segmentio/kafka-go"
)

type ApacheKafa struct {
	conn *kafka.Conn
}

func NewApacheKafa(cfg *config.Config, ctx context.Context) (*ApacheKafa, error) {

	conn, err := kafka.Dial("tcp", "localhost:9092")

	if err != nil {
		return nil, err
	}

	return &ApacheKafa{
		conn: conn,
	}, nil
}

func (a *ApacheKafa) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	return nil
}

func (a *ApacheKafa) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	return nil
}
