package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type AuthHandler struct {
	authStore *AuthStore
	pepper    string
	cookieCfg CookieConfig
	// You provide this (or call your users service/store):
	verifyUser func(r *http.Request, email, password string) (userID string, ok bool, err error)
}

func NewAuthHandler(authStore *AuthStore, pepper string, cookieCfg CookieConfig,
	verifyUser func(r *http.Request, email, password string) (string, bool, error),
) *AuthHandler {
	return &AuthHandler{
		authStore: authStore, pepper: pepper, cookieCfg: cookieCfg, verifyUser: verifyUser,
	}
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	userID, ok, err := h.verifyUser(r, req.Email, req.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	rawSID, err := NewToken(32)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	now := time.Now().UTC()
	session := Session{
		ID:        HashToken(rawSID, h.pepper),
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(h.cookieCfg.TTL),
	}

	if err := h.authStore.Create(r.Context(), session); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	SetSessionCookie(w, rawSID, h.cookieCfg)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SidCookieName)
	if err == nil && c.Value != "" {
		_ = h.authStore.Delete(r.Context(), HashToken(c.Value, h.pepper))
	}
	ClearSessionCookie(w, h.cookieCfg)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"userId": uid})
}
