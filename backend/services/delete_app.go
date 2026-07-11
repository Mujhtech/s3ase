package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

type AppObjectStore interface {
	DeleteBucket(ctx context.Context, bucket, region string) error
}

type DeleteAppService struct {
	AppID       string
	ConfirmName string
	AppRepo     store.AppRepository
	User        *models.User
	S3          AppObjectStore
}

func (s *DeleteAppService) Run(ctx context.Context) error {
	app, err := s.AppRepo.FindAppByID(ctx, s.AppID)
	if err != nil {
		return err
	}
	if s.User == nil || app.OwnerID != s.User.ID {
		return errs.ErrNotAuthorized
	}
	if strings.TrimSpace(s.ConfirmName) != app.Name {
		return fmt.Errorf("%w: app name does not match", errs.ErrInvalidInput)
	}
	if s.S3 == nil {
		return fmt.Errorf("object storage is not configured")
	}
	if err := s.S3.DeleteBucket(ctx, app.Bucket, app.Region.String); err != nil {
		return fmt.Errorf("failed to delete app storage: %w", err)
	}
	return s.AppRepo.DeleteApp(ctx, app.ID)
}
