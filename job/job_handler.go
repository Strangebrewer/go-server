package job

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	jobStore *JobStore
}

func NewJobHandler(store *JobStore) *JobHandler {
	return &JobHandler{jobStore: store}
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobStore.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		log.Printf("GetAllJobs: encode failed: %v", err)
	}
}

func (h *JobHandler) GetOneJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, err := h.jobStore.GetOne(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if job == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		log.Printf("GetOneJob: id=%s encode failed: %v", id, err)
	}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job Job
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userId, _ := token.UserIDFromContext(r.Context())
	job.User = userId

	created, err := h.jobStore.Create(r.Context(), &job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateJob: encode failed: %v", err)
	}
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var job Job
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updated, err := h.jobStore.Update(r.Context(), id, job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updated == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateJob: id=%s encode failed: %v", id, err)
	}
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.jobStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteJob: id=%s delete failed: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
