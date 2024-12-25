package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type CreateWebhookService struct {
	App         *models.App
	WebhookRepo store.WebhookRepository
	User        *models.User
	Body        *dto.CreateWebhookRequestDto
}

func (c *CreateWebhookService) Run(ctx context.Context) (*models.Webhook, error) {

	id := uuid.New().String()

	if err := c.WebhookRepo.CreateWebhook(ctx, &models.Webhook{
		ID:          id,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		URL:         c.Body.Url,
		CreatedBy:   c.User.ID,
		Metadata: map[string]interface{}{
			"events": c.Body.Events,
		},
	}); err != nil {
		return nil, err
	}

	return c.WebhookRepo.FindWebhookByID(ctx, id)
}
