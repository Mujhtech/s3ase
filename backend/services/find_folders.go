package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindFoldersService struct {
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
}

func (c *FindFoldersService) Run(ctx context.Context) ([]*models.Folder, error) {

	folders, err := c.FolderRepo.FindFolderByAppID(ctx, c.App.ID)

	if err != nil {
		return nil, err
	}

	return folders, nil
}
