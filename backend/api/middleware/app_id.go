package middleware

import (
	"context"
	"net/http"

	"github.com/mujhtech/s3ase/auth"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/errors"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/rs/zerolog"
)

const (
	appIDHeader     = "x-app-id"
	appKey      key = iota
)

func AppIdRequestHeader(store *store.Store) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// extract app id from request if exist
			appId := r.Header.Get(appIDHeader)

			if appId != "" {

				app, err := store.AppRepo.FindAppByID(ctx, appId)

				if err != nil {
					_ = response.BadRequest(w, r, err)
					return
				}

				ctx = context.WithValue(ctx, appKey, app)

				log := zerolog.Ctx(ctx)
				log.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str("app_id", appId)
				})

				// write request id to response header
				w.Header().Set(appIDHeader, appId)
			}

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AppFrom(ctx context.Context) (*models.App, bool) {
	v, ok := ctx.Value(appKey).(*models.App)
	return v, ok && v != nil
}

func AuthSessionAndAppFrom(ctx context.Context) (*auth.UserSession, *models.App, error) {
	session, ok := GetAuthSession(ctx)

	if !ok {
		return nil, nil, errors.ErrNotAuthorized
	}

	app, ok := AppFrom(ctx)

	if !ok {
		return nil, nil, errors.ErrNotAuthorized
	}

	return session, app, nil
}
