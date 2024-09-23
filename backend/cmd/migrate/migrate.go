package migrate

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/migrate"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func RegisterMigrateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "s3ase migration",
		Run: func(cmd *cobra.Command, args []string) {

		},
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

			_ = godotenv.Load(configFile)

			cfg, err := config.LoadConfig()

			if err != nil {
				log.Err(err).Msg("failed to load config")
			}

			ctx := context.Background()

			db, err := database.Connect(ctx, cfg)

			if err != nil {
				log.Err(err).Msg("failed to connect to database")
			}

			defer db.Close()

			migrator, err := migrate.Migrator(ctx, cfg, db)

			if err != nil {
				log.Err(err).Msg("failed to create migrator")
			}

			err = migrator.MigrateUp(ctx)

			if err != nil {
				log.Err(err).Msg("failed to migrate")
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

			_ = godotenv.Load(configFile)

			cfg, err := config.LoadConfig()

			if err != nil {
				log.Err(err).Msg("failed to load config")
			}

			ctx := context.Background()

			db, err := database.Connect(ctx, cfg)

			if err != nil {
				log.Err(err).Msg("failed to connect to database")
			}

			defer db.Close()

			migrator, err := migrate.Migrator(ctx, cfg, db)

			if err != nil {
				log.Err(err).Msg("failed to create migrator")
			}

			err = migrator.MigrateDown(ctx)

			if err != nil {
				log.Err(err).Msg("failed to migrate")
			}

			log.Info().Msg("migration completed")
		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")

	return cmd

}
