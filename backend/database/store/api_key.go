package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	apiKeyBaseTable    = "api_keys"
	apiKeySelectColumn = "id, app_id, created_by, name, description, access, expired_at, last_used, COALESCE(key_hash, '') AS key_hash, COALESCE(key_prefix, '') AS key_prefix, COALESCE(last_four, '') AS last_four, metadata, created_at, updated_at, deleted_at"
)

type apiKeyRepo struct {
	db *database.Database
}

func NewApiKeyRepository(db *database.Database) ApiKeyRepository {
	return &apiKeyRepo{
		db: db,
	}
}

// CreateApiKey implements ApiKeyRepository.
func (a *apiKeyRepo) CreateApiKey(ctx context.Context, apiKey *models.ApiKey) error {
	metadata := "{}"

	if apiKey.Metadata != nil {
		metadataByte, err := json.Marshal(apiKey.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(apiKeyBaseTable).
		Columns(
			"id",
			"app_id",
			"created_by",
			"name",
			"description",
			"access",
			"expired_at",
			"metadata",
			"key_hash",
			"key_prefix",
			"last_four",
		).
		Values(
			apiKey.ID,
			apiKey.AppID,
			apiKey.CreatedBy,
			apiKey.Name,
			apiKey.Description,
			apiKey.Access,
			apiKey.ExpiredAt,
			metadata,
			apiKey.KeyHash,
			apiKey.KeyPrefix,
			apiKey.LastFour,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create api key")
	}

	return nil
}

func (a *apiKeyRepo) FindApiKeyByHash(ctx context.Context, hash string) (*models.ApiKey, error) {
	stmt := Builder.Select(apiKeySelectColumn).From(apiKeyBaseTable).
		Where(squirrel.Eq{"key_hash": hash}).Where(excludeDeleted)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	apiKey := new(models.ApiKey)
	if err := a.db.GetDB().GetContext(ctx, apiKey, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find api key by hash")
	}
	return apiKey, nil
}

func (a *apiKeyRepo) TouchApiKey(ctx context.Context, id string) error {
	stmt := Builder.Update(apiKeyBaseTable).Set("last_used", squirrel.Expr("NOW()")).
		Set("updated_at", squirrel.Expr("NOW()")).Where(squirrel.Eq{"id": id}).Where(excludeDeleted)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}
	if _, err := a.db.GetDB().ExecContext(ctx, sql, args...); err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update api key usage")
	}
	return nil
}

// DeleteApiKey implements ApiKeyRepository.
func (a *apiKeyRepo) DeleteApiKey(ctx context.Context, id string) error {

	stmt := Builder.
		Update(apiKeyBaseTable).
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
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete api key")
	}

	return nil
}

// UpdateApiKey implements ApiKeyRepository.
func (a *apiKeyRepo) UpdateApiKey(ctx context.Context, apiKey *models.ApiKey) error {

	stmt := Builder.
		Update(apiKeyBaseTable).
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("name", apiKey.Name).
		Set("description", apiKey.Description).
		Set("access", apiKey.Access).
		Set("expired_at", apiKey.ExpiredAt).
		Where(squirrel.Eq{"id": apiKey.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update api key")
	}

	return nil
}

// FindApiKeyByID implements ApiKeyRepository.
func (a *apiKeyRepo) FindApiKeyByID(ctx context.Context, id string) (*models.ApiKey, error) {
	stmt := Builder.
		Select(apiKeySelectColumn).
		From(apiKeyBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	apiKey := new(models.ApiKey)
	if err := a.db.GetDB().GetContext(ctx, apiKey, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find api key by id")
	}

	return apiKey, nil
}

// FindApiKeysByAppID implements ApiKeyRepository.
func (a *apiKeyRepo) FindApiKeysByAppID(ctx context.Context, appId string) ([]*models.ApiKey, error) {
	stmt := Builder.
		Select(apiKeySelectColumn).
		From(apiKeyBaseTable).
		Where(squirrel.Eq{"app_id": appId}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	apiKeys := []*models.ApiKey{}
	if err := a.db.GetDB().SelectContext(ctx, &apiKeys, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find api keys by app id")
	}

	return apiKeys, nil
}
