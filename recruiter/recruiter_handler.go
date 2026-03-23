package recruiter

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type RecruiterHandler struct {
	recruiterStore *RecruiterStore
}

func NewRecruiterHandler(store *RecruiterStore) *RecruiterHandler {
	return &RecruiterHandler{recruiterStore: store}
}

func (h *RecruiterHandler) GetAllRecruiters(w http.ResponseWriter, r *http.Request) {
	recruiters, err := h.recruiterStore.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recruiters); err != nil {
		log.Printf("GetAllRecruiters: encode failed: %v", err)
	}
}

func (h *RecruiterHandler) GetOneRecruiter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	recruiter, err := h.recruiterStore.GetOne(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if recruiter == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recruiter); err != nil {
		log.Printf("GetOneRecruiter: id=%s encode failed: %v", id, err)
	}
}

func (h *RecruiterHandler) CreateRecruiter(w http.ResponseWriter, r *http.Request) {
	var recruiter Recruiter
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&recruiter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userId, _ := token.UserIDFromContext(r.Context())
	recruiter.User = userId

	recruiter.Email = strings.ToLower(recruiter.Email)

	created, err := h.recruiterStore.Create(r.Context(), &recruiter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateRecruiter: encode failed: %v", err)
	}
}

func (h *RecruiterHandler) UpdateRecruiter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var fields bson.M
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	delete(fields, "id")

	updated, err := h.recruiterStore.Update(r.Context(), id, fields)
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
		log.Printf("UpdateRecruiter: id=%s encode failed: %v", id, err)
	}
}

func (h *RecruiterHandler) DeleteRecruiter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.recruiterStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteRecruiter: id=%s delete failed: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
