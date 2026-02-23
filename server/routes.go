package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/auth"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/users"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	// Ensure TTL index exists (do once at startup)
	_ = application.AuthStore.EnsureIndexes(context.Background())
	_ = application.UserStore.EnsureIndexes(context.Background())

	cookeCfg := auth.CookieConfig{
		Secure:   false,                 // local http
		SameSite: http.SameSiteNoneMode, // fine for Postman + local dev
		TTL:      7 * 24 * time.Hour,
	}

	// IMPORTANT: pepper should come from env/config in real usage.
	// For now, set it to something non-empty.
	pepper := "dev-pepper-change-me"

	// Stub verifyUser for now (replace with real user lookup later)
	verifyUser := func(r *http.Request, email, password string) (string, bool, error) {
		return application.UserStore.VerifyCredentials(r.Context(), email, password)
	}

	authHandler := auth.NewAuthHandler(application.AuthStore, pepper, cookeCfg, verifyUser)
	authMiddleware := auth.NewAuthMiddleware(application.AuthStore, pepper, cookeCfg, 24*time.Hour)

	authDeps := auth.Deps{
		Handler:    authHandler,
		Middleware: authMiddleware,
	}

	r.Mount("/auth", auth.Routes(authDeps))
	r.Mount("/things", thing.Routes(application.ThingStore))
	r.Mount("/users", users.Routes(application.UserStore))
}
