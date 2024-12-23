package store

import (
	"context"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/models"
)

const (
	domainBaseTable    = "domains"
	domainSelectColumn = "id, app_id, created_by, domain, cname_record, txt_record, status, metadata, created_at, updated_at, deleted_at"
)

type domainRepo struct {
	db *database.Database
}

func NewDomainRepository(db *database.Database) DomainRepository {
	return &domainRepo{
		db: db,
	}
}

// CreateDomain implements DomainRepository.
func (d *domainRepo) CreateDomain(ctx context.Context, domain *models.Domain) error {
	metadata := "{}"

	if domain.Metadata != nil {
		metadataByte, err := json.Marshal(domain.Metadata)

		if err != nil {
			return err
		}

		metadata = string(metadataByte)
	}

	stmt := Builder.
		Insert(domainBaseTable).
		Columns(
			"id",
			"app_id",
			"created_by",
			"domain",
			"cname_record",
			"txt_record",
			"status",
			"metadata",
		).
		Values(
			domain.ID,
			domain.AppID,
			domain.CreatedBy,
			domain.Domain,
			domain.CnameRecord,
			domain.TxtRecord,
			domain.Status,
			metadata,
		)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = d.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to create domain")
	}

	return nil
}

// DeleteDomain implements DomainRepository.
func (d *domainRepo) DeleteDomain(ctx context.Context, id string) error {

	stmt := Builder.
		Update(domainBaseTable).
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = d.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to delete domain")
	}

	return nil
}

// UpdateDomain implements DomainRepository.
func (d *domainRepo) UpdateDomain(ctx context.Context, domain *models.Domain) error {

	stmt := Builder.
		Update(domainBaseTable).
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("domain", domain.Domain).
		Set("cname_record", domain.CnameRecord).
		Set("status", domain.Status).
		Set("txt_record", domain.TxtRecord).
		Where(squirrel.Eq{"id": domain.ID}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return err
	}

	_, err = d.db.GetDB().ExecContext(ctx, sql, args...)

	if err != nil {
		return ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to update domain")
	}

	return nil
}

// FindDomainByID implements DomainRepository.
func (d *domainRepo) FindDomainByID(ctx context.Context, id string) (*models.Domain, error) {
	stmt := Builder.
		Select(domainSelectColumn).
		From(domainBaseTable).
		Where(squirrel.Eq{"id": id}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	domain := new(models.Domain)
	if err := d.db.GetDB().GetContext(ctx, domain, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find domain by id")
	}

	return domain, nil
}

// FindDomainByDomain implements DomainRepository.
func (d *domainRepo) FindDomainByDomain(ctx context.Context, name string) (*models.Domain, error) {
	stmt := Builder.
		Select(domainSelectColumn).
		From(domainBaseTable).
		Where(squirrel.Eq{"domain": name}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	domain := new(models.Domain)
	if err := d.db.GetDB().GetContext(ctx, domain, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find domain by id")
	}

	return domain, nil
}

// FindDomainByAppID implements DomainRepository.
func (d *domainRepo) FindDomainByAppID(ctx context.Context, appId string) (*models.Domain, error) {
	stmt := Builder.
		Select(domainSelectColumn).
		From(domainBaseTable).
		Where(squirrel.Eq{"app_id": appId}).
		Where(excludeDeleted)

	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	domain := new(models.Domain)
	if err := d.db.GetDB().GetContext(ctx, domain, sql, args...); err != nil {
		return nil, ProcessSQLErrorfWithCtx(ctx, sql, err, "failed to find domain by app id")
	}

	return domain, nil
}
