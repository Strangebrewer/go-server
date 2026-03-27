package bill

import (
	"github.com/Strangebrewer/go-server/transaction"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bill struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	// Category id is a pointer, which allows for nil, and omitempty - this makes it an optional field
	//   this is not a general rule for optional fields, though - different types can need different
	//   setups to be optional.
	CategoryID  *primitive.ObjectID `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Description string              `json:"description" bson:"description"`
	DueDay      int                 `json:"dueDay" bson:"due_day"`
	Name        string              `json:"name" bson:"name"`
	Owner       string              `json:"owner" bson:"owner"` // 'mine' and 'hers', default to 'mine'
	Shared      bool                `json:"shared" bson:"shared"`
	SourceID    primitive.ObjectID  `json:"sourceId" bson:"source_id"`
	Status      string              `json:"status" bson:"status"` // 'active' or 'inactive', default to 'active'
	UserID      primitive.ObjectID  `json:"userId" bson:"user_id"`
}

type CreateBillRequest struct {
	// even though CategoryID is an ObjectID, it will come in from the frontend as a string
	//   it gets converted to an ObjectID in the handler
	CategoryID  string `json:"categoryId,omitempty"`
	Description string `json:"description"`
	DueDay      int    `json:"dueDay"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Shared      bool   `json:"shared"`
	SourceID    string `json:"sourceId"`
}

type UpdateBillRequest struct {
	CategoryID  string `json:"categoryId,omitempty"`
	Description string `json:"description"`
	DueDay      int    `json:"dueDay"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	SourceID    string `json:"sourceId"`
	Shared      bool   `json:"shared"`
	Status      string `json:"status"`
}

type PayBillRequest struct {
	Amount      int    `json:"amount"`
	BillMonth   string `json:"billMonth"`
	Date        string `json:"date"`
	Description string `json:"description,omitempty"`
}

type BillResponse struct {
	Bill
	Transactions []transaction.Transaction `json:"transactions"`
}
