package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type CreateOrUpdateDomainService struct {
	App        *models.App
	DomainRepo store.DomainRepository
	User       *models.User
	Body       *dto.CreateOrUpdateDomainRequestDto
}

func (c *CreateOrUpdateDomainService) Run(ctx context.Context) (*models.Domain, error) {

	domain, err := c.DomainRepo.FindDomainByAppID(ctx, c.App.ID)

	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}

	userDomain := c.Body.Domain

	// remove http:// or https://
	if len(userDomain) > 7 && userDomain[:7] == "http://" {
		userDomain = userDomain[7:]
	} else if len(userDomain) > 8 && userDomain[:8] == "https://" {
		userDomain = userDomain[8:]
	}

	if domain == nil {
		domain := &models.Domain{
			ID:          uuid.New().String(),
			AppID:       c.App.ID,
			Domain:      userDomain,
			CnameRecord: "127.0.0.1",
			TxtRecord:   fmt.Sprintf("_acme-challenge.%s", userDomain),
			Status:      models.DomainStatusPending,
			CreatedBy:   c.User.ID,
			Metadata:    map[string]interface{}{},
		}

		if err := c.DomainRepo.CreateDomain(ctx, domain); err != nil {
			return nil, err
		}
	} else {

		domain.Domain = userDomain
		domain.CnameRecord = "127.0.0.1"
		domain.Status = models.DomainStatusPending
		domain.TxtRecord = fmt.Sprintf("_acme-challenge.%s", userDomain)

		if err := c.DomainRepo.UpdateDomain(ctx, domain); err != nil {
			return nil, err
		}

	}

	return domain, nil
}
