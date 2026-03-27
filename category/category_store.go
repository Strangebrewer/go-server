package category

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryStore struct {
	collection *mongo.Collection
}

func NewCategoryStore(collection *mongo.Collection) *CategoryStore {
	return &CategoryStore{collection: collection}
}

func (s *CategoryStore) GetAll(ctx context.Context, userID primitive.ObjectID) ([]Category, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryStore) Create(ctx context.Context, category *Category) (*Category, error) {
	result, err := s.collection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		category.ID = oid
	}

	return category, nil
}

func (s *CategoryStore) Update(ctx context.Context, id string, fields bson.M) (*Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": fields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Category
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: category not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *CategoryStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: category not found for id: " + id)
	}

	return result, nil
}
