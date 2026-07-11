package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type DeleteApiKeyService struct {
	ApiKeyId      string
	App           *models.App
	AppMemberRepo store.AppMemberRepository
	ApiKeyRepo    store.ApiKeyRepository
	User          *models.User
}

func (c *DeleteApiKeyService) Run(ctx context.Context) error {
	if err := requireAppOwner(ctx, c.AppMemberRepo, c.App.ID, c.User.ID); err != nil {
		return err
	}

	apiKey, err := c.ApiKeyRepo.FindApiKeyByID(ctx, c.ApiKeyId)

	if err != nil {
		return err
	}
	if err := requireSameApp(c.App.ID, apiKey.AppID); err != nil {
		return err
	}

	if err := c.ApiKeyRepo.DeleteApiKey(ctx, apiKey.ID); err != nil {
		return err
	}

	return nil
}
