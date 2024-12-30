package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	fileBaseTable    = "files"
	fileSelectColumn = "id, uploaded_by, name, mime_type, extension, app_id, folder_id, metadata, size, is_public, public_id, created_at, updated_at, deleted_at"
)

type fileRepo struct {
	db *database.Database
}

func NewFileRepository(db *database.Database) FileRepository {
	return &fileRepo{
		db: db,
	}
}

// CreateFile implements FileRepository.
func (f *fileRepo) CreateFile(ctx context.Context, file *models.File) error {
	metadata := "{}"

	if file.Metadata != nil {
		metadataByte, err := json.Marshal(file.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(fileBaseTable).
		Columns(
			"id",
			"uploaded_by",
			"app_id",
			"folder_id",
			"metadata",
			"name",
			"mime_type",
			"extension",
			"size",
			"is_public",
			"public_id",
			"status",
		).
		Values(
			file.ID,
			file.UploadedBy,
			file.AppID,
			file.FolderID,
			metadata,
			file.Name,
			file.MimeType,
			file.Extension,
			file.Size,
			file.IsPublic,
			file.PublicID,
			file.Status,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = f.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create file")
	}

	return nil
}

// UpdateFile implements FileRepository.
func (f *fileRepo) UpdateFile(ctx context.Context, file *models.File) error {
	stmt := Builder.
		Update(fileBaseTable)

	if file.Name != "" {
		stmt = stmt.Set("name", file.Name)
	}

	if file.Status != "" {
		stmt = stmt.Set("status", file.Status)
	}

	if file.FolderID != "" {
		stmt = stmt.Set("folder_id", file.FolderID)
	}

	if file.Size != 0 {
		stmt = stmt.Set("size", file.Size)
	}

	if file.Metadata != nil {
		metadataByte, err := json.Marshal(file.Metadata)

		if err != nil {
			return err
		}

		stmt = stmt.Set("metadata", string(metadataByte))
	}

	stmt = stmt.Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": file.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = f.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update file")
	}

	return nil
}

// DeleteFile implements FileRepository.
func (f *fileRepo) DeleteFile(ctx context.Context, id string) error {
	stmt := Builder.
		Update(fileBaseTable).
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
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete file")
	}

	return nil
}

// FindFileByID implements FileRepository.
func (f *fileRepo) FindFileByID(ctx context.Context, id string) (*models.File, error) {
	stmt := Builder.
		Select(fileSelectColumn).
		From(fileBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	file := new(models.File)
	if err := f.db.GetDB().GetContext(ctx, file, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find file by id")
	}

	return file, nil
}

// FindFilesByAppID implements FileRepository.
func (f *fileRepo) FindFilesByAppID(ctx context.Context, appID string) ([]*models.File, error) {
	stmt := Builder.
		Select(fileSelectColumn).
		From(fileBaseTable).
		Where(squirrel.Eq{"app_id": appID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	files := []*models.File{}
	if err := f.db.GetDB().SelectContext(ctx, &files, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find files by app id")
	}

	return files, nil
}

// FindFilesByAppIDWithQuery implements FileRepository.
func (f *fileRepo) FindFilesByAppIDWithQuery(ctx context.Context, appID string, query *dto.FileQueryDto) ([]*models.File, error) {
	stmt := Builder.
		Select(fileSelectColumn).
		From(fileBaseTable).
		Where(squirrel.Eq{"app_id": appID}).
		Where(excludeDeleted)

	if query.FolderID != "" {
		stmt = stmt.Where(squirrel.Eq{"folder_id": query.FolderID})
	}

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	files := []*models.File{}
	if err := f.db.GetDB().SelectContext(ctx, &files, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find files by app id")
	}

	return files, nil
}

// FindFilesByFolderID implements FileRepository.
func (f *fileRepo) FindFilesByFolderID(ctx context.Context, folderID string) ([]*models.File, error) {
	stmt := Builder.
		Select(fileSelectColumn).
		From(fileBaseTable).
		Where(squirrel.Eq{"folder_id": folderID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	files := []*models.File{}
	if err := f.db.GetDB().SelectContext(ctx, &files, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find files by folder id")
	}

	return files, nil
}

// FindFilesByUserID implements FileRepository.
func (f *fileRepo) FindFilesByUserID(ctx context.Context, userID string) ([]*models.File, error) {
	stmt := Builder.
		Select(fileSelectColumn).
		From(fileBaseTable).
		Where(squirrel.Eq{"created_by": userID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	files := []*models.File{}
	if err := f.db.GetDB().SelectContext(ctx, &files, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find files by user id")
	}

	return files, nil
}
