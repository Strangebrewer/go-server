package transaction

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionHandler struct {
	transactionStore *TransactionStore
}

func NewTransactionHandler(store *TransactionStore) *TransactionHandler {
	return &TransactionHandler{transactionStore: store}
}

func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	filter := TransactionFilter{
		Month:  r.URL.Query().Get("month"),
		Owner:  r.URL.Query().Get("owner"),
		Shared: r.URL.Query().Get("shared"),
	}

	transactions, err := h.transactionStore.GetAll(r.Context(), userID, filter)
	if err != nil {
		log.Printf("GetAllTransactions: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("GetAllTransactions: encode failed: %v", err)
	}
}

func (h *TransactionHandler) GetOneTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	txn, err := h.transactionStore.GetOne(r.Context(), id)
	if err != nil {
		log.Printf("GetOneTransaction: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(txn); err != nil {
		log.Printf("GetOneTransaction: id=%s encode failed: %v", id, err)
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	txn := &Transaction{
		UserID:        userID,
		Amount:        req.Amount,
		BillID:        req.BillID,
		CategoryID:    req.CategoryID,
		Date:          req.Date,
		Description:   req.Description,
		DestinationID: req.DestinationID,
		Income:        req.Income,
		Owner:         req.Owner,
		Shared:        req.Shared,
		SourceID:      req.SourceID,
		Type:          req.Type,
	}

	created, err := h.transactionStore.Create(r.Context(), txn)
	if err != nil {
		log.Printf("CreateTransaction: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateTransaction: encode failed: %v", err)
	}
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fields := bson.M{
		"amount":         req.Amount,
		"bill_id":        req.BillID,
		"category_id":    req.CategoryID,
		"date":           req.Date,
		"description":    req.Description,
		"destination_id": req.DestinationID,
		"income":         req.Income,
		"owner":          req.Owner,
		"shared":         req.Shared,
		"source_id":      req.SourceID,
		"type":           req.Type,
	}

	updated, err := h.transactionStore.Update(r.Context(), id, fields)
	if err != nil {
		log.Printf("UpdateTransaction: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateTransaction: id=%s encode failed: %v", id, err)
	}
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.transactionStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteTransaction: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
