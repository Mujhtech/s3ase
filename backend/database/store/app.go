package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	appBaseTable    = "apps"
	appSelectColumn = "id, owner_id, name, slug, description, bucket, region, metadata, created_at, updated_at, deleted_at"
)

type appRepo struct {
	db *database.Database
}

func NewAppRepository(db *database.Database) AppRepository {
	return &appRepo{
		db: db,
	}
}

// CreateApp implements AppRepository.
func (a *appRepo) CreateApp(ctx context.Context, app *models.App) error {
	metadata := "{}"

	if app.Metadata != nil {
		metadataByte, err := json.Marshal(app.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(appBaseTable).
		Columns(
			"id",
			"owner_id",
			"name",
			"slug",
			"description",
			"bucket",
			"region",
			"metadata",
		).
		Values(
			app.ID,
			app.OwnerID,
			app.Name,
			app.Slug,
			app.Description,
			app.Bucket,
			app.Region,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create app")
	}

	return nil
}

// UpdateApp implements AppRepository.
func (a *appRepo) UpdateApp(ctx context.Context, app *models.App) error {
	stmt := Builder.
		Update(appBaseTable).
		Set("name", app.Name).
		Set("description", app.Description).
		Set("bucket", app.Bucket).
		Set("region", app.Region).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": app.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = a.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update app")
	}

	return nil
}

// DeleteApp implements AppRepository.
func (a *appRepo) DeleteApp(ctx context.Context, id string) error {
	stmt := Builder.
		Update(appBaseTable).
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
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete app")
	}

	return nil
}

// FindAppByID implements AppRepository.
func (a *appRepo) FindAppByID(ctx context.Context, id string) (*models.App, error) {
	stmt := Builder.
		Select(appSelectColumn).
		From(appBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	app := new(models.App)
	if err := a.db.GetDB().GetContext(ctx, app, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find app by id")
	}

	return app, nil
}

// FindAppsByUserID implements AppRepository.
func (a *appRepo) FindAppsByUserID(ctx context.Context, userID string) ([]*models.App, error) {
	stmt := Builder.
		Select(appSelectColumn).
		From(appBaseTable).
		Where(squirrel.Eq{"owner_id": userID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	apps := []*models.App{}
	if err := a.db.GetDB().SelectContext(ctx, &apps, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find apps by user id")
	}

	return apps, nil
}
