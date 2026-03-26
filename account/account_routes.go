package account

import "github.com/go-chi/chi/v5"

func AccountRoutes(store *AccountStore) chi.Router {
	r := chi.NewRouter()
	h := NewAccountHandler(store)

	r.Get("/", h.GetAllAccounts)
	r.Post("/", h.CreateAccount)
	r.Put("/{id}", h.UpdateAccount)
	r.Delete("/{id}", h.DeleteAccount)

	return r
}
