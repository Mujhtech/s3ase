package services

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

type CreateWebhookService struct {
	App           *models.App
	AppMemberRepo store.AppMemberRepository
	WebhookRepo   store.WebhookRepository
	User          *models.User
	Body          *dto.CreateWebhookRequestDto
}

func (c *CreateWebhookService) Run(ctx context.Context) (*models.Webhook, error) {
	if err := requireAppOwner(ctx, c.AppMemberRepo, c.App.ID, c.User.ID); err != nil {
		return nil, err
	}

	id := uuid.New().String()
	if err := validateWebhook(c.Body); err != nil {
		return nil, err
	}

	if err := c.WebhookRepo.CreateWebhook(ctx, &models.Webhook{
		ID:          id,
		AppID:       c.App.ID,
		Name:        c.Body.Name,
		Description: null.NewString(c.Body.Description, c.Body.Description != ""),
		URL:         c.Body.Url,
		CreatedBy:   c.User.ID,
		Metadata: map[string]interface{}{
			"events": c.Body.Events,
		},
	}); err != nil {
		return nil, err
	}

	return c.WebhookRepo.FindWebhookByID(ctx, id)
}

var validWebhookEvents = map[string]struct{}{
	"upload.started": {}, "upload.completed": {}, "upload.failed": {}, "upload.deleted": {},
	"file.created": {}, "file.updated": {}, "file.deleted": {},
}

func validateWebhook(body *dto.CreateWebhookRequestDto) error {
	if body == nil || strings.TrimSpace(body.Name) == "" {
		return fmt.Errorf("%w: webhook name is required", errs.ErrInvalidInput)
	}
	parsed, err := url.ParseRequestURI(strings.TrimSpace(body.Url))
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil {
		return fmt.Errorf("%w: webhook URL must be a valid HTTPS URL", errs.ErrInvalidInput)
	}
	if ip := net.ParseIP(parsed.Hostname()); ip != nil && (ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast()) {
		return fmt.Errorf("%w: webhook URL cannot target a private address", errs.ErrInvalidInput)
	}
	if len(body.Events) == 0 {
		return fmt.Errorf("%w: select at least one webhook event", errs.ErrInvalidInput)
	}
	for _, event := range body.Events {
		if _, ok := validWebhookEvents[event]; !ok {
			return fmt.Errorf("%w: unsupported webhook event %q", errs.ErrInvalidInput, event)
		}
	}
	return nil
}
