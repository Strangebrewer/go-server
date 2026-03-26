package template

import "go.mongodb.org/mongo-driver/bson/primitive"

type Template struct {
	ID            primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	UserID        primitive.ObjectID  `json:"user_id" bson:"user_id"`
	CategoryID    primitive.ObjectID  `json:"category_id" bson:"category_id"`
	Name          string              `json:"name" bson:"name"`
	Description   string              `json:"description,omitempty" bson:"description,omitempty"`
	Owner         string              `json:"owner" bson:"owner"`   // 'mine' or 'hers'
	Type          string              `json:"type" bson:"type"`     // 'expense' or 'deposit'
	Shared        bool                `json:"shared" bson:"shared"`
	SourceID      *primitive.ObjectID `json:"source_id,omitempty" bson:"source_id,omitempty"`
	DestinationID *primitive.ObjectID `json:"destination_id,omitempty" bson:"destination_id,omitempty"`
}

type CreateTemplateRequest struct {
	CategoryID    primitive.ObjectID  `json:"category_id"`
	Name          string              `json:"name"`
	Description   string              `json:"description,omitempty"`
	Owner         string              `json:"owner"`
	Type          string              `json:"type"`
	Shared        bool                `json:"shared"`
	SourceID      *primitive.ObjectID `json:"source_id,omitempty"`
	DestinationID *primitive.ObjectID `json:"destination_id,omitempty"`
}

type UpdateTemplateRequest struct {
	CategoryID    primitive.ObjectID  `json:"category_id"`
	Name          string              `json:"name"`
	Description   string              `json:"description,omitempty"`
	Owner         string              `json:"owner"`
	Type          string              `json:"type"`
	Shared        bool                `json:"shared"`
	SourceID      *primitive.ObjectID `json:"source_id,omitempty"`
	DestinationID *primitive.ObjectID `json:"destination_id,omitempty"`
}
