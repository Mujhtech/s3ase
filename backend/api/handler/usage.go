package handler

import (
	"net/http"

	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/services"
)

func (h *Handler) GetUsage(w http.ResponseWriter, r *http.Request) {
	_, app, err := middleware.AuthSessionAndAppFrom(r.Context())
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}
	usage, err := (&services.FindUsageService{App: app, FileRepo: h.store.FileRepo}).Run(r.Context())
	if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	_ = response.Ok(w, r, "usage retrieved", usage)
}
