package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/mujhtech/s3ase/config"
)

const (
	charSet = "UTF-8"
)

type Email struct {
	client      *ses.Client
	fromAddress string
}

func New(cfg *config.Config) (*Email, error) {
	config, err := awsConfig.LoadDefaultConfig(
		context.Background(),
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

	client := ses.NewFromConfig(config)

	return &Email{
		client:      client,
		fromAddress: cfg.Email.FromAddress,
	}, nil
}

func (e *Email) send(ctx context.Context, toAddresses []string, subject string, body string, htmlBody string) error {

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: toAddresses,
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(htmlBody),
				},
				Text: &types.Content{
					Charset: aws.String(charSet),
					Data:    aws.String(body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String(charSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(e.fromAddress),
	}

	// Send the email
	_, err := e.client.SendEmail(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
