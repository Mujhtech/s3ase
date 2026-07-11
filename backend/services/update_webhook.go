package services

import (
	"context"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type UpdateWebhookService struct {
	WebhookId     string
	App           *models.App
	AppMemberRepo store.AppMemberRepository
	WebhookRepo   store.WebhookRepository
	User          *models.User
	Body          *dto.CreateWebhookRequestDto
}

func (c *UpdateWebhookService) Run(ctx context.Context) error {
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
	if err := validateWebhook(c.Body); err != nil {
		return err
	}

	if err := c.WebhookRepo.UpdateWebhook(ctx, &models.Webhook{
		ID:          webhook.ID,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		URL:         c.Body.Url,
		CreatedBy:   c.User.ID,
		Metadata: map[string]interface{}{
			"events": c.Body.Events,
		},
	}); err != nil {
		return err
	}

	return nil
}
