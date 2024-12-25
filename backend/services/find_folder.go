package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindFolderService struct {
	FolderId   string
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
}

func (c *FindFolderService) Run(ctx context.Context) (*models.Folder, error) {

	folder, err := c.FolderRepo.FindFolderByID(ctx, c.FolderId)

	if err != nil {
		return nil, err
	}

	return folder, nil
}
