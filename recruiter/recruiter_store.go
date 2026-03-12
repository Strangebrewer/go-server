package recruiter

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RecruiterStore struct {
	collection *mongo.Collection
}

func NewRecruiterStore(collection *mongo.Collection) *RecruiterStore {
	return &RecruiterStore{collection: collection}
}

func (s *RecruiterStore) GetAll(ctx context.Context) ([]Recruiter, error) {
	cursor, err := s.collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recruiters []Recruiter
	if err := cursor.All(ctx, &recruiters); err != nil {
		return nil, err
	}

	return recruiters, nil
}

func (s *RecruiterStore) GetOne(ctx context.Context, id string) (*Recruiter, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var recruiter Recruiter
	err = s.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&recruiter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("FindOne: recruiter not found for id: " + id)
		}
		return nil, err
	}

	return &recruiter, nil
}

func (s *RecruiterStore) Create(ctx context.Context, recruiter *Recruiter) (*Recruiter, error) {
	result, err := s.collection.InsertOne(ctx, recruiter)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		recruiter.ID = oid
	}

	return recruiter, nil
}

func (s *RecruiterStore) Update(ctx context.Context, id string, recruiter Recruiter) (*Recruiter, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": recruiter}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Recruiter
	err = s.collection.
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: recruiter not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *RecruiterStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: recruiter not found for id: " + id)
	}

	return result, nil
}
