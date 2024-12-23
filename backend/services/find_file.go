package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindFileService struct {
	FileID     string
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
}

func (c *FindFileService) Run(ctx context.Context) (*models.File, error) {

	file, err := c.FileRepo.FindFileByID(ctx, c.FileID)

	if err != nil {
		return nil, err
	}

	return file, nil
}
