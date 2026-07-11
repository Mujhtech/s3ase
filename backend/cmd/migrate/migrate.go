package migrate

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/migrate"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Migration string

const (
	MigrationUp   Migration = "up"
	MigrationDown Migration = "down"
)

func RegisterMigrateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "s3ase migration",
		// Run: func(cmd *cobra.Command, args []string) {

		// },
	}

	cmd.AddCommand(addUpCommand())
	cmd.AddCommand(addDownCommand())

	return cmd

}

func addUpCommand() *cobra.Command {

	var (
		configFile string
	)

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := migration(configFile, MigrationUp)

			if err != nil {
				log.Err(err).Msg("failed to migrate")
				return
			}

			log.Info().Msg("migration completed")
		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")

	return cmd
}

func addDownCommand() *cobra.Command {

	var (
		configFile string
	)

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := migration(configFile, MigrationDown)

			if err != nil {
				log.Err(err).Msg("failed to migrate")
				return
			}

			log.Info().Msg("migration completed")
		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")

	return cmd

}

func migration(configFile string, migration Migration) error {

	_ = godotenv.Load(configFile)

	cfg, err := config.LoadConfig()

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, cfg)

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("failed to close migration database connection")
		}
	}()

	migrator, err := migrate.Migrator(ctx, cfg, db)

	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	switch migration {
	case MigrationUp:
		err = migrator.MigrateUp(ctx)
		if err != nil {
			return fmt.Errorf("failed to migrate: %w", err)
		}
		return nil
	case MigrationDown:
		err = migrator.MigrateDown(ctx)

		if err != nil {
			return fmt.Errorf("failed to migrate: %w", err)
		}
		return nil
	}
	return nil
}
