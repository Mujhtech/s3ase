package server

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mujhtech/s3ase/api/handler"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/server"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func RegisterServerCommand() *cobra.Command {

	var (
		configFile string
	)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start s3ase server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			_ = godotenv.Load(configFile)

			cfg, err := config.LoadConfig()

			if err != nil {
				panic(err)
			}

			db, err := database.Connect(ctx, cfg)

			if err != nil {
				panic(err)
			}

			defer db.Close()

			router, err := handler.New(cfg, db)

			if err != nil {
				panic(err)
			}

			server := server.New(cfg, router.BuildHandler())

			g, gCtx := errgroup.WithContext(ctx)

			gHTTP, shutdownHTTP := server.ListenAndServe()
			g.Go(gHTTP.Wait)

			<-gCtx.Done()

			stop()

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if sErr := shutdownHTTP(shutdownCtx); sErr != nil {
			}

			err = g.Wait()

			if err != nil {

			}

		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")

	return cmd

}
