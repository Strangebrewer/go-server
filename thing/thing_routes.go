package thing

import (
	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
)

func ThingRoutes(store *ThingStore, tokenService *token.TokenService) chi.Router {
	r := chi.NewRouter()
	h := NewThingHandler(store)

	r.Get("/", h.GetAllThings)
	r.Get("/{id}", h.GetOneThing)
	r.Group(func(r chi.Router) {
		r.Use(token.RequireAccess(tokenService))
		r.Post("/", h.CreateThing)
		r.Put("/{id}", h.UpdateThing)
		r.Delete("/{id}", h.DeleteThing)
	})

	return r
}
