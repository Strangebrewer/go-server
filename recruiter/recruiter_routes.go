package recruiter

import "github.com/go-chi/chi/v5"

func RecruiterRoutes(store *RecruiterStore) chi.Router {
	r := chi.NewRouter()
	h := NewRecruiterHandler(store)

	r.Get("/", h.GetAllRecruiters)
	r.Get("/{id}", h.GetOneRecruiter)
	r.Post("/", h.CreateRecruiter)
	r.Put("/{id}", h.UpdateRecruiter)
	r.Delete("/{id}", h.DeleteRecruiter)

	return r
}
