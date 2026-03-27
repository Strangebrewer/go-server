package thing

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ThingHandler struct {
	thingStore *ThingStore
}

func NewThingHandler(store *ThingStore) *ThingHandler {
	return &ThingHandler{thingStore: store}
}

func (h *ThingHandler) GetAllThings(w http.ResponseWriter, r *http.Request) {
	things, err := h.thingStore.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(things); err != nil {
		log.Printf("GetAllThings: encode failed: %v", err)
	}
}

func (h *ThingHandler) GetOneThing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	thing, err := h.thingStore.GetOne(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("GetOneThing: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(thing); err != nil {
		log.Printf("GetOneThing: id=%s encode failed: %v", id, err)
	}
}

func (h *ThingHandler) CreateThing(w http.ResponseWriter, r *http.Request) {
	var thing Thing
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&thing)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	created, err := h.thingStore.Create(r.Context(), &thing)
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

func (h *ThingHandler) UpdateThing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var fields bson.M
	if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	delete(fields, "id")

	updated, err := h.thingStore.Update(r.Context(), id, fields)
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

func (h *ThingHandler) DeleteThing(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.thingStore.Delete(r.Context(), id)
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
