package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	userBaseTable    = "users"
	userSelectColumn = "id, given_name, COALESCE(display_name, '') AS display_name, email, email_verified, COALESCE(avatar_url, '') AS avatar_url, authentication_method, password, created_at, updated_at, deleted_at"
)

type userRepo struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepo{
		db: db,
	}
}

// CreateUser implements UserRepository.
func (u *userRepo) CreateUser(ctx context.Context, user *models.User) error {

	metadata := "{}"

	if user.Metadata != nil {
		metadataByte, err := json.Marshal(user.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(userBaseTable).
		Columns(
			"given_name",
			"display_name",
			"email",
			"email_verified",
			"avatar_url",
			"authentication_method",
			"password",
			"metadata",
		).
		Values(
			user.Name,
			user.DisplayName,
			user.Email,
			user.EmailVerified,
			user.AvatarUrl,
			user.AuthenticationMethod,
			user.Password,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = u.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create user")
	}

	return nil
}

// FindUserByEmail implements UserRepository.
func (u *userRepo) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	stmt := Builder.
		Select(userSelectColumn).
		From(userBaseTable).
		Where(squirrel.Eq{"email": email}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	user := new(models.User)
	if err := u.db.GetDB().GetContext(ctx, user, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find user by email")
	}

	return user, nil
}

// FindUserByID implements UserRepository.
func (u *userRepo) FindUserByID(ctx context.Context, id string) (*models.User, error) {
	stmt := Builder.
		Select(userSelectColumn).
		From(userBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	user := new(models.User)
	if err := u.db.GetDB().GetContext(ctx, user, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find user by id")
	}

	return user, nil
}

// UpdateUser implements UserRepository.
func (u *userRepo) UpdateUser(ctx context.Context, user *models.User) error {

	stmt := Builder.
		Update(userBaseTable).
		Set("given_name", user.Name).
		Set("display_name", user.DisplayName).
		Set("email", user.Email).
		Set("email_verified", user.EmailVerified).
		Set("avatar_url", user.AvatarUrl).
		Set("authentication_method", user.AuthenticationMethod).
		Set("password", user.Password).
		Set("metadata", "{}").
		Where(squirrel.Eq{"id": user.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = u.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update user")
	}

	return nil
}
