package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"unicode/utf8"
)

type Handler struct {
	userStore *UserStore
}

func NewHandler(store *UserStore) *Handler {
	return &Handler{userStore: store}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Production-grade validation (keep it simple but real).
	req.Email = normalizeEmail(req.Email)
	if req.Email == "" || len(req.Email) > 254 {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(req.Password) < 12 {
		http.Error(w, "password must be at least 12 characters", http.StatusBadRequest)
		return
	}

	u, err := h.userStore.Create(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			http.Error(w, "email already in use", http.StatusConflict)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(u.Public())
}
