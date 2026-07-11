package middleware

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/mujhtech/s3ase/errors"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/rs/zerolog"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
)

type appMemberContextKey struct{}

var appMemberKey appMemberContextKey

func RequiredAppMember(strs *store.Store) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			session, app, err := AuthSessionAndAppFrom(ctx)

			if err != nil {
				_ = response.Unauthorized(w, r, err)
				return
			}

			member, err := strs.AppMemberRepo.FindAppMemberByAppIDAndUserId(ctx, app.ID, session.User.ID)
			if errors.Is(err, store.ErrNotFound) || member == nil {
				_ = response.Unauthorized(w, r, errs.ErrNotAuthorized)
				return
			}
			if err != nil {
				_ = response.Unauthorized(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, appMemberKey, member)

			log := zerolog.Ctx(ctx)
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("member_id", member.ID)
			})

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func MemberFrom(ctx context.Context) (*models.AppMember, bool) {
	v, ok := ctx.Value(appMemberKey).(*models.AppMember)
	return v, ok && v != nil
}

func MemberFromRequired(ctx context.Context) (*models.AppMember, error) {
	v, ok := MemberFrom(ctx)
	if !ok || v == nil {
		return nil, errs.ErrNotAuthorized
	}
	return v, nil
}
