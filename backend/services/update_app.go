package services

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type UpdateAppService struct {
	AppId   string
	AppRepo store.AppRepository
	User    *models.User
	Body    *dto.CreateAppRequestDto
}

func (c *UpdateAppService) Run(ctx context.Context) error {
	app, err := c.AppRepo.FindAppByID(ctx, c.AppId)

	if err != nil {
		return err
	}

	if app.OwnerID != c.User.ID {
		return fmt.Errorf("user is not owner of the app")
	}

	app.Name = c.Body.Name
	app.Description = null.NewString(c.Body.Description, c.Body.Description != "")

	if err = c.AppRepo.UpdateApp(ctx, app); err != nil {
		return err
	}

	return nil
}
