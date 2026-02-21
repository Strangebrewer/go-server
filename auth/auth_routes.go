package auth

import "github.com/go-chi/chi/v5"

type Deps struct {
	Handler    *AuthHandler
	Middleware *AuthMiddleware
}

func Routes(deps Deps) chi.Router {
	r := chi.NewRouter()

	// Public auth endpoints
	r.Post("/login", deps.Handler.Login)
	r.Post("/logout", deps.Handler.Logout)

	// Protected auth endpoints
	r.Group(func(r chi.Router) {
		r.Use(deps.Middleware.RequireAuth)
		r.Get("/me", deps.Handler.Me)
	})

	return r
}
