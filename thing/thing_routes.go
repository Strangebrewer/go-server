package thing

import "github.com/go-chi/chi/v5"

func Routes(store *ThingStore) chi.Router {
	r := chi.NewRouter()
	h := NewThingHandler(store)

	r.Get("/", h.GetAllThings)
	r.Get("/{id}", h.GetOneThing)
	r.Post("/", h.CreateThing)
	r.Put("/{id}", h.UpdateThing)
	r.Delete("/{id}", h.DeleteThing)

	return r
}
