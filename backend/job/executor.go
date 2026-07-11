package job

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/mujhtech/s3ase/config"
)

type Executor struct {
	mux *asynq.ServeMux
	srv *asynq.Server
}

func NewExecutor(cfg *config.Config, opts asynq.RedisConnOpt) *Executor {

	srv := asynq.NewServer(
		opts,
		asynq.Config{
			Concurrency: cfg.Job.Concurrency,
			BaseContext: context.Background,
		},
	)

	mux := asynq.NewServeMux()

	return &Executor{
		mux: mux,
		srv: srv,
	}
}

func (e *Executor) Start() error {
	return e.srv.Start(e.mux)
}

func (e *Executor) Stop() {
	e.srv.Stop()
	e.srv.Shutdown()
}

func (e *Executor) RegisterJobHandler(name JobName, handler asynq.Handler) {
	e.mux.HandleFunc(string(name), handler.ProcessTask)
}
