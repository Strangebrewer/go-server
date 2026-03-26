package category

import "github.com/go-chi/chi/v5"

func CategoryRoutes(store *CategoryStore) chi.Router {
	r := chi.NewRouter()
	h := NewCategoryHandler(store)

	r.Get("/", h.GetAllCategories)
	r.Post("/", h.CreateCategory)
	r.Put("/{id}", h.UpdateCategory)
	r.Delete("/{id}", h.DeleteCategory)

	return r
}
