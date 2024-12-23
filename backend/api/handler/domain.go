package handler

import (
	"net/http"

	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/internal/pkg/request"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/services"
)

func (h *Handler) GetDomain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	findDomainService := services.FindDomainService{
		App:        app,
		DomainRepo: h.store.DomainRepo,
		User:       session.User,
	}

	domain, err := findDomainService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Ok(w, r, "domain retrieved", domain)
}

func (h *Handler) CreateOrUpdateDomain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	dst := new(dto.CreateOrUpdateDomainRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	createOrUpdateDomainService := services.CreateOrUpdateDomainService{
		App:        app,
		DomainRepo: h.store.DomainRepo,
		User:       session.User,
		Body:       dst,
	}

	domain, err := createOrUpdateDomainService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Created(w, r, "domain created", domain)
}
