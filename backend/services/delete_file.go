package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FileObjectStore interface {
	DeleteFile(ctx context.Context, bucket, key, region string) error
}

type DeleteFileService struct {
	App      *models.App
	FileID   string
	FileRepo store.FileRepository
	S3       FileObjectStore
}

func (s *DeleteFileService) Run(ctx context.Context) (*models.File, error) {
	file, err := s.FileRepo.FindFileByID(ctx, s.FileID)
	if err != nil {
		return nil, err
	}
	if err := requireSameApp(s.App.ID, file.AppID); err != nil {
		return nil, err
	}
	if err := s.S3.DeleteFile(ctx, s.App.Bucket, file.ObjectKey(), s.App.Region.String); err != nil {
		return nil, err
	}
	if err := s.FileRepo.DeleteFile(ctx, file.ID); err != nil {
		return nil, err
	}
	return file, nil
}
