package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

type CreateAppService struct {
	Body                *dto.CreateAppRequestDto
	DefaultRegion       string
	ForceDefaultRegion  bool
	AppRepo             store.AppRepository
	AppMemberRepo       store.AppMemberRepository
	AppSubscriptionRepo store.AppSubscriptionRepository
	User                *models.User
	S3                  AppBucketStore
}

type AppBucketStore interface {
	CheckOrCreateNewBucket(ctx context.Context, bucket string, region string) (string, error)
}

func (c *CreateAppService) Run(ctx context.Context) (*models.App, error) {

	region := c.DefaultRegion

	if !c.ForceDefaultRegion && c.Body.Region != "" {
		region = c.Body.Region
	}

	slug := slugify(c.Body.Name)
	if len(slug) < 3 {
		return nil, fmt.Errorf("%w: app name must contain at least three letters or numbers", errs.ErrInvalidInput)
	}

	if c.S3 == nil {
		return nil, fmt.Errorf("object storage is not configured")
	}
	if _, err := c.S3.CheckOrCreateNewBucket(ctx, slug, region); err != nil {
		return nil, fmt.Errorf("failed to prepare app bucket: %w", err)
	}

	app := &models.App{
		ID:          uuid.New().String(),
		OwnerID:     c.User.ID,
		Name:        c.Body.Name,
		Slug:        slug,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		Metadata:    models.Metadata{},
		Bucket:      slug,
		Region:      null.NewString(region, region != ""),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Persist the app only after its object-storage bucket is ready.
	err := c.AppRepo.CreateApp(ctx, app)

	if err != nil {
		return nil, err
	}

	// create app member
	if err := c.AppMemberRepo.CreateAppMember(ctx, &models.AppMember{
		AppID:    app.ID,
		UserId:   c.User.ID,
		Role:     models.AppMemberRoleOwner,
		Metadata: map[string]interface{}{},
	}); err != nil {
		return nil, err
	}
	if c.AppSubscriptionRepo != nil {
		if err := c.AppSubscriptionRepo.CreateAppSubscription(ctx, &models.AppSubscription{
			ID: uuid.NewString(), AppID: app.ID, Plan: "basic", Status: "active", Metadata: models.Metadata{},
		}); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func slugify(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(slug, "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 63 {
		slug = strings.TrimRight(slug[:63], "-")
	}
	return slug
}
