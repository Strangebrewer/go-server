package template

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type TemplateHandler struct {
	templateStore *TemplateStore
}

func NewTemplateHandler(store *TemplateStore) *TemplateHandler {
	return &TemplateHandler{templateStore: store}
}

func (h *TemplateHandler) GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	templates, err := h.templateStore.GetAll(r.Context(), userID)
	if err != nil {
		log.Printf("GetAllTemplates: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(templates); err != nil {
		log.Printf("GetAllTemplates: encode failed: %v", err)
	}
}

func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	t := &Template{
		UserID:        userID,
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Description:   req.Description,
		Owner:         req.Owner,
		Type:          req.Type,
		Shared:        req.Shared,
		SourceID:      req.SourceID,
		DestinationID: req.DestinationID,
	}

	created, err := h.templateStore.Create(r.Context(), t)
	if err != nil {
		log.Printf("CreateTemplate: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateTemplate: encode failed: %v", err)
	}
}

func (h *TemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fields := bson.M{
		"category_id":    req.CategoryID,
		"name":           req.Name,
		"description":    req.Description,
		"owner":          req.Owner,
		"type":           req.Type,
		"shared":         req.Shared,
		"source_id":      req.SourceID,
		"destination_id": req.DestinationID,
	}

	updated, err := h.templateStore.Update(r.Context(), id, fields)
	if err != nil {
		log.Printf("UpdateTemplate: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateTemplate: id=%s encode failed: %v", id, err)
	}
}

func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.templateStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteTemplate: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
