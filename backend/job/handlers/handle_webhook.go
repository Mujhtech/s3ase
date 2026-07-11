package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/encrypt"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
	"github.com/rs/zerolog"
)

func HandleWebhook(aesCfb encrypt.Encrypt, store *store.Store, s3 *s3store.S3Store) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		payload, err := aesCfb.Decrypt(string(t.Payload()))
		if err != nil {
			return err
		}
		var data WebhookPayload
		if err := json.Unmarshal([]byte(payload), &data); err != nil {
			return err
		}
		webhooks, err := store.WebhookRepo.FindWebhooksByAppID(ctx, data.AppID)
		if err != nil {
			return err
		}
		for _, webhook := range webhooks {
			if !subscribesTo(webhook.Metadata, data.Event) {
				continue
			}
			if err := sendWebhookNotification(ctx, webhook.URL, data); err != nil {
				return err
			}
		}
		return nil
	}
}

type WebhookPayload struct {
	ID        string          `json:"id"`
	AppID     string          `json:"app_id"`
	Event     string          `json:"event"`
	CreatedAt time.Time       `json:"created_at"`
	Data      json.RawMessage `json:"data"`
}

func subscribesTo(metadata map[string]interface{}, wanted string) bool {
	events, ok := metadata["events"].([]interface{})
	if !ok {
		if strings, ok := metadata["events"].([]string); ok {
			for _, event := range strings {
				if event == wanted {
					return true
				}
			}
		}
		return false
	}
	for _, event := range events {
		if value, ok := event.(string); ok && value == wanted {
			return true
		}
	}
	return false
}

func sendWebhookNotification(ctx context.Context, endpoint string, data WebhookPayload) error {
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme != "https" {
		return fmt.Errorf("invalid webhook URL")
	}
	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, parsed.Hostname())
	if err != nil {
		return err
	}
	for _, address := range addresses {
		ip := address.IP
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("webhook target resolves to a private address")
		}
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "S3ase-Webhooks/1.0")
	response, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Warn().Err(closeErr).Msg("failed to close webhook response body")
		}
	}()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("webhook returned HTTP %d", response.StatusCode)
	}
	return nil
}
