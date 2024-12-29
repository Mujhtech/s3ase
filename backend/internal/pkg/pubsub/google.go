package pubsub

import (
	"context"
	"fmt"
	"sync"

	gpubsub "cloud.google.com/go/pubsub"
	"github.com/mujhtech/s3ase/config"
	"github.com/rs/zerolog/log"
)

type GooglePubsub struct {
	config   Config
	client   *gpubsub.Client
	mutex    sync.RWMutex
	registry []Consumer
}

func NewGooglePubsub(cfg *config.Config, ctx context.Context) (*GooglePubsub, error) {
	client, err := gpubsub.NewClient(ctx, cfg.Pubsub.Google.ProjectID)
	if err != nil {
		return nil, err
	}

	return &GooglePubsub{
		config: Config{
			App:            cfg.Pubsub.App,
			Namespace:      cfg.Pubsub.Namespace,
			HealthInterval: cfg.Pubsub.HealthInterval,
			SendTimeout:    cfg.Pubsub.SendTimeout,
			ChannelSize:    cfg.Pubsub.ChannelSize,
		},
		client:   client,
		registry: make([]Consumer, 0, 16),
	}, nil
}

func (g *GooglePubsub) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	pubConfig := PublishConfig{
		app:       g.config.App,
		namespace: g.config.Namespace,
	}

	for _, f := range opts {
		f.Apply(&pubConfig)
	}

	topicName := formatTopic(pubConfig.app, pubConfig.namespace, topic)
	t := g.client.Topic(topicName)
	defer t.Stop()

	result := t.Publish(ctx, &gpubsub.Message{
		Data: payload,
	})

	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topicName, err)
	}

	return nil
}

func (g *GooglePubsub) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	config := SubscribeConfig{
		topics:         make([]string, 0, 8),
		app:            g.config.App,
		namespace:      g.config.Namespace,
		healthInterval: g.config.HealthInterval,
		sendTimeout:    g.config.SendTimeout,
		channelSize:    g.config.ChannelSize,
	}

	for _, f := range opts {
		f.Apply(&config)
	}

	subscriber := &googleSubscriber{
		config:       &config,
		handler:      handler,
		client:       g.client,
		subsriptions: make(map[string]*gpubsub.Subscription),
	}

	config.topics = append(config.topics, topic)

	// Create or get subscription for topic
	for _, t := range subscriber.formatTopics(config.topics...) {
		topic := g.client.Topic(t)
		subName := fmt.Sprintf("%s-sub", t)
		sub := g.client.Subscription(subName)
		exists, err := sub.Exists(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to check subscription existence for %s", t)
			continue
		}

		if !exists {
			sub, err = g.client.CreateSubscription(ctx, subName, gpubsub.SubscriptionConfig{
				Topic: topic,
			})
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf("failed to create subscription for %s", t)
				continue
			}
		}
		subscriber.subsriptions[t] = sub
	}

	// Start subscriber
	go subscriber.start(ctx)

	// Register subscriber
	g.registry = append(g.registry, subscriber)

	return subscriber
}

func (g *GooglePubsub) Close(ctx context.Context) error {
	for _, subscriber := range g.registry {
		if err := subscriber.Close(); err != nil {
			return err
		}
	}
	return g.client.Close()
}

type googleSubscriber struct {
	config       *SubscribeConfig
	handler      func([]byte) error
	client       *gpubsub.Client
	subsriptions map[string]*gpubsub.Subscription
	done         chan struct{}
}

func (s *googleSubscriber) start(ctx context.Context) {
	s.done = make(chan struct{})

	for _, sub := range s.subsriptions {
		go func(subscription *gpubsub.Subscription) {
			err := subscription.Receive(ctx, func(msgCtx context.Context, msg *gpubsub.Message) {
				if err := s.handler(msg.Data); err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("received an error from handler function")
					msg.Nack()
					return
				}
				msg.Ack()
			})
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("error receiving messages")
			}
		}(sub)
	}

	<-s.done
}

func (s *googleSubscriber) Subscribe(ctx context.Context, topics ...string) error {
	for _, topic := range s.formatTopics(topics...) {
		if _, exists := s.subsriptions[topic]; exists {
			continue
		}

		t := s.client.Topic(topic)
		subName := fmt.Sprintf("%s-sub", topic)
		sub := s.client.Subscription(subName)
		exists, err := sub.Exists(ctx)
		if err != nil {
			return fmt.Errorf("failed to check subscription existence for %s: %w", topic, err)
		}

		if !exists {
			sub, err = s.client.CreateSubscription(ctx, subName, gpubsub.SubscriptionConfig{
				Topic: t,
			})
			if err != nil {
				return fmt.Errorf("failed to create subscription for %s: %w", topic, err)
			}
		}
		s.subsriptions[topic] = sub
	}
	return nil
}

func (s *googleSubscriber) Unsubscribe(ctx context.Context, topics ...string) error {
	for _, topic := range s.formatTopics(topics...) {
		if sub, exists := s.subsriptions[topic]; exists {
			sub.Delete(ctx)
			delete(s.subsriptions, topic)
		}
	}
	return nil
}

func (s *googleSubscriber) Close() error {
	if s.done != nil {
		close(s.done)
	}
	return nil
}

func (s *googleSubscriber) formatTopics(topics ...string) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		result[i] = formatTopic(s.config.app, s.config.namespace, topic)
	}
	return result
}
