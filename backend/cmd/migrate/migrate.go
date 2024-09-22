package migrate

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/migrate"
	"github.com/spf13/cobra"
)

func RegisterMigrateCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "S3ase migration",
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
				panic(err)
			}

			ctx := context.Background()

			db, err := database.Connect(ctx, cfg)

			if err != nil {
				panic(err)
			}

			defer db.Close()

			migrator, err := migrate.Migrator(ctx, db)

			if err != nil {
				panic(err)
			}

			err = migrator.MigrateUp(ctx)

			if err != nil {
				panic(err)
			}
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
				panic(err)
			}

			ctx := context.Background()

			db, err := database.Connect(ctx, cfg)

			if err != nil {
				panic(err)
			}

			defer db.Close()

			migrator, err := migrate.Migrator(ctx, db)

			if err != nil {
				panic(err)
			}

			err = migrator.MigrateDown(ctx)

			if err != nil {
				panic(err)
			}

		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")

	return cmd

}
