package transaction

import "github.com/go-chi/chi/v5"

func TransactionRoutes(store *TransactionStore) chi.Router {
	r := chi.NewRouter()
	h := NewTransactionHandler(store)

	r.Get("/", h.GetAllTransactions)
	r.Get("/{id}", h.GetOneTransaction)
	r.Post("/", h.CreateTransaction)
	r.Put("/{id}", h.UpdateTransaction)
	r.Delete("/{id}", h.DeleteTransaction)

	return r
}
