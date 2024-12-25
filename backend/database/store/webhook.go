package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	webhookBaseTable    = "webhooks"
	webhookSelectColumn = "id, app_id, created_by, name, description, url, metadata, created_at, updated_at, deleted_at"
)

type webhookRepo struct {
	db *database.Database
}

func NewWebhookRepository(db *database.Database) WebhookRepository {
	return &webhookRepo{
		db: db,
	}
}

// CreateWebhook implements WebhookRepository.
func (a *webhookRepo) CreateWebhook(ctx context.Context, apiKey *models.Webhook) error {
	metadata := "{}"

	if apiKey.Metadata != nil {
		metadataByte, err := json.Marshal(apiKey.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(webhookBaseTable).
		Columns(
			"id",
			"app_id",
			"created_by",
			"name",
			"description",
			"url",
			"metadata",
		).
		Values(
			apiKey.ID,
			apiKey.AppID,
			apiKey.CreatedBy,
			apiKey.Name,
			apiKey.Description,
			apiKey.URL,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create webhook")
	}

	return nil
}

// DeleteWebhook implements WebhookRepository.
func (a *webhookRepo) DeleteWebhook(ctx context.Context, id string) error {

	stmt := Builder.
		Update(webhookBaseTable).
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
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete webhook")
	}

	return nil
}

// UpdateWebhook implements WebhookRepository.
func (a *webhookRepo) UpdateWebhook(ctx context.Context, apiKey *models.Webhook) error {

	stmt := Builder.
		Update(webhookBaseTable).
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("name", apiKey.Name).
		Set("description", apiKey.Description).
		Set("url", apiKey.URL).
		Where(squirrel.Eq{"id": apiKey.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update webhook")
	}

	return nil
}

// FindWebhookByID implements WebhookRepository.
func (a *webhookRepo) FindWebhookByID(ctx context.Context, id string) (*models.Webhook, error) {
	stmt := Builder.
		Select(webhookSelectColumn).
		From(webhookBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	webhook := new(models.Webhook)
	if err := a.db.GetDB().GetContext(ctx, webhook, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find webhook by id")
	}

	return webhook, nil
}

// FindWebhooksByAppID implements WebhookRepository.
func (a *webhookRepo) FindWebhooksByAppID(ctx context.Context, appId string) ([]*models.Webhook, error) {
	stmt := Builder.
		Select(webhookSelectColumn).
		From(webhookBaseTable).
		Where(squirrel.Eq{"app_id": appId}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	webhooks := []*models.Webhook{}
	if err := a.db.GetDB().SelectContext(ctx, &webhooks, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find webhooks by app id")
	}

	return webhooks, nil
}
