package pubsub

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/mujhtech/s3ase/config"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type ApacheKafka struct {
	config   Config
	mutex    sync.RWMutex
	dialer   *kafka.Dialer
	brokers  []string
	registry []Consumer
}

func NewApacheKafka(cfg *config.Config, ctx context.Context) (*ApacheKafka, error) {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if err := validateConnection(ctx, dialer, cfg.Pubsub.Brokers); err != nil {
		return nil, fmt.Errorf("failed to connect to kafka: %w", err)
	}

	return &ApacheKafka{
		config: Config{
			App:            cfg.Pubsub.App,
			Namespace:      cfg.Pubsub.Namespace,
			HealthInterval: cfg.Pubsub.HealthInterval,
			SendTimeout:    cfg.Pubsub.SendTimeout,
			ChannelSize:    cfg.Pubsub.ChannelSize,
		},
		dialer:   dialer,
		brokers:  cfg.Pubsub.Brokers,
		registry: make([]Consumer, 0, 16),
	}, nil
}

func validateConnection(ctx context.Context, dialer *kafka.Dialer, brokers []string) error {
	if len(brokers) == 0 {
		return fmt.Errorf("no kafka brokers configured")
	}

	conn, err := dialer.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}

func (a *ApacheKafka) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	pubConfig := PublishConfig{
		app:       a.config.App,
		namespace: a.config.Namespace,
	}

	for _, f := range opts {
		f.Apply(&pubConfig)
	}

	topic = formatTopic(pubConfig.app, pubConfig.namespace, topic)

	writer := &kafka.Writer{
		Addr:         kafka.TCP(a.brokers...),
		Topic:        topic,
		BatchTimeout: time.Second,
		Transport: &kafka.Transport{
			//Dial: a.dialer.DialContext,
			TLS: &tls.Config{},
		},
	}
	defer writer.Close()

	err := writer.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})

	if err != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, err)
	}

	return nil
}

func (a *ApacheKafka) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	config := SubscribeConfig{
		topics:         make([]string, 0, 8),
		app:            a.config.App,
		namespace:      a.config.Namespace,
		healthInterval: a.config.HealthInterval,
		sendTimeout:    a.config.SendTimeout,
		channelSize:    a.config.ChannelSize,
	}

	for _, f := range opts {
		f.Apply(&config)
	}

	subscriber := &kafkaSubscriber{
		config:  &config,
		handler: handler,
		dialer:  a.dialer,
		brokers: a.brokers,
	}

	config.topics = append(config.topics, topic)
	topics := subscriber.formatTopics(config.topics...)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     a.brokers,
		Topic:       topics[0],
		GroupID:     config.app,
		Dialer:      a.dialer,
		StartOffset: kafka.FirstOffset,
	})

	subscriber.reader = reader

	// start subscriber
	go subscriber.start(ctx)

	// register subscriber
	a.registry = append(a.registry, subscriber)

	return subscriber
}

func (a *ApacheKafka) Close(ctx context.Context) error {
	for _, subscriber := range a.registry {
		if err := subscriber.Close(); err != nil {
			return err
		}
	}
	return nil
}

type kafkaSubscriber struct {
	config  *SubscribeConfig
	reader  *kafka.Reader
	handler func([]byte) error
	dialer  *kafka.Dialer
	brokers []string
}

func (k *kafkaSubscriber) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("error reading message")
				continue
			}

			if err := k.handler(msg.Value); err != nil {
				log.Ctx(ctx).Err(err).Msg("received an error from handler function")
			}
		}
	}
}

func (k *kafkaSubscriber) Subscribe(ctx context.Context, topics ...string) error {
	// Kafka reader is single-topic, would need to create new readers for additional topics
	return fmt.Errorf("multiple topic subscription not supported")
}

func (k *kafkaSubscriber) Unsubscribe(ctx context.Context, topics ...string) error {
	// No direct unsubscribe in kafka-go, would need to close and recreate reader
	return nil
}

func (k *kafkaSubscriber) Close() error {
	return k.reader.Close()
}

func (k *kafkaSubscriber) formatTopics(topics ...string) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		result[i] = formatTopic(k.config.app, k.config.namespace, topic)
	}
	return result
}
