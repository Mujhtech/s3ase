package handler

import (
	"net/http"

	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/rs/zerolog/log"
)

func (h *Handler) Event(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	chEvents, chErr, sseCancel := h.sse.Subscribe(ctx, app.ID)

	defer func() {
		err := sseCancel(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to cancel sse subscription")
		}
	}()

	response.Stream(ctx, w, h.ctx.Done(), chEvents, chErr)
}
