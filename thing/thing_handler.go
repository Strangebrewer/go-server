package thing

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetAllThings(w http.ResponseWriter, r *http.Request) {
	things, err := h.store.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(things); err != nil {
		log.Printf("GetAllThings: encode failed: %v", err)
	}
}

func (h *Handler) GetOneThing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	thing, err := h.store.GetOne(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if thing == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(thing); err != nil {
		log.Printf("GetOneThing: id=%s encode failed: %v", id, err)
	}
}

func (h *Handler) CreateThing(w http.ResponseWriter, r *http.Request) {
	var thing Thing
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	created, err := h.store.Create(r.Context(), &thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateThing: encode failed: %v", err)
	}
}

func (h *Handler) UpdateThing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var thing Thing
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	updated, err := h.store.Update(r.Context(), id, thing)
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
		log.Printf("UpdateThing: id=%s encode failed: %v", id, err)
	}
}

func (h *Handler) DeleteThing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.store.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteThing: id=%s delete failed: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
