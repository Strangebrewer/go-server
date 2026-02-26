package server

import (
	"net/http"

	"github.com/Strangebrewer/go-server/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	HTTPServer *http.Server
}

func New(addr string, application *app.Application) *Server {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	registerRoutes(r, application)

	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{HTTPServer: s}
}
