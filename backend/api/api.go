package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mujhtech/s3ase/api/handler"
	"github.com/mujhtech/s3ase/api/middleware"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/database/store"
	"github.com/mujhtech/s3ase/internal/pkg/s3store"
	"github.com/mujhtech/s3ase/internal/pkg/sse"
	"github.com/mujhtech/s3ase/job"
	"github.com/rs/zerolog/hlog"
)

type Api struct {
	handler *handler.Handler
	cfg     *config.Config
	store   *store.Store
	job     *job.Job
}

func New(
	cfg *config.Config,
	ctx context.Context,
	job *job.Job,
	store *store.Store,
	s3 *s3store.S3Store,
	sse sse.Streamer,
) (*Api, error) {

	h, err := handler.New(cfg, ctx, job, store, s3, sse)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %w", err)
	}

	return &Api{
		handler: h,
		cfg:     cfg,
		store:   store,
		job:     job,
	}, nil

}

func (a *Api) BuildRouter() *chi.Mux {
	router := chi.NewMux()

	router.Use(middleware.SetupRequestLog())
	router.Use(chiMiddleware.NoCache)
	router.Use(chiMiddleware.Recoverer)
	router.Use(hlog.URLHandler("http.url"))
	router.Use(hlog.MethodHandler("http.method"))
	router.Use(middleware.WriteRequestIDHeader())
	router.Use(middleware.HLogAccessLogHandler())

	router.Route("/api", func(r chi.Router) {

		// v1 route
		r.Route("/v1", func(r chi.Router) {})

		// ui route
		r.Route("/ui", func(r chi.Router) {

			r.Use(
				chiMiddleware.Maybe(middleware.RequiredUserAuth(a.cfg, a.store), shouldAllowAuth),
				middleware.AppIdRequestHeader(a.store),
				chiMiddleware.Maybe(middleware.RequiredAppMember(a.cfg, a.store), shouldAllowMember),
			)
			// r.Use(middleware.AppIdRequestHeader(a.store))
			// r.Use(chiMiddleware.Maybe(middleware.RequiredAppMember(a.cfg, a.store), shouldAllowMember))

			// features
			r.Get("/features", a.handler.GetFeatures)

			// auth route
			r.Route("/auth", func(r chi.Router) {
				r.Get(fmt.Sprintf("/{%s}", handler.AuthProviderKey), a.handler.Authenticate)
				r.Get(fmt.Sprintf("/{%s}/callback", handler.AuthProviderKey), a.handler.AuthenticateCallback)
				r.Post(fmt.Sprintf("/{%s}/callback", handler.AuthProviderKey), a.handler.AuthenticateCallbackPost)
			})

			// user route
			r.Route("/user", func(r chi.Router) {
				r.Get("/", a.handler.GetUser)
			})

			// apps route
			r.Route("/apps", func(r chi.Router) {
				r.Get("/", a.handler.GetApps)
				r.Post("/", a.handler.CreateApp)
				r.Put(fmt.Sprintf("/{%s}", handler.AppParamId), a.handler.UpdateApp)
			})

			// api keys
			r.Route("/api_keys", func(r chi.Router) {
				r.Get("/", a.handler.GetApiKeys)
				r.Post("/", a.handler.CreateApiKey)
				r.Put(fmt.Sprintf("/{%s}", handler.ApiKeyParamId), a.handler.UpdateApiKey)
				r.Delete(fmt.Sprintf("/{%s}", handler.ApiKeyParamId), a.handler.DeleteApiKey)
			})

			// api keys
			r.Route("/webhooks", func(r chi.Router) {
				r.Get("/", a.handler.GetWebhooks)
				r.Post("/", a.handler.CreateWebhook)
				r.Put(fmt.Sprintf("/{%s}", handler.WebhookParamId), a.handler.UpdateWebhook)
				r.Delete(fmt.Sprintf("/{%s}", handler.WebhookParamId), a.handler.DeleteWebhook)
			})

			// members
			r.Route("/members", func(r chi.Router) {
				r.Get("/", a.handler.GetMembers)
			})

			// files
			r.Route("/files", func(r chi.Router) {
				r.Get("/", a.handler.GetFiles)
				r.Get(fmt.Sprintf("/{%s}", handler.FileParamID), a.handler.GetFile)
			})

			// domain
			r.Route("/domain", func(r chi.Router) {
				r.Get("/", a.handler.GetDomain)
				r.Post("/", a.handler.CreateOrUpdateDomain)
			})

			// folders
			r.Route("/folders", func(r chi.Router) {
				r.Get("/", a.handler.GetFolders)
				r.Get(fmt.Sprintf("/{%s}", handler.FolderParamID), a.handler.GetFolder)
				r.Post("/", a.handler.CreateFolder)
				r.Put(fmt.Sprintf("/{%s}", handler.FolderParamID), a.handler.UpdateFolder)
				r.Delete(fmt.Sprintf("/{%s}", handler.FolderParamID), a.handler.DeleteFolder)
			})

			// sse event
			r.Get("/sse", a.handler.Event)
		})
	})

	// job route
	router.Route("/job", func(r chi.Router) {
		r.Handle("/monitoring/*", a.job.Client.Monitor())
	})

	return router
}

var guestRoutes = []string{
	"/auth/github",
	"/auth/github/callback",
	"/auth/google",
	"/auth/google/callback",
	"/features",
}

func shouldAllowAuth(r *http.Request) bool {

	for _, route := range guestRoutes {
		if strings.HasSuffix(r.URL.Path, route) {
			return false
		}
	}

	return true
}

func shouldAllowMember(r *http.Request) bool {

	userRoute := []string{
		"/user",
		"/apps",
		"/auth/github",
		"/auth/github/callback",
		"/auth/google",
		"/auth/google/callback",
		"/features",
	}

	for _, route := range userRoute {
		if strings.HasSuffix(r.URL.Path, route) {
			return false
		}
	}

	return true
}
