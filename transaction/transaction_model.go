package transaction

import "go.mongodb.org/mongo-driver/bson/primitive"

type Transaction struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Amount        int                 `json:"amount" bson:"amount"`
	BillMonth     string              `json:"billMonth,omitempty" bson:"bill_month,omitempty"`
	BillID        *primitive.ObjectID `json:"billId,omitempty" bson:"bill_id,omitempty"`
	CategoryID    *primitive.ObjectID `json:"categoryId,omitempty" bson:"category_id,omitempty"`
	Date          string              `json:"date" bson:"date"`
	Description   string              `json:"description,omitempty" bson:"description,omitempty"`
	DestinationID *primitive.ObjectID `json:"destinationId,omitempty" bson:"destination_id,omitempty"`
	Income        bool                `json:"income" bson:"income"` // default to false
	Owner         string              `json:"owner" bson:"owner"`   // 'mine' and 'hers', default to 'mine'
	Shared        bool                `json:"shared" bson:"shared"` // default to false
	SourceID      *primitive.ObjectID `json:"sourceId,omitempty" bson:"source_id,omitempty"`
	Type          string              `json:"type" bson:"type"` // "expense" or "deposit", default to 'expense'
	UserID        primitive.ObjectID  `json:"userId" bson:"user_id"`
}

type CreateTransactionRequest struct {
	Amount        int    `json:"amount"`
	BillID        string `json:"billId,omitempty"`
	CategoryID    string `json:"categoryId,omitempty"`
	Date          string `json:"date"`
	Description   string `json:"description,omitempty"`
	DestinationID string `json:"destinationId,omitempty"`
	Income        bool   `json:"income"`
	Owner         string `json:"owner"`
	Shared        bool   `json:"shared"`
	SourceID      string `json:"sourceId,omitempty"`
	Type          string `json:"type"`
}

type UpdateTransactionRequest struct {
	Amount        int    `json:"amount"`
	BillID        string `json:"billId,omitempty"`
	CategoryID    string `json:"categoryId,omitempty"`
	Date          string `json:"date"`
	Description   string `json:"description,omitempty"`
	DestinationID string `json:"destinationId,omitempty"`
	Income        bool   `json:"income"`
	Owner         string `json:"owner"`
	Shared        bool   `json:"shared"`
	SourceID      string `json:"sourceId,omitempty"`
	Type          string `json:"type"`
}

type TransactionFilter struct {
	Month  string
	Owner  string // 'mine', 'hers', or '' for all
	Shared string // "true", "false", or "" for all
}
