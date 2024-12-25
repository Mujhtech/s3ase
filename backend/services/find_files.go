package services

import (
	"context"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindFilesService struct {
	App        *models.App
	FolderRepo store.FolderRepository
	FileRepo   store.FileRepository
	User       *models.User
	Query      *dto.FileQueryDto
}

func (c *FindFilesService) Run(ctx context.Context) ([]*models.File, error) {

	files, err := c.FileRepo.FindFilesByAppIDWithQuery(ctx, c.App.ID, c.Query)

	if err != nil {
		return nil, err
	}

	return files, nil
}
