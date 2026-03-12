package recruiter

import "go.mongodb.org/mongo-driver/bson/primitive"

type Recruiter struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Company  string             `json:"company" bson:"company"`
	Phone    string             `json:"phone" bson:"phone"`
	Email    string             `json:"email" bson:"email"`
	Rating   int32              `json:"rating" bson:"rating"`
	Comments []string           `json:"comments" bson:"comments"`
	Archived bool               `json:"archived" bson:"archived"`
	User     primitive.ObjectID `json:"user" bson:"user"`
}
