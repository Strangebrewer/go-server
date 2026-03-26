package job

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	JobTitle          string             `json:"jobTitle" bson:"job_title"`
	WorkFrom          string             `json:"workFrom" bson:"work_from"`
	Recruiter         primitive.ObjectID `json:"recruiter" bson:"recruiter"`
	DateApplied       string             `json:"dateApplied" bson:"date_applied"`
	CompanyName       string             `json:"companyName" bson:"company_name"`
	CompanyAddress    string             `json:"companyAddress" bson:"company_address"`
	CompanyCity       string             `json:"companyCity" bson:"company_city"`
	CompanyState      string             `json:"companyState" bson:"company_state"`
	PointOfContact    string             `json:"pointOfContact" bson:"point_of_contact"`
	PocTitle          string             `json:"pocTitle" bson:"poc_title"`
	Interviews        []string           `json:"interviews" bson:"interviews"`
	Comments          []string           `json:"comments" bson:"comments"`
	Status            string             `json:"status" bson:"status"`
	Archived          bool               `json:"archived" bson:"archived"`
	PrimaryLink       string             `json:"primaryLink" bson:"primary_link"`
	PrimaryLinkText   string             `json:"primaryLinkText" bson:"primary_link_text"`
	SecondaryLink     string             `json:"secondaryLink" bson:"secondary_link"`
	SecondaryLinkText string             `json:"secondaryLinkText" bson:"secondary_link_text"`
	User              primitive.ObjectID `json:"user" bson:"user"`
}

type CreateJobRequest struct {
	JobTitle          string `json:"jobTitle"`
	WorkFrom          string `json:"workFrom"`
	Recruiter         string `json:"recruiter"`
	DateApplied       string `json:"dateApplied"`
	CompanyName       string `json:"companyName"`
	CompanyAddress    string `json:"companyAddress"`
	CompanyCity       string `json:"companyCity"`
	CompanyState      string `json:"companyState"`
	PointOfContact    string `json:"poc"`
	PocTitle          string `json:"pocTitle"`
	PrimaryLink       string `json:"primaryLink"`
	PrimaryLinkText   string `json:"primaryLinkText"`
	SecondaryLink     string `json:"secondaryLink"`
	SecondaryLinkText string `json:"secondaryLinkText"`
	Status            string `json:"status"`
}

type UpdateJobRequest struct {
	JobTitle          string   `json:"jobTitle"`
	WorkFrom          string   `json:"workFrom"`
	Recruiter         string   `json:"recruiter"`
	DateApplied       string   `json:"dateApplied"`
	CompanyName       string   `json:"companyName"`
	CompanyAddress    string   `json:"companyAddress"`
	CompanyCity       string   `json:"companyCity"`
	CompanyState      string   `json:"companyState"`
	PointOfContact    string   `json:"poc"`
	PocTitle          string   `json:"pocTitle"`
	Interviews        []string `json:"interviews"`
	Comments          []string `json:"comments"`
	Status            string   `json:"status"`
	Archived          bool     `json:"archived"`
	PrimaryLink       string   `json:"primaryLink"`
	PrimaryLinkText   string   `json:"primaryLinkText"`
	SecondaryLink     string   `json:"secondaryLink"`
	SecondaryLinkText string   `json:"secondaryLinkText"`
}

type JobFilter struct {
	Company         string
	Recruiter       string
	Status          string
	WorkFrom        string
	DateMin         string
	DateMax         string
	IncludeArchived bool
	IncludeDeclined bool
	SortBy          string
	SortDir         string
}

func (rBody CreateJobRequest) ToJob(recruiterId primitive.ObjectID, userId primitive.ObjectID) (Job, error) {
	return Job{
		JobTitle:          rBody.JobTitle,
		WorkFrom:          rBody.WorkFrom,
		Recruiter:         recruiterId,
		DateApplied:       rBody.DateApplied,
		CompanyName:       rBody.CompanyName,
		CompanyAddress:    rBody.CompanyAddress,
		CompanyCity:       rBody.CompanyCity,
		CompanyState:      rBody.CompanyState,
		PointOfContact:    rBody.PointOfContact,
		PocTitle:          rBody.PocTitle,
		PrimaryLink:       rBody.PrimaryLink,
		PrimaryLinkText:   rBody.PrimaryLinkText,
		SecondaryLink:     rBody.SecondaryLink,
		SecondaryLinkText: rBody.SecondaryLinkText,
		Status:            rBody.Status,
		User:              userId,
	}, nil
}
