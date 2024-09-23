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
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func RegisterServerCommand() *cobra.Command {

	var (
		configFile string
		logLevel   string
	)

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start s3ase server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			switch logLevel {
			case "debug":
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			case "trace":
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}

			zerolog.TimeFieldFormat = time.RFC3339Nano

			// attach logger to context
			logger := log.Logger.With().Logger()
			ctx = logger.WithContext(ctx)

			_ = godotenv.Load(configFile)

			cfg, err := config.LoadConfig()

			if err != nil {
				log.Err(err).Msg("failed to load config")
			}

			db, err := database.Connect(ctx, cfg)

			if err != nil {
				log.Err(err).Msg("failed to connect to database")
			}

			defer db.Close()

			router, err := handler.New(cfg, db)

			if err != nil {
				log.Err(err).Msg("failed to create handler")
			}

			server := server.New(cfg, router.BuildHandler())

			g, gCtx := errgroup.WithContext(ctx)

			gHTTP, shutdownHTTP := server.ListenAndServe()
			g.Go(gHTTP.Wait)

			logger.Info().Msgf("server started on port %d", cfg.Server.Port)

			<-gCtx.Done()

			stop()

			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if sErr := shutdownHTTP(shutdownCtx); sErr != nil {
				log.Err(sErr).Msg("failed to shutdown server gracefully")
			}

			logger.Info().Msg("waiting for all goroutines to finish")
			err = g.Wait()

			if err != nil {
				log.Err(err).Msg("failed to wait for all goroutines to finish")
			}

		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "log level")

	return cmd

}
