package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
	"github.com/mujhtech/s3ase/config"
)

func ApplyCORS(config *config.Config) func(http.Handler) http.Handler {
	return cors.New(
		cors.Options{
			AllowedOrigins:   config.Cors.AllowedOrigins,
			AllowedMethods:   config.Cors.AllowedMethods,
			AllowedHeaders:   config.Cors.AllowedHeaders,
			ExposedHeaders:   config.Cors.ExposedHeaders,
			AllowCredentials: config.Cors.AllowCredentials,
			MaxAge:           config.Cors.MaxAge,
		},
	).Handler
}
