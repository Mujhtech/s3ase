package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mujhtech/s3ase/auth"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/response"
)

type apiKeyContextKey struct{}

var apiKeyKey apiKeyContextKey

func RequiredAPIKeyAuth(strs *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer ") {
				_ = response.Unauthorized(w, r, nil)
				return
			}
			raw := strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
			if raw == "" || !strings.HasPrefix(raw, "s3ase_") {
				_ = response.Unauthorized(w, r, nil)
				return
			}
			digest := sha256.Sum256([]byte(raw))
			key, err := strs.ApiKeyRepo.FindApiKeyByHash(r.Context(), hex.EncodeToString(digest[:]))
			if err != nil || (key.ExpiredAt.Valid && key.ExpiredAt.Int64 <= time.Now().Unix()) {
				_ = response.Unauthorized(w, r, nil)
				return
			}
			app, err := strs.AppRepo.FindAppByID(r.Context(), key.AppID)
			if err != nil {
				_ = response.Unauthorized(w, r, nil)
				return
			}
			user, err := strs.UserRepo.FindUserByID(r.Context(), key.CreatedBy)
			if err != nil {
				_ = response.Unauthorized(w, r, nil)
				return
			}
			if err := strs.ApiKeyRepo.TouchApiKey(r.Context(), key.ID); err != nil {
				_ = response.Error(w, r, err)
				return
			}
			ctx := context.WithValue(r.Context(), apiKeyKey, key)
			ctx = context.WithValue(ctx, appKey, app)
			ctx = withAuthSession(ctx, &auth.UserSession{User: user})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAPIKeyAccess(write bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, ok := r.Context().Value(apiKeyKey).(*models.ApiKey)
			allowed := ok && key != nil && (key.Access == models.ApiKeyAccessFull ||
				(write && key.Access == models.ApiKeyAccessWrite) || (!write && key.Access == models.ApiKeyAccessRead))
			if !allowed {
				_ = response.Forbidden(w, r, fmt.Errorf("API key access does not permit this operation"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
