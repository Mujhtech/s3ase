package services

import (
	"context"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type FindMembersService struct {
	App           *models.App
	AppRepo       store.AppRepository
	AppMemberRepo store.AppMemberRepository
	UserRepo      store.UserRepository
	User          *models.User
}

func (c *FindMembersService) Run(ctx context.Context) ([]*models.AppMember, error) {

	appMembers, err := c.AppMemberRepo.FindAppMembersByAppID(ctx, c.App.ID)

	if err != nil {
		return nil, err
	}

	for _, appMember := range appMembers {
		user, err := c.UserRepo.FindUserByID(ctx, appMember.UserId)

		if err != nil {
			return nil, err
		}

		appMember.User = user
	}

	return appMembers, nil
}
