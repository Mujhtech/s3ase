package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type CreateFolderService struct {
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
	Body       *dto.CreateFolderRequestDto
}

func (c *CreateFolderService) Run(ctx context.Context) (*models.Folder, error) {

	id := uuid.New().String()

	if err := c.FolderRepo.CreateFolder(ctx, &models.Folder{
		ID:          id,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		CreatedBy:   c.User.ID,
		Metadata:    map[string]interface{}{},
	}); err != nil {
		return nil, err
	}

	return c.FolderRepo.FindFolderByID(ctx, id)
}
