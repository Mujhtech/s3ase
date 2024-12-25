package pubsub

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/mujhtech/s3ase/config"
)

type AwsSqs struct {
	client *sqs.Client
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
		client: client,
	}, nil
}

func (a *AwsSqs) Publish(ctx context.Context, topic string, payload []byte, opts ...PublishOption) error {
	return nil
}

func (a *AwsSqs) Subscribe(ctx context.Context, topic string, handler func(payload []byte) error, opts ...SubscribeOption) Consumer {
	return nil
}
