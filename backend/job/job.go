package job

import (
	"github.com/hibiken/asynq"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/encrypt"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
	"github.com/mujhtech/s3ase/internal/redis"
	"github.com/mujhtech/s3ase/job/handlers"
	rdsv9 "github.com/redis/go-redis/v9"
)

type Job struct {
	Client    *Client
	Executor  *Executor
	Scheduler *Scheduler
	aesCfb    encrypt.Encrypt
}

func NewJob(cfg *config.Config, redis *redis.Redis) (*Job, error) {

	var c asynq.RedisConnOpt
	var _ = redis.MakeRedisClient().(rdsv9.UniversalClient)
	c = redis

	aesCfb, err := encrypt.NewAesCfb(cfg.EncryptionKey)

	if err != nil {
		return nil, err
	}

	return &Job{
		aesCfb:    aesCfb,
		Client:    NewClient(c, aesCfb),
		Executor:  NewExecutor(cfg, c),
		Scheduler: NewScheduler(cfg, c),
	}, nil
}

func (j *Job) RegisterAndStart(store *store.Store, s3 *s3store.S3Store) error {
	j.Executor.RegisterJobHandler(JobNameWebhook, asynq.HandlerFunc(handlers.HandleWebhook(j.aesCfb, store, s3)))

	return j.Executor.Start()
}
