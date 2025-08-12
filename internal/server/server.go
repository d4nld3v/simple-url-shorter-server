package server

import (
	"net/http"

	"github.com/d4nld3v/url-shortener-go/internal/config"
	"github.com/d4nld3v/url-shortener-go/internal/handler"
)

type Server struct {
	cfg config.Config
}

func New(cfg config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {

	mux := http.NewServeMux()

	handler.RegisterUrlRoutes(mux)

	srv := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
