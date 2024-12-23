package services

import (
	"context"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type DeleteFolderService struct {
	FolderId   string
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
	Body       *dto.CreateFolderRequestDto
}

func (c *DeleteFolderService) Run(ctx context.Context) error {

	folder, err := c.FolderRepo.FindFolderByID(ctx, c.FolderId)

	if err != nil {
		return err
	}

	if err := c.FolderRepo.DeleteFolder(ctx, folder.ID); err != nil {
		return err
	}

	return nil
}
