package job

import (
	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
)

func JobRoutes(jobStore *JobStore, tokenService *token.TokenService) chi.Router {
	r := chi.NewRouter()
	h := NewJobHandler(jobStore)

	r.Get("/", h.GetAllJobs)
	r.Get("/{id}", h.GetOneJob)
	r.Post("/", h.CreateJob)
	r.Put("/{id}", h.UpdateJob)
	r.Delete("/{id}", h.DeleteJob)

	return r
}
