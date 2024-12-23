package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type DeleteWebhookService struct {
	WebhookId string
	App       *models.App

	WebhookRepo store.WebhookRepository
	User        *models.User
}

func (c *DeleteWebhookService) Run(ctx context.Context) error {

	webhook, err := c.WebhookRepo.FindWebhookByID(ctx, c.WebhookId)

	if err != nil {
		return err
	}

	if err := c.WebhookRepo.DeleteWebhook(ctx, webhook.ID); err != nil {
		return err
	}

	return nil
}
