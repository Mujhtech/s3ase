package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindWebhooksService struct {
	App         *models.App
	WebhookRepo store.WebhookRepository
	User        *models.User
}

func (c *FindWebhooksService) Run(ctx context.Context) ([]*models.Webhook, error) {

	webhooks, err := c.WebhookRepo.FindWebhookByAppID(ctx, c.App.ID)

	if err != nil {
		return nil, err
	}

	return webhooks, nil
}
