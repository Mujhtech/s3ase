package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type DeleteWebhookService struct {
	WebhookId string
	App       *models.App

	AppMemberRepo store.AppMemberRepository
	WebhookRepo   store.WebhookRepository
	User          *models.User
}

func (c *DeleteWebhookService) Run(ctx context.Context) error {
	if err := requireAppOwner(ctx, c.AppMemberRepo, c.App.ID, c.User.ID); err != nil {
		return err
	}

	webhook, err := c.WebhookRepo.FindWebhookByID(ctx, c.WebhookId)

	if err != nil {
		return err
	}
	if err := requireSameApp(c.App.ID, webhook.AppID); err != nil {
		return err
	}

	if err := c.WebhookRepo.DeleteWebhook(ctx, webhook.ID); err != nil {
		return err
	}

	return nil
}
