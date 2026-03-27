package template

import "github.com/go-chi/chi/v5"

func TemplateRoutes(store *TemplateStore) chi.Router {
	r := chi.NewRouter()
	h := NewTemplateHandler(store)

	r.Get("/", h.GetAllTemplates)
	r.Post("/", h.CreateTemplate)
	r.Put("/{id}", h.UpdateTemplate)
	r.Delete("/{id}", h.DeleteTemplate)

	return r
}
