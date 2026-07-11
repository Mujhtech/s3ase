package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	if c.AuthUser == nil || c.AuthUser.Metadata == nil {
		return nil, fmt.Errorf("authentication provider returned an invalid user")
	}

	email := strings.TrimSpace(strings.ToLower(c.AuthUser.Metadata.Email))
	for _, candidate := range c.AuthUser.Emails {
		if candidate.Primary {
			email = strings.TrimSpace(strings.ToLower(candidate.Email))
			break
		}
		if email == "" {
			email = strings.TrimSpace(strings.ToLower(candidate.Email))
		}
	}
	if email == "" {
		return nil, fmt.Errorf("authentication provider did not return an email address")
	}

	dst := &models.User{
		Email:                email,
		EmailVerified:        c.AuthUser.Metadata.EmailVerified,
		Name:                 c.AuthUser.Metadata.Name,
		DisplayName:          c.AuthUser.Metadata.Username,
		AvatarUrl:            c.AuthUser.Metadata.AvatarUrl,
		AuthenticationMethod: models.AuthMethod(c.AuthUser.AuthenticationMethod),
		Metadata:             c.AuthUser.Metadata,
		Password:             null.NewString("", true),
	}
	for _, candidate := range c.AuthUser.Emails {
		if strings.EqualFold(candidate.Email, email) {
			dst.EmailVerified = candidate.Verified
			break
		}
	}
	if !dst.EmailVerified {
		return nil, fmt.Errorf("authentication provider did not verify the email address")
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
