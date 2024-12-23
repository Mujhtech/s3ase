package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	tokenBaseTable    = "tokens"
	tokenSelectColumn = "id, value, is_app, type, expired_at, issued_at, metadata, created_at, updated_at, deleted_at"
)

type tokenRepo struct {
	db *database.Database
}

func NewTokenRepository(db *database.Database) TokenRepository {
	return &tokenRepo{
		db: db,
	}
}

// CreateToken implements TokenRepository.
func (t *tokenRepo) CreateToken(ctx context.Context, token *models.Token) error {
	metadata := "{}"

	if token.Metadata != nil {
		metadataByte, err := json.Marshal(token.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(tokenBaseTable).
		Columns(
			"id",
			"value",
			"is_app",
			"type",
			"expired_at",
			"issued_at",
			"metadata",
		).
		Values(
			token.ID,
			token.Value,
			token.IsApp,
			token.Type,
			token.ExpiredAt,
			token.IssuedAt,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = t.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create token")
	}

	return nil
}

// DeleteToken implements TokenRepository.
func (t *tokenRepo) DeleteToken(ctx context.Context, id string) error {
	stmt := Builder.
		Update(tokenBaseTable).
		Set("deleted_at", "NOW()").
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = t.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete token")
	}

	return nil
}

// FindTokenByID implements TokenRepository.
func (t *tokenRepo) FindTokenByID(ctx context.Context, id string) (*models.Token, error) {
	stmt := Builder.
		Select(tokenSelectColumn).
		From(tokenBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	token := new(models.Token)
	if err := t.db.GetDB().GetContext(ctx, token, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find token by id")
	}

	return token, nil
}
