package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryHandler struct {
	categoryStore *CategoryStore
}

func NewCategoryHandler(store *CategoryStore) *CategoryHandler {
	return &CategoryHandler{categoryStore: store}
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	categories, err := h.categoryStore.GetAll(r.Context(), userID)
	if err != nil {
		log.Printf("GetAllCategories: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		log.Printf("GetAllCategories: encode failed: %v", err)
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	category := &Category{
		UserID: userID,
		Name:   req.Name,
	}

	created, err := h.categoryStore.Create(r.Context(), category)
	if err != nil {
		log.Printf("CreateCategory: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateCategory: encode failed: %v", err)
	}
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fields := bson.M{"name": req.Name}

	updated, err := h.categoryStore.Update(r.Context(), id, fields)
	if err != nil {
		log.Printf("UpdateCategory: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateCategory: id=%s encode failed: %v", id, err)
	}
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.categoryStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteCategory: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
