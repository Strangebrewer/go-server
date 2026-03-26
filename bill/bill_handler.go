package bill

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/transaction"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillHandler struct {
	billStore        *BillStore
	transactionStore *transaction.TransactionStore
}

func NewBillHandler(billStore *BillStore, transactionStore *transaction.TransactionStore) *BillHandler {
	return &BillHandler{billStore: billStore, transactionStore: transactionStore}
}

func (h *BillHandler) GetAllBills(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	bills, err := h.billStore.GetAll(r.Context(), userID)
	if err != nil {
		log.Printf("GetAllBills: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bills); err != nil {
		log.Printf("GetAllBills: encode failed: %v", err)
	}
}

func (h *BillHandler) CreateBill(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	var req CreateBillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	bill := &Bill{
		UserID:      userID,
		Description: req.Description,
		DueDay:      req.DueDay,
		Name:        req.Name,
		Owner:       req.Owner,
		Shared:      req.Shared,
		SourceID:    req.SourceID,
		Status:      "active",
	}

	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			http.Error(w, "invalid category_id", http.StatusBadRequest)
			return
		}
		bill.CategoryID = &catID
	}

	created, err := h.billStore.Create(r.Context(), bill)
	if err != nil {
		log.Printf("CreateBill: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateBill: encode failed: %v", err)
	}
}

func (h *BillHandler) UpdateBill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateBillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fields := bson.M{
		"description": req.Description,
		"due_day":     req.DueDay,
		"name":        req.Name,
		"owner":       req.Owner,
		"shared":      req.Shared,
		"source_id":   req.SourceID,
		"status":      req.Status,
	}

	if req.CategoryID != "" {
		catID, err := primitive.ObjectIDFromHex(req.CategoryID)
		if err != nil {
			http.Error(w, "invalid category_id", http.StatusBadRequest)
			return
		}
		fields["category_id"] = &catID
	}

	updated, err := h.billStore.Update(r.Context(), id, fields)
	if err != nil {
		log.Printf("UpdateBill: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateBill: id=%s encode failed: %v", id, err)
	}
}

func (h *BillHandler) DeleteBill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := h.billStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteBill: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BillHandler) PayBill(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, _ := token.UserIDFromContext(r.Context())

	bill, err := h.billStore.GetOne(r.Context(), id)
	if err != nil {
		log.Printf("PayBill: GetOne id=%s %v", id, err)
		http.Error(w, "bill not found", http.StatusNotFound)
		return
	}

	var data PayBillRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	date := data.Date
	billMonth := data.BillMonth
	if _, err := time.Parse("2006-01", billMonth); err != nil {
		http.Error(w, "invalid bill_month format, expected YYYY-MM", http.StatusBadRequest)
		return
	}

	amount := data.Amount

	billID := bill.ID
	sourceID := bill.SourceID
	txn := &transaction.Transaction{
		UserID:      userID,
		Amount:      amount,
		BillMonth:   billMonth,
		BillID:      &billID,
		CategoryID:  bill.CategoryID,
		Date:        date,
		Description: data.Description,
		Income:      false,
		Owner:       bill.Owner,
		Shared:      bill.Shared,
		SourceID:    &sourceID,
		Type:        "expense",
	}

	created, err := h.transactionStore.Create(r.Context(), txn)
	if err != nil {
		log.Printf("PayBill: create transaction id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("PayBill: encode failed: %v", err)
	}
}
