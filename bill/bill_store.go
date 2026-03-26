package bill

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BillStore struct {
	collection *mongo.Collection
}

func NewBillStore(collection *mongo.Collection) *BillStore {
	return &BillStore{collection: collection}
}

func (s *BillStore) GetAll(ctx context.Context, userID primitive.ObjectID) ([]Bill, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bills []Bill
	if err := cursor.All(ctx, &bills); err != nil {
		return nil, err
	}

	return bills, nil
}

func (s *BillStore) GetOne(ctx context.Context, id string) (*Bill, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var bill Bill
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&bill)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("GetOne: bill not found for id: " + id)
		}
		return nil, err
	}

	return &bill, nil
}

func (s *BillStore) Create(ctx context.Context, bill *Bill) (*Bill, error) {
	result, err := s.collection.InsertOne(ctx, bill)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		bill.ID = oid
	}

	return bill, nil
}

func (s *BillStore) Update(ctx context.Context, id string, fields bson.M) (*Bill, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": fields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Bill
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: bill not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *BillStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: bill not found for id: " + id)
	}

	return result, nil
}
