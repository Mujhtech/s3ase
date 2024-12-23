package server

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/mujhtech/s3ase/api"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/http"
	"github.com/mujhtech/s3ase/internal/pkg/pubsub"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/internal/redis"
	"github.com/mujhtech/s3ase/job"
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

			err := startServer(configFile, logLevel)

			if err != nil {
				log.Err(err).Msg("failed to start server")
			}

		},
	}

	cmd.Flags().StringVar(&configFile, "config", config.DefaultConfigFilePath, "configuration file")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "log level")

	return cmd

}

func startServer(configFile string, logLevel string) error {

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
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := database.Connect(ctx, cfg)

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	redis, err := redis.NewRedis(cfg)

	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	store := store.NewStore(db)

	job, err := job.NewJob(cfg, redis)

	if err != nil {
		return fmt.Errorf("failed to initialize job: %w", err)
	}

	s3, err := s3store.NewS3Store(cfg)

	if err != nil {
		return fmt.Errorf("failed to create s3 store client: %w", err)
	}

	pubsub, err := pubsub.NewPubsub(cfg, ctx, redis)

	if err != nil {
		return fmt.Errorf("failed to create pubsub: %w", err)
	}

	sse := sse.NewStreamer(pubsub)

	app, err := api.New(cfg, ctx, job, store, s3, sse)

	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}

	server := http.NewServer(cfg, app.BuildRouter())

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return job.RegisterAndStart(store, s3)
	})

	gHTTP, shutdownHTTP := server.ListenAndServe()
	g.Go(gHTTP.Wait)

	logger.Info().Msgf("server started on port %d", cfg.Server.Port)

	<-gCtx.Done()

	stop()
	logger.Info().Msg("shutting down gracefully (press Ctrl+C again to force)")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	if shutdownErr := shutdownHTTP(shutdownCtx); shutdownErr != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", shutdownErr)
	}

	job.Executor.Stop()

	logger.Info().Msg("waiting for all goroutines to finish")
	err = g.Wait()

	return err

}
