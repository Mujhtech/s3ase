package main

import (
	"log"
	"os"

	"github.com/mujhtech/s3ase/cmd/migrate"
	"github.com/mujhtech/s3ase/cmd/server"
	"github.com/mujhtech/s3ase/cmd/version"
	"github.com/spf13/cobra"
)

func main() {
	err := os.Setenv("TZ", "")
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:   "s3ase",
		Short: "Simplifying S3 usage through open source.",
	}

	// Add subcommands
	cmd.AddCommand(server.RegisterServerCommand())
	cmd.AddCommand(version.RegisterVersionCommand())
	cmd.AddCommand(migrate.RegisterMigrateCommand())

	err = cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
