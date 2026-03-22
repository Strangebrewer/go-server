package thing

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ThingStore struct {
	collection *mongo.Collection
}

func NewThingStore(collection *mongo.Collection) *ThingStore {
	return &ThingStore{collection: collection}
}

func (s *ThingStore) GetAll(ctx context.Context) ([]Thing, error) {
	cursor, err := s.collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var things []Thing
	if err := cursor.All(ctx, &things); err != nil {
		return nil, err
	}

	return things, nil
}

func (s *ThingStore) GetOne(ctx context.Context, id string) (*Thing, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var thing Thing
	err = s.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&thing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("FindOne: thing not found for id: " + id)
		}
		return nil, err
	}

	return &thing, nil
}

func (s *ThingStore) Create(ctx context.Context, thing *Thing) (*Thing, error) {
	result, err := s.collection.InsertOne(ctx, thing)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		thing.ID = oid
	}

	return thing, nil
}

func (s *ThingStore) Update(ctx context.Context, id string, thing Thing) (*Thing, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": thing}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Thing
	err = s.collection.
		FindOneAndUpdate(ctx, filter, update, opts).
		Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: thing not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *ThingStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: thing not found for id: " + id)
	}

	return result, nil
}
