package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindDomainService struct {
	App        *models.App
	DomainRepo store.DomainRepository
	User       *models.User
}

func (c *FindDomainService) Run(ctx context.Context) (*models.Domain, error) {

	domain, err := c.DomainRepo.FindDomainByAppID(ctx, c.App.ID)

	if err != nil {
		return nil, err
	}

	return domain, nil
}
