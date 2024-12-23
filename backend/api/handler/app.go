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

const (
	AppParamId = "app_param_id"
)

func getAppIdFromPath(r *http.Request) (string, error) {
	rawRef, err := pathParamOrError(r, AppParamId)
	if err != nil {
		return "", err
	}

	return url.PathUnescape(rawRef)
}

func (h *Handler) GetApps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, ok := middleware.GetAuthSession(ctx)

	if !ok {
		_ = response.Unauthorized(w, r, nil)
		return
	}

	findAppsService := services.FindAppsService{
		AppRepo: h.store.AppRepo,
		User:    session.User,
	}

	apps, err := findAppsService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Ok(w, r, "apps retrieved", apps)
}

func (h *Handler) CreateApp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, ok := middleware.GetAuthSession(ctx)

	if !ok {
		_ = response.Unauthorized(w, r, nil)
		return
	}

	dst := new(dto.CreateAppRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	createAppService := services.CreateAppService{
		Body:          dst,
		DefaultRegion: h.cfg.Aws.DefaultRegion,
		AppRepo:       h.store.AppRepo,
		AppMemberRepo: h.store.AppMemberRepo,
		ApiKeyRepo:    h.store.ApiKeyRepo,
		S3:            h.s3,
		User:          session.User,
	}

	app, err := createAppService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Ok(w, r, "app created successfully", app)
}

func (h *Handler) UpdateApp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, ok := middleware.GetAuthSession(ctx)

	if !ok {
		_ = response.Unauthorized(w, r, nil)
		return
	}

	appId, err := getAppIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	dst := new(dto.CreateAppRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	updateAppService := services.UpdateAppService{
		AppId:   appId,
		AppRepo: h.store.AppRepo,
		User:    session.User,
		Body:    dst,
	}

	if err = updateAppService.Run(ctx); err != nil {

		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Ok(w, r, "app updated successfully", nil)
}
