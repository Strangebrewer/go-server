package thing

import "go.mongodb.org/mongo-driver/bson/primitive"

type Thing struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
}
