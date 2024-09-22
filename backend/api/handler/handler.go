package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database"
)

type AppHandler struct {
	cfg *config.Config
	db  *database.Database
}

func New(cfg *config.Config, db *database.Database) (*AppHandler, error) {

	return &AppHandler{
		cfg: cfg,
		db:  db,
	}, nil
}

func (a *AppHandler) BuildHandler() *chi.Mux {
	router := chi.NewMux()

	return router
}
