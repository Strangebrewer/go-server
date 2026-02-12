package server

import (
	"net/http"

	"github.com/Strangebrewer/go-server/app"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	HTTPServer *http.Server
}

func New(addr string, application *app.Application) *Server {
	r := chi.NewRouter()

	registerRoutes(r, application)

	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{HTTPServer: s}
}
