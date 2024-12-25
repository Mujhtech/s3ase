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

		_, err := aesCfb.Decrypt(string(t.Payload()))

		if err != nil {
			return err
		}

		return nil
	}
}
