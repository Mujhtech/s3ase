package handler

import (
	"net/http"
	"net/url"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/internal/pkg/request"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/services"
)

const MemberParamID = "member_param_id"

func getMemberIDFromPath(r *http.Request) (string, error) {
	raw, err := pathParamOrError(r, MemberParamID)
	if err != nil {
		return "", err
	}
	return url.PathUnescape(raw)
}

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
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "members retrieved", members)
}

func (h *Handler) CreateMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, app, err := middleware.AuthSessionAndAppFrom(ctx)
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}
	dst := new(dto.CreateMemberRequestDto)
	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}
	member, err := (&services.CreateMemberService{
		App: app, User: session.User, Body: dst, UserRepo: h.store.UserRepo,
		AppMemberRepo: h.store.AppMemberRepo,
	}).Run(ctx)
	if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	_ = response.Created(w, r, "member added", member)
}

func (h *Handler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, app, err := middleware.AuthSessionAndAppFrom(ctx)
	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}
	memberID, err := getMemberIDFromPath(r)
	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}
	err = (&services.DeleteMemberService{
		App: app, User: session.User, MemberID: memberID, AppMemberRepo: h.store.AppMemberRepo,
	}).Run(ctx)
	if err != nil {
		_ = response.Error(w, r, err)
		return
	}
	_ = response.Ok(w, r, "member removed", nil)
}
