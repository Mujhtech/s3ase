package pubsub

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/internal/redis"
	redV9 "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Redis struct {
	config   Config
	mutex    sync.RWMutex
	client   redV9.UniversalClient
	registry []Consumer
}

func NewRedis(cfg *config.Config, redis *redis.Redis) (*Redis, error) {

	client := redis.Client()

	return &Redis{
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

func (r *Redis) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {

	pubConfig := PublishConfig{
		app:       r.config.App,
		namespace: r.config.Namespace,
	}

	for _, f := range opts {
		f.Apply(&pubConfig)
	}

	topic = formatTopic(pubConfig.app, pubConfig.namespace, topic)

	err := r.client.Publish(ctx, topic, payload).Err()

	if err != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, err)
	}

	return nil
}

func (r *Redis) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	config := SubscribeConfig{
		topics:         make([]string, 0, 8),
		app:            r.config.App,
		namespace:      r.config.Namespace,
		healthInterval: r.config.HealthInterval,
		sendTimeout:    r.config.SendTimeout,
		channelSize:    r.config.ChannelSize,
	}

	for _, f := range opts {
		f.Apply(&config)
	}

	subscriber := &redisSubscriber{
		config:  &config,
		handler: handler,
	}

	config.topics = append(config.topics, topic)

	topics := subscriber.formatTopics(config.topics...)
	subscriber.rdb = r.client.Subscribe(ctx, topics...)

	// start subscriber
	go subscriber.start(ctx)

	// register subscriber
	r.registry = append(r.registry, subscriber)

	return subscriber
}

func (r *Redis) Close(_ context.Context) error {
	for _, subscriber := range r.registry {
		err := subscriber.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

type redisSubscriber struct {
	config  *SubscribeConfig
	rdb     *redV9.PubSub
	handler func([]byte) error
}

func (r *redisSubscriber) start(ctx context.Context) {
	// Go channel which receives messages.
	ch := r.rdb.Channel(
		redV9.WithChannelHealthCheckInterval(r.config.healthInterval),
		redV9.WithChannelSendTimeout(r.config.sendTimeout),
		redV9.WithChannelSize(r.config.channelSize),
	)
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				log.Ctx(ctx).Debug().Msg("redis channel was closed")
				return
			}
			if err := r.handler([]byte(msg.Payload)); err != nil {
				log.Ctx(ctx).Err(err).Msg("received an error from handler function")
			}
		case <-time.After(5 * time.Second):
			log.Ctx(ctx).Debug().Msg("pubsub blocked writing to subscription channel for >5 seconds")
		}
	}
}

func (r *redisSubscriber) Subscribe(ctx context.Context, topics ...string) error {
	err := r.rdb.Subscribe(ctx, r.formatTopics(topics...)...)
	if err != nil {
		return fmt.Errorf("subscribe failed for chanels %v with error: %w",
			strings.Join(topics, ","), err)
	}
	return nil
}

func (r *redisSubscriber) Unsubscribe(ctx context.Context, topics ...string) error {
	err := r.rdb.Unsubscribe(ctx, r.formatTopics(topics...)...)
	if err != nil {
		return fmt.Errorf("unsubscribe failed for chanels %v with error: %w",
			strings.Join(topics, ","), err)
	}
	return nil
}

func (r *redisSubscriber) Close() error {
	err := r.rdb.Close()
	if err != nil {
		return fmt.Errorf("failed while closing subscriber with error: %w", err)
	}
	return nil
}

func (s *redisSubscriber) formatTopics(topics ...string) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		result[i] = formatTopic(s.config.app, s.config.namespace, topic)
	}
	return result
}
