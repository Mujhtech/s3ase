package sse

import (
	"context"
	"encoding/json"

	"github.com/mujhtech/s3ase/internal/pkg/pubsub"
)

type pubsubSse struct {
	pubsub pubsub.Pubsub
}

type Streamer interface {
	Publish(ctx context.Context, id string, eventType EventType, data interface{}) error
	Subscribe(ctx context.Context, id string) (<-chan *Event, <-chan error, func(context.Context) error)
}

func NewStreamer(pubsub pubsub.Pubsub) Streamer {
	return &pubsubSse{
		pubsub: pubsub,
	}
}

func (s *pubsubSse) Publish(ctx context.Context, id string, eventType EventType, data interface{}) error {
	serializedData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	event := &Event{
		Type: eventType,
		Data: serializedData,
	}

	payload, err := json.Marshal(event)

	if err != nil {
		return err
	}

	namespaceOption := pubsub.WithPublishNamespace("sse")
	if err = s.pubsub.Publish(ctx, id, payload, namespaceOption); err != nil {
		return err
	}

	return nil
}

func (s *pubsubSse) Subscribe(ctx context.Context, id string) (<-chan *Event, <-chan error, func(context.Context) error) {
	eventChs := make(chan *Event, 100)
	errCh := make(chan error)

	handler := func(payload []byte) error {
		event := &Event{}
		err := json.Unmarshal(payload, event)

		if err != nil {
			return err
		}

		select {
		case eventChs <- event:
		default:
		}

		return nil
	}

	namespaceOption := pubsub.WithChannelNamespace("sse")
	consumer := s.pubsub.Subscribe(ctx, id, handler, namespaceOption)

	cleanupFn := func(ctx context.Context) error {
		return consumer.Close()
	}

	return eventChs, errCh, cleanupFn
}
