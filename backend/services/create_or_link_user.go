package services

import (
	"context"
	"errors"

	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/auth"
)

type CreateOrLinkUserService struct {
	UserRepo store.UserRepository
	AuthUser *auth.User
}

func (c *CreateOrLinkUserService) Run(ctx context.Context) (*models.User, error) {

	dst := &models.User{
		Email:                c.AuthUser.Emails[0].Email,
		EmailVerified:        c.AuthUser.Emails[0].Verified,
		Name:                 c.AuthUser.Metadata.Name,
		DisplayName:          c.AuthUser.Metadata.Username,
		AvatarUrl:            c.AuthUser.Metadata.AvatarUrl,
		AuthenticationMethod: models.AuthMethod(c.AuthUser.AuthenticationMethod),
		Metadata:             c.AuthUser.Metadata,
		Password:             null.NewString("", true),
	}

	user, err := c.UserRepo.FindUserByEmail(ctx, dst.Email)

	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}

	if user == nil {
		err = c.UserRepo.CreateUser(ctx, dst)

		if err != nil {
			return nil, err
		}
	} else {
		dst.ID = user.ID
		err = c.UserRepo.UpdateUser(ctx, dst)

		if err != nil {
			return nil, err
		}
	}

	return c.UserRepo.FindUserByEmail(ctx, dst.Email)
}
