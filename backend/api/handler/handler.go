package handler

import (
	"context"

	"github.com/mujhtech/s3ase/cache"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/internal/protocol"
	"github.com/mujhtech/s3ase/job"
)

type Handler struct {
	cfg      *config.Config
	ctx      context.Context
	cache    cache.Cache
	store    *store.Store
	job      *job.Job
	s3       *s3store.S3Store
	sse      sse.Streamer
	protocol *protocol.Protocol
}

func New(
	cfg *config.Config,
	ctx context.Context,
	cache cache.Cache,
	job *job.Job,
	store *store.Store,
	s3 *s3store.S3Store,
	sse sse.Streamer,
	protocol *protocol.Protocol,
) (*Handler, error) {

	return &Handler{
		cfg:      cfg,
		ctx:      ctx,
		cache:    cache,
		store:    store,
		job:      job,
		s3:       s3,
		sse:      sse,
		protocol: protocol,
	}, nil
}
