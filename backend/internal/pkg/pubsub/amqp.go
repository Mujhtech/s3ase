package pubsub

import (
	"context"
	"fmt"
	"sync"

	"github.com/mujhtech/s3ase/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQP struct {
	config    Config
	conn      *amqp.Connection
	channel   *amqp.Channel
	mutex     sync.RWMutex
	registry  []Consumer
	closeOnce sync.Once
}

func NewAMQP(cfg *config.Config) (*AMQP, error) {
	conn, err := amqp.Dial(cfg.Pubsub.Amqp.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return &AMQP{
		config: Config{
			App:            cfg.Pubsub.App,
			Namespace:      cfg.Pubsub.Namespace,
			HealthInterval: cfg.Pubsub.HealthInterval,
			SendTimeout:    cfg.Pubsub.SendTimeout,
			ChannelSize:    cfg.Pubsub.ChannelSize,
		},
		conn:     conn,
		channel:  ch,
		registry: make([]Consumer, 0),
	}, nil
}

func (a *AMQP) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	pubConfig := PublishConfig{
		app:       a.config.App,
		namespace: a.config.Namespace,
	}

	for _, f := range opts {
		f.Apply(&pubConfig)
	}

	exchangeName := formatTopic(pubConfig.app, pubConfig.namespace, topic)

	err := a.channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	return a.channel.PublishWithContext(
		ctx,
		exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        payload,
		},
	)
}

func (a *AMQP) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	config := SubscribeConfig{
		topics:      []string{topic},
		app:         a.config.App,
		namespace:   a.config.Namespace,
		channelSize: a.config.ChannelSize,
		sendTimeout: a.config.SendTimeout,
	}

	for _, f := range opts {
		f.Apply(&config)
	}

	subscriber := &amqpSubscriber{
		config:  &config,
		handler: handler,
		amqp:    a,
		queues:  make(map[string]string),
		done:    make(chan struct{}),
	}

	go subscriber.start(ctx)

	a.mutex.Lock()
	a.registry = append(a.registry, subscriber)
	a.mutex.Unlock()

	_ = subscriber.Subscribe(ctx, topic)
	return subscriber
}

func (a *AMQP) Close(ctx context.Context) error {
	var err error
	a.closeOnce.Do(func() {
		a.mutex.Lock()
		defer a.mutex.Unlock()

		for _, subscriber := range a.registry {
			if cerr := subscriber.Close(); cerr != nil {
				err = cerr
			}
		}

		if a.channel != nil {
			if cerr := a.channel.Close(); cerr != nil {
				err = cerr
			}
		}

		if a.conn != nil {
			if cerr := a.conn.Close(); cerr != nil {
				err = cerr
			}
		}
	})
	return err
}

type amqpSubscriber struct {
	config  *SubscribeConfig
	handler func([]byte) error
	amqp    *AMQP
	queues  map[string]string
	done    chan struct{}
	mutex   sync.RWMutex
}

func (s *amqpSubscriber) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		}
	}
}

func (s *amqpSubscriber) Subscribe(ctx context.Context, topics ...string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, topic := range topics {
		exchangeName := formatTopic(s.config.app, s.config.namespace, topic)

		err := s.amqp.channel.ExchangeDeclare(
			exchangeName,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange: %v", err)
		}

		q, err := s.amqp.channel.QueueDeclare(
			"",    // name
			false, // durable
			true,  // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue: %v", err)
		}

		err = s.amqp.channel.QueueBind(
			q.Name,
			"",
			exchangeName,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue: %v", err)
		}

		msgs, err := s.amqp.channel.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to register consumer: %v", err)
		}

		s.queues[topic] = q.Name

		go func() {
			for {
				select {
				case <-s.done:
					return
				case msg, ok := <-msgs:
					if !ok {
						return
					}
					if err := s.handler(msg.Body); err != nil {
						// Log error but continue processing
						continue
					}
				}
			}
		}()
	}
	return nil
}

func (s *amqpSubscriber) Unsubscribe(ctx context.Context, topics ...string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, topic := range topics {
		if queueName, exists := s.queues[topic]; exists {
			_, err := s.amqp.channel.QueueDelete(
				queueName,
				false,
				false,
				false,
			)
			if err != nil {
				return fmt.Errorf("failed to delete queue: %v", err)
			}
			delete(s.queues, topic)
		}
	}
	return nil
}

func (s *amqpSubscriber) Close() error {
	close(s.done)
	return s.Unsubscribe(context.Background(), s.config.topics...)
}
