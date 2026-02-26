package server

import (
	"context"

	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/users"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	// Ensure TTL index exists (do once at startup)
	_ = application.TokenStore.EnsureIndexes(context.Background())
	_ = application.UserStore.EnsureIndexes(context.Background())

	r.With(token.RequireAccess(application.TokenService)).
		Mount("/things", thing.Routes(application.ThingStore))
	r.Mount("/tokens", token.Routes(application.TokenService))
	r.Mount("/users", users.Routes(application.UserStore, application.TokenService))
}
