package job

import (
	"context"
	"errors"

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

func (s *JobStore) GetAll(ctx context.Context) ([]Job, error) {
	cursor, err := s.collection.Find(ctx, map[string]interface{}{})
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

func (s *JobStore) GetOne(ctx context.Context, id string) (*Job, error) {
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

func (s *JobStore) Update(ctx context.Context, id string, job Job) (*Job, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": job}
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
