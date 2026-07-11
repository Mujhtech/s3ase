package services

import (
	"context"
	"errors"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

func requireSameApp(expectedAppID, actualAppID string) error {
	if expectedAppID == "" || actualAppID == "" || expectedAppID != actualAppID {
		return errs.ErrNotAuthorized
	}
	return nil
}

func requireAppOwner(
	ctx context.Context,
	repo store.AppMemberRepository,
	appID string,
	userID string,
) error {
	member, err := repo.FindAppMemberByAppIDAndUserId(ctx, appID, userID)
	if errors.Is(err, store.ErrNotFound) || member == nil {
		return errs.ErrNotAuthorized
	}
	if err != nil {
		return err
	}
	if member.Role != models.AppMemberRoleOwner {
		return errs.ErrNotAuthorized
	}
	return nil
}
