package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Strangebrewer/go-server/token"
)

type UserHandler struct {
	userStore    *UserStore
	tokenService *token.TokenService
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(store *UserStore, tokenService *token.TokenService) *UserHandler {
	return &UserHandler{userStore: store, tokenService: tokenService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	u, err := h.userStore.FindByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	ok, err := VerifyPassword(req.Password, u.PasswordHash)
	if err != nil || !ok {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	res, err := h.tokenService.IssueForUser(r.Context(), u.ID)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	response := &LoginResponse{
		User: PublicUser{
			Email: u.Email,
			ID:    u.ID.Hex(),
		},
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := token.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	u, err := h.userStore.FindByID(r.Context(), uid.Hex())
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(u.Public())
}
