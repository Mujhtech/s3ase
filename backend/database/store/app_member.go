package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	appMemberBaseTable    = "app_members"
	appMemberSelectColumn = "id, user_id, app_id, role, metadata, created_at, updated_at, deleted_at"
)

type appMemberRepo struct {
	db *database.Database
}

func NewAppMemberRepository(db *database.Database) AppMemberRepository {
	return &appMemberRepo{
		db: db,
	}
}

// CreateAppMember implements AppMemberRepository.
func (a *appMemberRepo) CreateAppMember(ctx context.Context, appMember *models.AppMember) error {
	metadata := "{}"

	if appMember.Metadata != nil {
		metadataByte, err := json.Marshal(appMember.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(appMemberBaseTable).
		Columns(
			"user_id",
			"app_id",
			"role",
			"metadata",
		).
		Values(
			appMember.UserId,
			appMember.AppID,
			appMember.Role,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create app member")
	}

	return nil
}

// DeleteAppMember implements AppMemberRepository.
func (a *appMemberRepo) DeleteAppMember(ctx context.Context, id string) error {
	stmt := Builder.
		Update(appMemberBaseTable).
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete app member")
	}

	return nil
}

// FindAppMembersByAppID implements AppMemberRepository.
func (a *appMemberRepo) FindAppMembersByAppID(ctx context.Context, appID string) ([]*models.AppMember, error) {
	stmt := Builder.
		Select(appMemberSelectColumn).
		From(appMemberBaseTable).
		Where(squirrel.Eq{"app_id": appID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	members := []*models.AppMember{}
	if err := a.db.GetDB().SelectContext(ctx, &members, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find app member by app id")
	}

	return members, nil
}

// FindAppMemberByID implements AppMemberRepository.
func (a *appMemberRepo) FindAppMemberByID(ctx context.Context, id string) (*models.AppMember, error) {
	stmt := Builder.
		Select(appMemberSelectColumn).
		From(appMemberBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	member := new(models.AppMember)
	if err := a.db.GetDB().GetContext(ctx, member, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find app member by id")
	}

	return member, nil
}

// FindAppMemberByAppIDAndUserId implements AppMemberRepository.
func (a *appMemberRepo) FindAppMemberByAppIDAndUserId(ctx context.Context, appID string, userId string) (*models.AppMember, error) {
	stmt := Builder.
		Select(appMemberSelectColumn).
		From(appMemberBaseTable).
		Where(squirrel.Eq{"app_id": appID}).
		Where(squirrel.Eq{"user_id": userId}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	member := new(models.AppMember)
	if err := a.db.GetDB().GetContext(ctx, member, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find app member by app id and user id")
	}

	return member, nil
}

// FindAppMembersByUserID implements AppMemberRepository.
func (a *appMemberRepo) FindAppMembersByUserID(ctx context.Context, userID string) ([]*models.AppMember, error) {
	stmt := Builder.
		Select(appMemberSelectColumn).
		From(appMemberBaseTable).
		Where(squirrel.Eq{"user_id": userID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	members := []*models.AppMember{}
	if err := a.db.GetDB().SelectContext(ctx, &members, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find app member by user id")
	}

	return members, nil
}
