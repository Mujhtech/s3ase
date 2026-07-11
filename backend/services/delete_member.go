package services

import (
	"context"
	"fmt"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

type DeleteMemberService struct {
	App           *models.App
	User          *models.User
	MemberID      string
	AppMemberRepo store.AppMemberRepository
}

func (s *DeleteMemberService) Run(ctx context.Context) error {
	if err := requireAppOwner(ctx, s.AppMemberRepo, s.App.ID, s.User.ID); err != nil {
		return err
	}
	member, err := s.AppMemberRepo.FindAppMemberByID(ctx, s.MemberID)
	if err != nil {
		return err
	}
	if err := requireSameApp(s.App.ID, member.AppID); err != nil {
		return err
	}
	if member.Role == models.AppMemberRoleOwner || member.UserId == s.App.OwnerID {
		return fmt.Errorf("%w: the app owner cannot be removed", errs.ErrInvalidInput)
	}
	return s.AppMemberRepo.DeleteAppMember(ctx, member.ID)
}
