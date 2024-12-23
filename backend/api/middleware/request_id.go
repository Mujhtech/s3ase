package middleware

import (
	"context"
	"net/http"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

type key int

const (
	requestIDHeader     = "X-Request-Id"
	requestIDKey    key = iota
)

func WriteRequestIDHeader() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// extract request id from request if exist
			var requestID string
			requestIDs, ok := r.Header[requestIDHeader]

			if ok && len(requestIDs) > 0 && len(requestIDs[0]) > 0 {
				requestID = requestIDs[0]
			} else {
				// generate new request id
				requestID = xid.New().String()
			}

			ctx = context.WithValue(ctx, requestIDKey, requestID)

			log := zerolog.Ctx(ctx)
			log.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("request_id", requestID)
			})

			// write request id to response header
			w.Header().Set(requestIDHeader, requestID)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequestIDFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok && v != ""
}
