package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindAppsService struct {
	AppRepo store.AppRepository
	User    *models.User
}

func (c *FindAppsService) Run(ctx context.Context) ([]*models.App, error) {
	apps, err := c.AppRepo.FindAppsByUserID(ctx, c.User.ID)

	if err != nil {
		return nil, err
	}

	return apps, nil
}
