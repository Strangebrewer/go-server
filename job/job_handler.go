package job

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Strangebrewer/go-server/recruiter"
	"github.com/Strangebrewer/go-server/token"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobHandler struct {
	jobStore       *JobStore
	recruiterStore *recruiter.RecruiterStore
}

func NewJobHandler(store *JobStore, recruiterStore *recruiter.RecruiterStore) *JobHandler {
	return &JobHandler{jobStore: store, recruiterStore: recruiterStore}
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	userID, _ := token.UserIDFromContext(r.Context())

	filter := JobFilter{
		Company:         r.URL.Query().Get("company"),
		Recruiter:       r.URL.Query().Get("recruiter"),
		Status:          r.URL.Query().Get("status"),
		WorkFrom:        r.URL.Query().Get("workFrom"),
		DateMin:         r.URL.Query().Get("dateMin"),
		DateMax:         r.URL.Query().Get("dateMax"),
		IncludeArchived: r.URL.Query().Get("archived") == "true",
		IncludeDeclined: r.URL.Query().Get("includeDeclined") == "true",
		SortBy:          r.URL.Query().Get("sortBy"),
		SortDir:         r.URL.Query().Get("sortDir"),
	}

	jobs, err := h.jobStore.Find(r.Context(), userID, filter)
	if err != nil {
		log.Printf("GetAllJobs: search failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		log.Printf("GetAllJobs: encode failed: %v", err)
	}
}

func (h *JobHandler) GetOneJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	job, err := h.jobStore.FindOne(r.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		log.Printf("GetOneJob: id=%s %v", id, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(job); err != nil {
		log.Printf("GetOneJob: id=%s encode failed: %v", id, err)
	}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateJobRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var recruiterId primitive.ObjectID
	recruiterId, err = primitive.ObjectIDFromHex(reqBody.Recruiter)
	if err != nil {
		recruiter, err := h.recruiterStore.FindByName(r.Context(), "No Recruiter")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		recruiterId = recruiter.ID
	}

	userId, _ := token.UserIDFromContext(r.Context())
	job, err := reqBody.ToJob(recruiterId, userId)

	created, err := h.jobStore.Create(r.Context(), &job)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf("CreateJob: encode failed: %v", err)
	}
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fields := bson.M{
		"job_title":           req.JobTitle,
		"work_from":           req.WorkFrom,
		"date_applied":        req.DateApplied,
		"company_name":        req.CompanyName,
		"company_address":     req.CompanyAddress,
		"company_city":        req.CompanyCity,
		"company_state":       req.CompanyState,
		"point_of_contact":    req.PointOfContact,
		"poc_title":           req.PocTitle,
		"interviews":          req.Interviews,
		"comments":            req.Comments,
		"status":              req.Status,
		"archived":            req.Archived,
		"primary_link":        req.PrimaryLink,
		"primary_link_text":   req.PrimaryLinkText,
		"secondary_link":      req.SecondaryLink,
		"secondary_link_text": req.SecondaryLinkText,
	}

	if req.Recruiter != "" {
		recruiterID, err := primitive.ObjectIDFromHex(req.Recruiter)
		if err != nil {
			http.Error(w, "invalid recruiter id", http.StatusBadRequest)
			return
		}
		fields["recruiter"] = recruiterID
	}

	updated, err := h.jobStore.Update(r.Context(), id, fields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updated == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updated); err != nil {
		log.Printf("UpdateJob: id=%s encode failed: %v", id, err)
	}
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.jobStore.Delete(r.Context(), id)
	if err != nil {
		log.Printf("DeleteJob: id=%s delete failed: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.DeletedCount == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
