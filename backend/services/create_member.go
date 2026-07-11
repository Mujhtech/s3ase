package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

type CreateMemberService struct {
	App           *models.App
	User          *models.User
	Body          *dto.CreateMemberRequestDto
	UserRepo      store.UserRepository
	AppMemberRepo store.AppMemberRepository
}

func (s *CreateMemberService) Run(ctx context.Context) (*models.AppMember, error) {
	if s.App == nil || s.User == nil || s.Body == nil {
		return nil, fmt.Errorf("%w: missing member context", errs.ErrInvalidInput)
	}
	if err := requireAppOwner(ctx, s.AppMemberRepo, s.App.ID, s.User.ID); err != nil {
		return nil, err
	}
	email := strings.ToLower(strings.TrimSpace(s.Body.Email))
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", errs.ErrInvalidInput)
	}
	if s.Body.Role == "" {
		s.Body.Role = models.AppMemberRoleMember
	}
	if s.Body.Role != models.AppMemberRoleMember {
		return nil, fmt.Errorf("%w: only the member role can be assigned", errs.ErrInvalidInput)
	}
	user, err := s.UserRepo.FindUserByEmail(ctx, email)
	if errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("%w: this person must sign in to S3ase before they can be added", store.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	if _, err := s.AppMemberRepo.FindAppMemberByAppIDAndUserId(ctx, s.App.ID, user.ID); err == nil {
		return nil, fmt.Errorf("%w: user is already a member", errs.ErrConflict)
	} else if !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}
	member := &models.AppMember{
		AppID: s.App.ID, UserId: user.ID, Role: s.Body.Role,
		Metadata: models.Metadata{}, User: user,
	}
	if err := s.AppMemberRepo.CreateAppMember(ctx, member); err != nil {
		return nil, err
	}
	created, err := s.AppMemberRepo.FindAppMemberByAppIDAndUserId(ctx, s.App.ID, user.ID)
	if err != nil {
		return nil, err
	}
	created.User = user
	return created, nil
}
