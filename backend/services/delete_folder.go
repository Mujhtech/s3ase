package services

import (
	"context"
	"fmt"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
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
	if err := requireSameApp(c.App.ID, folder.AppID); err != nil {
		return err
	}
	files, err := c.FileRepo.FindFilesByFolderID(ctx, folder.ID)
	if err != nil {
		return err
	}
	if len(files) > 0 {
		return fmt.Errorf("%w: folder must be empty before deletion", errs.ErrConflict)
	}

	if err := c.FolderRepo.DeleteFolder(ctx, folder.ID); err != nil {
		return err
	}

	return nil
}
