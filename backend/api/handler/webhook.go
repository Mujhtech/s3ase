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
	WebhookParamId = "webhook_param_id"
)

func getWebhookIdFromPath(r *http.Request) (string, error) {
	rawRef, err := pathParamOrError(r, WebhookParamId)
	if err != nil {
		return "", err
	}

	return url.PathUnescape(rawRef)
}

func (h *Handler) GetWebhooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	findWebhooksService := services.FindWebhooksService{
		App:         app,
		WebhookRepo: h.store.WebhookRepo,
		User:        session.User,
	}

	webhooks, err := findWebhooksService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "webhooks retrieved", webhooks)
}

func (h *Handler) CreateWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	dst := new(dto.CreateWebhookRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	createWebhookService := services.CreateWebhookService{
		App:           app,
		AppMemberRepo: h.store.AppMemberRepo,
		WebhookRepo:   h.store.WebhookRepo,
		User:          session.User,
		Body:          dst,
	}

	webhook, err := createWebhookService.Run(ctx)

	if err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Created(w, r, "webhook created", webhook)
}

func (h *Handler) UpdateWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	webhookId, err := getWebhookIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	dst := new(dto.CreateWebhookRequestDto)

	if err := request.ReadBody(r, dst); err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	updateWebhookService := services.UpdateWebhookService{
		App:           app,
		WebhookId:     webhookId,
		AppMemberRepo: h.store.AppMemberRepo,
		WebhookRepo:   h.store.WebhookRepo,
		User:          session.User,
		Body:          dst,
	}

	if err = updateWebhookService.Run(ctx); err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "webhook updated", nil)
}

func (h *Handler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, app, err := middleware.AuthSessionAndAppFrom(ctx)

	if err != nil {
		_ = response.Unauthorized(w, r, err)
		return
	}

	webhookId, err := getWebhookIdFromPath(r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	deleteWebhookService := services.DeleteWebhookService{
		App:           app,
		WebhookId:     webhookId,
		AppMemberRepo: h.store.AppMemberRepo,
		WebhookRepo:   h.store.WebhookRepo,
		User:          session.User,
	}

	if err = deleteWebhookService.Run(ctx); err != nil {
		_ = response.Error(w, r, err)
		return
	}

	_ = response.Ok(w, r, "webhook deleted", nil)
}
