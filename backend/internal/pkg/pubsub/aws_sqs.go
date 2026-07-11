package pubsub

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mujhtech/s3ase/config"
	"github.com/rs/zerolog/log"
)

type AwsSqs struct {
	config   Config
	client   *sqs.Client
	mutex    sync.RWMutex
	registry []Consumer
}

func NewAwsSqs(cfg *config.Config, ctx context.Context) (*AwsSqs, error) {
	config, err := awsConfig.LoadDefaultConfig(
		ctx,
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.Aws.AccessKey,
				cfg.Aws.SecretKey,
				"",
			),
		),
		awsConfig.WithRegion(cfg.Aws.DefaultRegion),
	)

	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(config)

	return &AwsSqs{
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

func (a *AwsSqs) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	pubConfig := PublishConfig{
		app:       a.config.App,
		namespace: a.config.Namespace,
	}

	for _, f := range opts {
		f.Apply(&pubConfig)
	}

	queueName := formatTopic(pubConfig.app, pubConfig.namespace, topic)

	// Get queue URL
	urlResult, err := a.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		return fmt.Errorf("failed to get queue URL for %s: %w", queueName, err)
	}

	// Send message
	_, err = a.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    urlResult.QueueUrl,
		MessageBody: aws.String(string(payload)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish to queue %s: %w", queueName, err)
	}

	return nil
}

func (a *AwsSqs) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
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

	subscriber := &sqsSubscriber{
		config:  &config,
		handler: handler,
		client:  a.client,
	}

	config.topics = append(config.topics, topic)
	subscriber.queueURLs = make(map[string]string)

	// Get queue URLs for all topics
	for _, t := range subscriber.formatTopics(config.topics...) {
		urlResult, err := a.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: &t,
		})
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msgf("failed to get queue URL for %s", t)
			continue
		}
		subscriber.queueURLs[t] = *urlResult.QueueUrl
	}

	// Start subscriber
	go subscriber.start(ctx)

	// Register subscriber
	a.registry = append(a.registry, subscriber)

	return subscriber
}

func (a *AwsSqs) Close(ctx context.Context) error {
	for _, subscriber := range a.registry {
		if err := subscriber.Close(); err != nil {
			return err
		}
	}
	return nil
}

type sqsSubscriber struct {
	config    *SubscribeConfig
	handler   func([]byte) error
	client    *sqs.Client
	queueURLs map[string]string
	done      chan struct{}
}

func (s *sqsSubscriber) start(ctx context.Context) {
	s.done = make(chan struct{})

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.done:
			return
		default:
			for _, url := range s.queueURLs {
				result, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            &url,
					MaxNumberOfMessages: 10,
					WaitTimeSeconds:     20,
				})
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("error receiving messages")
					continue
				}

				for _, msg := range result.Messages {
					if err := s.handler([]byte(*msg.Body)); err != nil {
						log.Ctx(ctx).Error().Err(err).Msg("received an error from handler function")
						continue
					}

					// Delete message after successful processing
					_, err = s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
						QueueUrl:      &url,
						ReceiptHandle: msg.ReceiptHandle,
					})
					if err != nil {
						log.Ctx(ctx).Error().Err(err).Msg("error deleting message")
					}
				}
			}
		}
	}
}

func (s *sqsSubscriber) Subscribe(ctx context.Context, topics ...string) error {
	for _, topic := range s.formatTopics(topics...) {
		if _, exists := s.queueURLs[topic]; exists {
			continue
		}

		urlResult, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
			QueueName: &topic,
		})
		if err != nil {
			return fmt.Errorf("failed to get queue URL for %s: %w", topic, err)
		}
		s.queueURLs[topic] = *urlResult.QueueUrl
	}
	return nil
}

func (s *sqsSubscriber) Unsubscribe(ctx context.Context, topics ...string) error {
	for _, topic := range s.formatTopics(topics...) {
		delete(s.queueURLs, topic)
	}
	return nil
}

func (s *sqsSubscriber) Close() error {
	if s.done != nil {
		close(s.done)
	}
	return nil
}

func (s *sqsSubscriber) formatTopics(topics ...string) []string {
	result := make([]string, len(topics))
	for i, topic := range topics {
		result[i] = formatTopic(s.config.app, s.config.namespace, topic)
	}
	return result
}
