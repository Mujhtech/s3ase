package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mujhtech/s3ase/config"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	cfg     *config.Config
	handler http.Handler
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		cfg:     cfg,
		handler: handler,
	}
}

func (s *Server) ListenAndServe() (*errgroup.Group, func(ctx context.Context) error) {
	if s.cfg.Server.SSL {
		return s.listenAndServeTLS()
	}

	return s.listenAndServe()
}

func (s *Server) listenAndServe() (*errgroup.Group, func(ctx context.Context) error) {

	var g errgroup.Group
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.Server.Port),
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           s.handler,
	}
	g.Go(func() error {
		return server.ListenAndServe()
	})

	return &g, server.Shutdown
}

func (s *Server) listenAndServeTLS() (*errgroup.Group, func(ctx context.Context) error) {
	var g errgroup.Group
	server1 := &http.Server{
		Addr:              ":http",
		ReadHeaderTimeout: 2 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			target := "https://" + req.Host + "/" + strings.TrimPrefix(req.URL.Path, "/")
			http.Redirect(w, req, target, http.StatusTemporaryRedirect)
		}),
	}
	server2 := &http.Server{
		Addr:              ":https",
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           s.handler,
	}
	g.Go(func() error {
		return server1.ListenAndServe()
	})
	g.Go(func() error {
		return server2.ListenAndServeTLS(
			s.cfg.Server.SSLCertFile,
			s.cfg.Server.SSLKeyFile,
		)
	})

	return &g, func(ctx context.Context) error {
		var sg errgroup.Group
		sg.Go(func() error {
			return server1.Shutdown(ctx)
		})
		sg.Go(func() error {
			return server2.Shutdown(ctx)
		})
		return sg.Wait()
	}
}
