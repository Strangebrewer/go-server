package server

import (
	"context"

	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/job"
	"github.com/Strangebrewer/go-server/recruiter"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/user"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	// Ensure TTL index exists (do once at startup)
	_ = application.TokenStore.EnsureIndexes(context.Background())
	_ = application.UserStore.EnsureIndexes(context.Background())

	r.With(token.RequireAccess(application.TokenService)).
		Mount("/things", thing.ThingRoutes(application.ThingStore))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/jobs", job.JobRoutes(application.JobStore, application.TokenService))
	r.With(token.RequireAccess(application.TokenService)).
		Mount("/recruiters", recruiter.RecruiterRoutes(application.RecruiterStore))
	r.Mount("/tokens", token.TokenRoutes(application.TokenService))
	r.Mount("/users", user.UserRoutes(application.UserStore, application.TokenService))
}
