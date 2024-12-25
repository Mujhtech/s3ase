package services

import (
	"context"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type UpdateFolderService struct {
	FolderId   string
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
	Body       *dto.CreateFolderRequestDto
}

func (c *UpdateFolderService) Run(ctx context.Context) error {

	folder, err := c.FolderRepo.FindFolderByID(ctx, c.FolderId)

	if err != nil {
		return err
	}

	if err := c.FolderRepo.UpdateFolder(ctx, &models.Folder{
		ID:          folder.ID,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		CreatedBy:   c.User.ID,
		Metadata:    map[string]interface{}{},
	}); err != nil {
		return err
	}

	return nil
}
