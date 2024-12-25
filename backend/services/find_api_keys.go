package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindApiKeysService struct {
	App        *models.App
	AppRepo    store.AppRepository
	ApiKeyRepo store.ApiKeyRepository
	User       *models.User
}

func (c *FindApiKeysService) Run(ctx context.Context) ([]*models.ApiKey, error) {

	apiKeys, err := c.ApiKeyRepo.FindApiKeysByAppID(ctx, c.App.ID)

	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}
