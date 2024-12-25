package handlers

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/encrypt"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
)

func HandleWebhook(aesCfb encrypt.Encrypt, store *store.Store, s3 *s3store.S3Store) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {

		// payload, err := aesCfb.Decrypt(string(t.Payload()))
		// if err != nil {
		// 	return err
		// }

		// var data struct {
		// 	Event string          `json:"event"`
		// 	Data  json.RawMessage `json:"data"`
		// }

		// if err := json.Unmarshal([]byte(payload), &data); err != nil {
		// 	return err
		// }

		// // Get webhooks for the app that match this event
		// webhooks, err := store.WebhookRepo.FindWebhooksByEvent(ctx, data.Event)
		// if err != nil {
		// 	return err
		// }

		// // Send webhook notifications
		// for _, webhook := range webhooks {
		// 	if err := sendWebhookNotification(webhook.URL, data); err != nil {
		// 		// Log error but continue with other webhooks
		// 		log.Error().Err(err).Msg("failed to send webhook notification")
		// 	}
		// }

		return nil
	}
}

// func sendWebhookNotification(url string, data interface{}) error {
// 	// Implement HTTP POST to webhook URL
// 	// Add retry logic, timeout etc
// 	return nil
// }
