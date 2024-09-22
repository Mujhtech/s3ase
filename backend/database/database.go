package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/mujhtech/s3ase/config"
)

const (
	dbTx string = "db_tx"
)

type Database struct {
	dbx *sqlx.DB
}

func Connect(ctx context.Context, cfg *config.Config) (*Database, error) {
	db, err := sql.Open(string(cfg.Database.Driver), cfg.Database.BuildDsn())
	if err != nil {
		return nil, fmt.Errorf("failed to open the db: %w", err)
	}

	dbx := sqlx.NewDb(db, cfg.Database.BuildDsn())

	if err = pingDatabase(ctx, dbx); err != nil {
		return nil, fmt.Errorf("failed to ping the db: %w", err)
	}

	return &Database{
		dbx: dbx,
	}, nil
}

func (d *Database) GetDB() *sqlx.DB {
	return d.dbx
}

func (d *Database) Close() error {
	return d.dbx.Close()
}

func (d *Database) StartTx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return d.dbx.BeginTxx(ctx, opts)
}

func (d *Database) RollbackTx(tx *sqlx.Tx, err error) error {
	if err != nil {
		return tx.Rollback()
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		return tx.Rollback()
	}

	return nil
}

func (d *Database) GetTx(ctx context.Context, db *sqlx.DB) (*sqlx.Tx, error) {
	ctxTx, ok := ctx.Value(dbTx).(*sqlx.Tx)
	if ok {
		return ctxTx, nil
	}

	tx, err := db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to start a transaction: %w", err)
	}

	return tx, nil
}

func pingDatabase(ctx context.Context, db *sqlx.DB) error {
	var err error
	for i := 1; i <= 30; i++ {
		err = db.PingContext(ctx)

		if errors.Is(err, context.Canceled) {
			return err
		}

		if err == nil {
			return nil
		}

		time.Sleep(time.Second)
	}

	return fmt.Errorf("all 30 tries failed, last failure: %w", err)
}
