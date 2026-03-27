package bill

import (
	"github.com/Strangebrewer/go-server/transaction"
	"github.com/go-chi/chi/v5"
)

func BillRoutes(billStore *BillStore, transactionStore *transaction.TransactionStore) chi.Router {
	r := chi.NewRouter()
	h := NewBillHandler(billStore, transactionStore)

	r.Get("/", h.GetAllBills)
	r.Post("/", h.CreateBill)
	r.Put("/{id}", h.UpdateBill)
	r.Delete("/{id}", h.DeleteBill)
	r.Post("/{id}/pay", h.PayBill)

	return r
}
