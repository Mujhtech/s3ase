package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/response"
)

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	_, app, err := middleware.AuthSessionAndAppFrom(r.Context())
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}
	subscription, err := h.store.AppSubscriptionRepo.FindAppSubscriptionByAppID(r.Context(), app.ID)
	if errors.Is(err, store.ErrNotFound) {
		subscription = &models.AppSubscription{ID: uuid.NewString(), AppID: app.ID, Plan: "basic", Status: "active", Metadata: models.Metadata{}}
		if err := h.store.AppSubscriptionRepo.CreateAppSubscription(r.Context(), subscription); err != nil {
			_ = response.Error(w, r, err)
			return
		}
	} else if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	_ = response.Ok(w, r, "subscription retrieved", subscription)
}
