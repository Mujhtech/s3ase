package job

import (
	"github.com/hibiken/asynq"
	"github.com/mujhtech/s3ase/config"
)

type Scheduler struct {
	scheduler *asynq.Scheduler
}

func NewScheduler(cfg *config.Config, opts asynq.RedisConnOpt) *Scheduler {
	scheduler := asynq.NewScheduler(opts, &asynq.SchedulerOpts{})

	return &Scheduler{
		scheduler: scheduler,
	}
}

func (s *Scheduler) Start() error {
	return s.scheduler.Start()
}

func (s *Scheduler) Stop() {
	s.scheduler.Shutdown()
}
