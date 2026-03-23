package job

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobStore struct {
	collection *mongo.Collection
}

func NewJobStore(collection *mongo.Collection) *JobStore {
	return &JobStore{collection: collection}
}

func (s *JobStore) Find(ctx context.Context, userID primitive.ObjectID, f JobFilter) ([]Job, error) {
	filter := bson.D{{Key: "user", Value: userID}}

	if !f.IncludeArchived {
		filter = append(filter, bson.E{Key: "archived", Value: false})
	}

	if f.Company != "" {
		filter = append(filter, bson.E{Key: "company_name", Value: primitive.Regex{Pattern: f.Company, Options: "i"}})
	}

	if f.Recruiter != "" {
		recruiterID, err := primitive.ObjectIDFromHex(f.Recruiter)
		if err != nil {
			return nil, fmt.Errorf("GetAll: invalid recruiter id: %w", err)
		}
		filter = append(filter, bson.E{Key: "recruiter", Value: recruiterID})
	}

	if f.Status != "" {
		filter = append(filter, bson.E{Key: "status", Value: f.Status})
	}

	if f.WorkFrom != "" {
		filter = append(filter, bson.E{Key: "work_from", Value: primitive.Regex{Pattern: f.WorkFrom, Options: "i"}})
	}

	if f.DateMin != "" || f.DateMax != "" {
		dateFilter := bson.D{}
		if f.DateMin != "" {
			dateFilter = append(dateFilter, bson.E{Key: "$gte", Value: f.DateMin})
		}
		if f.DateMax != "" {
			dateFilter = append(dateFilter, bson.E{Key: "$lte", Value: f.DateMax})
		}
		filter = append(filter, bson.E{Key: "dateApplied", Value: dateFilter})
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *JobStore) FindOne(ctx context.Context, id string) (*Job, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var job Job
	err = s.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&job)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("FindOne: job not found for id: " + id)
		}
		return nil, err
	}

	return &job, nil
}

func (s *JobStore) Create(ctx context.Context, job *Job) (*Job, error) {
	result, err := s.collection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		job.ID = oid
	}

	return job, nil
}

func (s *JobStore) Update(ctx context.Context, id string, fields bson.M) (*Job, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": fields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Job
	err = s.collection.
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: job not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *JobStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: job not found for id: " + id)
	}

	return result, nil
}
