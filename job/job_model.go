package job

import "go.mongodb.org/mongo-driver/bson/primitive"

type Job struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DateApplied    string             `json:"dateApplied" bson:"dateApplied"`
	Interviews     []string           `json:"interviews" bson:"interviews"`
	PointOfContact string             `json:"pointOfContact" bson:"point_of_contact"`
	PocTitle       string             `json:"pocTitle" bson:"poc_title"`
	CompanyName    string             `json:"companyName" bson:"company_name"`
	CompanyAddress string             `json:"companyAddress" bson:"company_address"`
	CompanyCity    string             `json:"companyCity" bson:"company_city"`
	CompanyState   string             `json:"companyState" bson:"company_state"`
	Archived       bool               `json:"archived" bson:"archived"`
	JobTitle       string             `json:"jobTitle" bson:"job_title"`
	Recruiter      primitive.ObjectID `json:"recruiter" bson:"recruiter"`
	WorkFrom       string             `json:"workFrom" bson:"work_from"`
	Status         string             `json:"status" bson:"status"`
	Comments       []string           `json:"comments" bson:"comments"`
	User           primitive.ObjectID `json:"user" bson:"user"`
}
