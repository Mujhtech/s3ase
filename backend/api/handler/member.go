package handler

import (
	"net/http"

	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/services"
)

func (h *Handler) GetMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	findMembersService := services.FindMembersService{
		App:           app,
		AppRepo:       h.store.AppRepo,
		AppMemberRepo: h.store.AppMemberRepo,
		UserRepo:      h.store.UserRepo,
		User:          session.User,
	}

	members, err := findMembersService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Ok(w, r, "members retrieved", members)
}
