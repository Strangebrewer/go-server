package account

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountHandler struct {
	accountStore *AccountStore
}

func NewAccountHandler(store *AccountStore) *AccountHandler {
	return &AccountHandler{accountStore: store}
}

func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	accounts, err := h.accountStore.GetAll(r.Context(), userID)
	if err != nil {
		log.Printf("GetAllAccounts: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		log.Printf("GetAllAccounts: encode failed: %v", err)
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	account := &Account{
		UserID:      userID,
		Balance:     req.Balance,
		Description: req.Description,
		Name:        req.Name,
		Owner:       req.Owner,
		Type:        req.Type,
		Status:      "active",
	}

	created, err := h.accountStore.Create(r.Context(), account)
	if err != nil {
		log.Printf("CreateAccount: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateAccount: encode failed: %v", err)
	}
}

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fields := bson.M{
		"balance":     req.Balance,
		"description": req.Description,
		"name":        req.Name,
		"owner":       req.Owner,
		"status":      req.Status,
		"type":        req.Type,
	}

	updated, err := h.accountStore.Update(r.Context(), id, fields)
	if err != nil {
		log.Printf("UpdateAccount: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateAccount: id=%s encode failed: %v", id, err)
	}
}

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.accountStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteAccount: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
