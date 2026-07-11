package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
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
	name := strings.TrimSpace(c.Body.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: API key name is required", errors.ErrInvalidInput)
	}
	if c.Body.Access != models.ApiKeyAccessRead && c.Body.Access != models.ApiKeyAccessWrite && c.Body.Access != models.ApiKeyAccessFull {
		return nil, fmt.Errorf("%w: invalid API key access", errors.ErrInvalidInput)
	}
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}
	secret := "s3ase_" + base64.RawURLEncoding.EncodeToString(random)
	digest := sha256.Sum256([]byte(secret))

	apiKey := &models.ApiKey{
		ID:          id,
		AppID:       c.App.ID,
		Name:        name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		Access:      c.Body.Access,
		ExpiredAt:   null.NewInt(c.Body.ExpiredAt, c.Body.ExpiredAt != 0),
		CreatedBy:   c.User.ID,
		Metadata:    map[string]interface{}{},
		KeyHash:     hex.EncodeToString(digest[:]),
		KeyPrefix:   secret[:12],
		LastFour:    secret[len(secret)-4:],
		Secret:      secret,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := c.ApiKeyRepo.CreateApiKey(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}
