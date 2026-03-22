package job

import (
	"github.com/Strangebrewer/go-server/recruiter"
	"github.com/go-chi/chi/v5"
)

func JobRoutes(jobStore *JobStore, recruiterStore *recruiter.RecruiterStore) chi.Router {
	r := chi.NewRouter()
	h := NewJobHandler(jobStore, recruiterStore)

	r.Get("/", h.GetAllJobs)
	r.Get("/{id}", h.GetOneJob)
	r.Post("/", h.CreateJob)
	r.Put("/{id}", h.UpdateJob)
	r.Delete("/{id}", h.DeleteJob)

	return r
}
