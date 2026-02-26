package server

import (
	"context"

	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/users"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	// Ensure TTL index exists (do once at startup)
	_ = application.UserStore.EnsureIndexes(context.Background())

	r.Mount("/things", thing.Routes(application.ThingStore))
	r.Mount("/users", users.Routes(application.UserStore))
}
