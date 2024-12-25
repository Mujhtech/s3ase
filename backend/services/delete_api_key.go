package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type DeleteApiKeyService struct {
	ApiKeyId   string
	App        *models.App
	ApiKeyRepo store.ApiKeyRepository
	User       *models.User
}

func (c *DeleteApiKeyService) Run(ctx context.Context) error {

	apiKey, err := c.ApiKeyRepo.FindApiKeyByID(ctx, c.ApiKeyId)

	if err != nil {
		return err
	}

	if err := c.ApiKeyRepo.DeleteApiKey(ctx, apiKey.ID); err != nil {
		return err
	}

	return nil
}
