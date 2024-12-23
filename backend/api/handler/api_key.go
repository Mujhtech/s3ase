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
	ApiKeyParamId = "api_key_param_id"
)

func getApiKeyIdFromPath(r *http.Request) (string, error) {
	rawRef, err := pathParamOrError(r, WebhookParamId)
	if err != nil {
		return "", err
	}

	return url.PathUnescape(rawRef)
}

func (h *Handler) GetApiKeys(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	findApiKeysService := services.FindApiKeysService{
		App:        app,
		AppRepo:    h.store.AppRepo,
		ApiKeyRepo: h.store.ApiKeyRepo,
		User:       session.User,
	}

	apiKeys, err := findApiKeysService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Ok(w, r, "api keys retrieved", apiKeys)
}

func (h *Handler) CreateApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	dst := new(dto.CreateApiKeyRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	createApiKeyService := services.CreateApiKeyService{
		App:           app,
		AppMemberRepo: h.store.AppMemberRepo,
		ApiKeyRepo:    h.store.ApiKeyRepo,
		User:          session.User,
		Body:          dst,
	}

	apiKey, err := createApiKeyService.Run(ctx)

	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Created(w, r, "api key created", apiKey)
}

func (h *Handler) UpdateApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	apiKeyId, err := getApiKeyIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	dst := new(dto.CreateApiKeyRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	updateApiKeyService := services.UpdateApiKeyService{
		ApiKeyId:   apiKeyId,
		App:        app,
		ApiKeyRepo: h.store.ApiKeyRepo,
		User:       session.User,
		Body:       dst,
	}

	if err = updateApiKeyService.Run(ctx); err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Created(w, r, "api key updated", nil)
}

func (h *Handler) DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	apiKeyId, err := getApiKeyIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	deleteApiKeyService := services.DeleteApiKeyService{
		ApiKeyId:   apiKeyId,
		App:        app,
		ApiKeyRepo: h.store.ApiKeyRepo,
		User:       session.User,
	}

	if err = deleteApiKeyService.Run(ctx); err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	_ = response.Created(w, r, "api key deleted", nil)
}
