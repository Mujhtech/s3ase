package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/errors"
)

type CreateApiKeyService struct {
	App           *models.App
	AppMemberRepo store.AppMemberRepository
	ApiKeyRepo    store.ApiKeyRepository
	User          *models.User
	Body          *dto.CreateApiKeyRequestDto
}

func (c *CreateApiKeyService) Run(ctx context.Context) (*models.ApiKey, error) {

	appMember, err := c.AppMemberRepo.FindAppMemberByAppIDAndUserId(ctx, c.App.ID, c.User.ID)

	if err != nil {
		return nil, err
	}

	if appMember.Role != models.AppMemberRoleOwner {
		return nil, errors.ErrNotAuthorized
	}

	id := uuid.New().String()

	apiKey := &models.ApiKey{
		ID:          id,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		Access:      c.Body.Access,
		ExpiredAt:   null.NewInt(c.Body.ExpiredAt, c.Body.ExpiredAt != 0),
		CreatedBy:   c.User.ID,
		Metadata:    map[string]interface{}{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.ApiKeyRepo.CreateApiKey(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}
