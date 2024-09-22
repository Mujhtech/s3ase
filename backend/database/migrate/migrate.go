package migrate

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"
	"github.com/mujhtech/s3ase/database"
	"maragu.dev/migrate"
)

//go:embed postgres/*.sql
var postgres embed.FS

//go:embed sqlite/*.sql
var sqlite embed.FS

const (
	tableName = "migrations"

	postgresDriverName = "postgres"
	postgresSourceDir  = "postgres"

	sqliteDriverName = "sqlite3"
	sqliteSourceDir  = "sqlite"
)

func Migrator(ctx context.Context, db *database.Database) (*migrate.Migrator, error) {
	opts, err := getMigratorOpt(db.GetDB())
	if err != nil {
		return nil, fmt.Errorf("failed to get migrator opt: %w", err)
	}
	return migrate.New(opts), nil
}

func getMigratorOpt(db *sqlx.DB) (migrate.Options, error) {

	opts := migrate.Options{
		FS: postgres,
		DB: db.DB,
	}

	switch db.DriverName() {
	case sqliteDriverName:
		folder, _ := fs.Sub(sqlite, sqliteSourceDir)
		opts.FS = folder
	case postgresDriverName:
		folder, _ := fs.Sub(postgres, postgresSourceDir)
		opts.FS = folder

	default:
		return migrate.Options{}, fmt.Errorf("unsupported driver '%s'", db.DriverName())
	}

	return opts, nil
}
