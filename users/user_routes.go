package users

import "github.com/go-chi/chi/v5"

func Routes(userStore *UserStore) chi.Router {
	r := chi.NewRouter()

	h := NewHandler(userStore)
	r.Post("/", h.CreateUser)

	return r
}
