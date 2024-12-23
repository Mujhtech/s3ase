package services

import (
	"context"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type UpdateApiKeyService struct {
	ApiKeyId   string
	App        *models.App
	ApiKeyRepo store.ApiKeyRepository
	User       *models.User
	Body       *dto.CreateApiKeyRequestDto
}

func (c *UpdateApiKeyService) Run(ctx context.Context) error {

	apiKey, err := c.ApiKeyRepo.FindApiKeyByID(ctx, c.ApiKeyId)

	if err != nil {
		return err
	}

	if err := c.ApiKeyRepo.UpdateApiKey(ctx, &models.ApiKey{
		ID:          apiKey.ID,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		Access:      c.Body.Access,
		ExpiredAt:   null.NewInt(c.Body.ExpiredAt, c.Body.ExpiredAt != 0),
		CreatedBy:   c.User.ID,
		Metadata:    map[string]interface{}{},
	}); err != nil {
		return err
	}

	return nil
}
