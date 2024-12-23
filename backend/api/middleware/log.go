package middleware

import (
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func SetupRequestLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log := log.Logger.With().Logger()
			ctx := log.WithContext(r.Context())
			ctx = logr.NewContext(ctx, zerologr.New(&log))
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.
					Str("http.original_url", r.URL.String())
			})

			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("http.handler", "api")
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}

func HLogAccessLogHandler() func(http.Handler) http.Handler {
	return hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Int("http.status_code", status).
				Int("http.response_size_bytes", size).
				Dur("http.elapsed_ms", duration).
				Msg("http request completed.")
		},
	)
}
