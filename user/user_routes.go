package user

import (
	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(userStore *UserStore, tokenService *token.TokenService) chi.Router {
	r := chi.NewRouter()
	h := NewUserHandler(userStore, tokenService)

	r.Post("/", h.CreateUser)
	r.Post("/login", h.Login)

	r.Group(func(pr chi.Router) {
		pr.Use(token.RequireAccess(tokenService))
		pr.Get("/me", h.Me)
	})

	return r
}
