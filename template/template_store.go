package template

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TemplateStore struct {
	collection *mongo.Collection
}

func NewTemplateStore(collection *mongo.Collection) *TemplateStore {
	return &TemplateStore{collection: collection}
}

func (s *TemplateStore) GetAll(ctx context.Context, userID primitive.ObjectID) ([]Template, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var templates []Template
	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

func (s *TemplateStore) GetOne(ctx context.Context, id string) (*Template, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var t Template
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("GetOne: template not found for id %s: %w", id, mongo.ErrNoDocuments)
		}
		return nil, err
	}

	return &t, nil
}

func (s *TemplateStore) Create(ctx context.Context, t *Template) (*Template, error) {
	result, err := s.collection.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid
	}

	return t, nil
}

func (s *TemplateStore) Update(ctx context.Context, id string, fields bson.M) (*Template, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": fields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Template
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: template not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *TemplateStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: template not found for id: " + id)
	}

	return result, nil
}
