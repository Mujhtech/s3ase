package services

import (
	"context"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type UpdateWebhookService struct {
	WebhookId   string
	App         *models.App
	WebhookRepo store.WebhookRepository
	User        *models.User
	Body        *dto.CreateWebhookRequestDto
}

func (c *UpdateWebhookService) Run(ctx context.Context) error {

	webhook, err := c.WebhookRepo.FindWebhookByID(ctx, c.WebhookId)

	if err != nil {
		return err
	}

	if err := c.WebhookRepo.CreateWebhook(ctx, &models.Webhook{
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
