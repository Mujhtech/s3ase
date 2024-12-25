package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	folderBaseTable    = "folders"
	folderSelectColumn = "id, app_id, created_by, parent_id, name, description, metadata, created_at, updated_at, deleted_at"
)

type folderRepo struct {
	db *database.Database
}

func NewFolderRepository(db *database.Database) FolderRepository {
	return &folderRepo{
		db: db,
	}
}

// CreateFolder implements FolderRepository.
func (f *folderRepo) CreateFolder(ctx context.Context, folder *models.Folder) error {
	metadata := "{}"

	if folder.Metadata != nil {
		metadataByte, err := json.Marshal(folder.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(folderBaseTable).
		Columns(
			"id",
			"app_id",
			"created_by",
			"parent_id",
			"name",
			"description",
			"metadata",
		).
		Values(
			folder.ID,
			folder.AppID,
			folder.CreatedBy,
			folder.ParentID,
			folder.Name,
			folder.Description,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = f.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create folder")
	}

	return nil
}

// UpdateFolder implements FolderRepository.
func (f *folderRepo) UpdateFolder(ctx context.Context, folder *models.Folder) error {
	stmt := Builder.
		Update(folderBaseTable).
		Set("name", folder.Name).
		Set("description", folder.Description).
		Set("parent_id", folder.ParentID).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": folder.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = f.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update folder")
	}

	return nil
}

// DeleteFolder implements FolderRepository.
func (f *folderRepo) DeleteFolder(ctx context.Context, id string) error {
	stmt := Builder.
		Update(folderBaseTable).
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = f.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete folder")
	}

	return nil
}

// FindFoldersByAppID implements FolderRepository.
func (f *folderRepo) FindFoldersByAppID(ctx context.Context, appID string) ([]*models.Folder, error) {
	stmt := Builder.
		Select(folderSelectColumn).
		From(folderBaseTable).
		Where(squirrel.Eq{"app_id": appID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	folders := []*models.Folder{}
	if err := f.db.GetDB().SelectContext(ctx, &folders, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find folders by app id")
	}

	return folders, nil
}

// FindFolderByID implements FolderRepository.
func (f *folderRepo) FindFolderByID(ctx context.Context, id string) (*models.Folder, error) {
	stmt := Builder.
		Select(folderSelectColumn).
		From(folderBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	folder := new(models.Folder)
	if err := f.db.GetDB().GetContext(ctx, folder, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find folder by id")
	}

	return folder, nil
}
