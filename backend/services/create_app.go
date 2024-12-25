package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type CreateAppService struct {
	Body          *dto.CreateAppRequestDto
	DefaultRegion string
	AppRepo       store.AppRepository
	AppMemberRepo store.AppMemberRepository
	ApiKeyRepo    store.ApiKeyRepository
	User          *models.User
	// S3            *s3store.S3Store
	// Job           *job.Job
}

func (c *CreateAppService) Run(ctx context.Context) (*models.App, error) {

	region := c.DefaultRegion

	if c.Body.Region != "" {
		region = c.Body.Region
	}

	slug := slugify(c.Body.Name)

	// check if bucket already exists
	// existBucketLocation, err := c.S3.CheckOrCreateNewBucket(ctx, slug, region)

	// if err != nil {
	// 	return nil, err
	// }

	app := &models.App{
		ID:          uuid.New().String(),
		OwnerID:     c.User.ID,
		Name:        c.Body.Name,
		Slug:        slug,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		Metadata:    null.NewString("{}", true),
		Bucket:      slug,
		Region:      null.NewString(region, region != ""),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// create bucket on aws s3
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

	// create onboarding api key
	// createApiKeyService := CreateApiKeyService{
	// 	App:           app,
	// 	AppMemberRepo: c.AppMemberRepo,
	// 	ApiKeyRepo:    c.ApiKeyRepo,
	// 	User:          c.User,
	// 	Body: &dto.CreateApiKeyRequestDto{
	// 		Name:        "Onboarding",
	// 		Description: "",
	// 		Access:      models.ApiKeyAccessRead,
	// 		ExpiredAt:   0,
	// 	},
	// }

	// if _, err := createApiKeyService.Run(ctx); err != nil {
	// 	return nil, err
	// }

	//if existBucketLocation != "" {
	// if err = c.Job.Client.Enqueue(job.QueueNameDefault, job.JobNameAppSync, &job.ClientPayload{
	// 	Data: []byte(app.ID),
	// }); err != nil {
	// 	return nil, err
	// }
	//}

	return app, nil
}

func slugify(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
