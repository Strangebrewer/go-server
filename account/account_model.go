package account

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Balance     int32              `json:"balance" bson:"balance"`
	Description string             `json:"description" bson:"description"`
	Name        string             `json:"name" bson:"name"`
	Owner       string             `json:"owner" bson:"owner"`   // 'mine' and 'hers', default to 'mine'
	Status      string             `json:"status" bson:"status"` // 'active' or 'closed', default to 'active'
	Type        string             `json:"type" bson:"type"`     // 'asset' or 'debt', default to 'debt' since most accounts will be CC
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type CreateAccountRequest struct {
	Balance     int32  `json:"balance"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Type        string `json:"type"`
}

type UpdateAccountRequest struct {
	Balance     int32  `json:"balance"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Type        string `json:"type"`
	Status      string `json:"status"`
}
