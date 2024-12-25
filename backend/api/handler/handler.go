package handler

import (
	"context"

	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/job"
)

type Handler struct {
	cfg   *config.Config
	ctx   context.Context
	store *store.Store
	job   *job.Job
	s3    *s3store.S3Store
	sse   sse.Streamer
}

func New(
	cfg *config.Config,
	ctx context.Context,
	job *job.Job,
	store *store.Store,
	s3 *s3store.S3Store,
	sse sse.Streamer,
) (*Handler, error) {

	return &Handler{
		cfg:   cfg,
		ctx:   ctx,
		store: store,
		job:   job,
		s3:    s3,
		sse:   sse,
	}, nil
}
