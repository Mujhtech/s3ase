package pubsub

import (
	"context"
	"sync"

	"github.com/mujhtech/s3ase/config"
)

type InMemory struct {
	config   Config
	mutex    sync.RWMutex
	registry []Consumer
	topics   map[string][]chan []byte
}

func NewInMemory(cfg *config.Config) (*InMemory, error) {
	return &InMemory{
		config: Config{
			App:            cfg.Pubsub.App,
			Namespace:      cfg.Pubsub.Namespace,
			HealthInterval: cfg.Pubsub.HealthInterval,
			SendTimeout:    cfg.Pubsub.SendTimeout,
			ChannelSize:    cfg.Pubsub.ChannelSize,
		},
		registry: make([]Consumer, 0, 16),
		topics:   make(map[string][]chan []byte),
	}, nil
}

func (i *InMemory) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	pubConfig := PublishConfig{
		app:       i.config.App,
		namespace: i.config.Namespace,
	}

	for _, f := range opts {
		f.Apply(&pubConfig)
	}

	topicName := formatTopic(pubConfig.app, pubConfig.namespace, topic)

	i.mutex.RLock()
	channels := i.topics[topicName]
	i.mutex.RUnlock()

	for _, ch := range channels {
		select {
		case ch <- payload:
		default:
			// Channel is full, skip this subscriber
		}
	}

	return nil
}

func (i *InMemory) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	config := SubscribeConfig{
		topics:         make([]string, 0, 8),
		app:            i.config.App,
		namespace:      i.config.Namespace,
		healthInterval: i.config.HealthInterval,
		sendTimeout:    i.config.SendTimeout,
		channelSize:    i.config.ChannelSize,
	}

	for _, f := range opts {
		f.Apply(&config)
	}

	subscriber := &inMemorySubscriber{
		config:   &config,
		handler:  handler,
		inmem:    i,
		channels: make(map[string]chan []byte),
		done:     make(chan struct{}),
	}

	config.topics = append(config.topics, topic)

	// Subscribe to initial topic
	topicName := formatTopic(config.app, config.namespace, topic)
	ch := make(chan []byte, config.channelSize)
	subscriber.channels[topicName] = ch

	if i.topics[topicName] == nil {
		i.topics[topicName] = make([]chan []byte, 0)
	}
	i.topics[topicName] = append(i.topics[topicName], ch)

	// Start message handling
	go subscriber.start(ctx)

	i.registry = append(i.registry, subscriber)
	return subscriber
}

func (i *InMemory) Close(ctx context.Context) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	for _, subscriber := range i.registry {
		if err := subscriber.Close(); err != nil {
			return err
		}
	}
	i.topics = make(map[string][]chan []byte)
	return nil
}

type inMemorySubscriber struct {
	config   *SubscribeConfig
	handler  func([]byte) error
	inmem    *InMemory
	channels map[string]chan []byte
	done     chan struct{}
}

func (s *inMemorySubscriber) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		default:
			for _, ch := range s.channels {
				select {
				case msg := <-ch:
					if err := s.handler(msg); err != nil {
						// Log error but continue processing
						continue
					}
				default:
					// No message available
				}
			}
		}
	}
}

func (s *inMemorySubscriber) Subscribe(ctx context.Context, topics ...string) error {
	s.inmem.mutex.Lock()
	defer s.inmem.mutex.Unlock()

	for _, topic := range topics {
		topicName := formatTopic(s.config.app, s.config.namespace, topic)
		if _, exists := s.channels[topicName]; exists {
			continue
		}

		ch := make(chan []byte, s.config.channelSize)
		s.channels[topicName] = ch

		if s.inmem.topics[topicName] == nil {
			s.inmem.topics[topicName] = make([]chan []byte, 0)
		}
		s.inmem.topics[topicName] = append(s.inmem.topics[topicName], ch)
	}
	return nil
}

func (s *inMemorySubscriber) Unsubscribe(ctx context.Context, topics ...string) error {
	s.inmem.mutex.Lock()
	defer s.inmem.mutex.Unlock()

	for _, topic := range topics {
		topicName := formatTopic(s.config.app, s.config.namespace, topic)
		if ch, exists := s.channels[topicName]; exists {
			// Remove channel from topics
			channels := s.inmem.topics[topicName]
			for i, c := range channels {
				if c == ch {
					s.inmem.topics[topicName] = append(channels[:i], channels[i+1:]...)
					break
				}
			}
			delete(s.channels, topicName)
		}
	}
	return nil
}

func (s *inMemorySubscriber) Close() error {
	close(s.done)
	return nil
}
