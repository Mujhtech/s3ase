package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

func InstrumentRequests(serverName string, r chi.Router) func(next http.Handler) http.Handler {
	return otelchi.Middleware(serverName, otelchi.WithChiRoutes(r))
}
