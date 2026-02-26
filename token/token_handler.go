package token

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenHandler struct {
	svc *TokenService
}

func NewTokenHandler(svc *TokenService) *TokenHandler {
	return &TokenHandler{svc: svc}
}

type issueReq struct {
	UserID string `json:"userId"`
}

func (h *TokenHandler) Issue(w http.ResponseWriter, r *http.Request) {
	var req issueReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	uid, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	res, err := h.svc.IssueForUser(r.Context(), uid)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (h *TokenHandler) Exchange(w http.ResponseWriter, r *http.Request) {
	refresh, ok := bearerToken(r)
	if !ok {
		http.Error(w, "missing bearer token", http.StatusUnauthorized)
		return
	}

	res, err := h.svc.Exchange(r.Context(), refresh)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (h *TokenHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	refresh, ok := bearerToken(r)
	if !ok {
		http.Error(w, "missing bearer token", http.StatusUnauthorized)
		return
	}

	if err := h.svc.Revoke(r.Context(), refresh); err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
