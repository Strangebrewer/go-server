package account

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccountStore struct {
	collection *mongo.Collection
}

func NewAccountStore(collection *mongo.Collection) *AccountStore {
	return &AccountStore{collection: collection}
}

func (s *AccountStore) GetAll(ctx context.Context, userID primitive.ObjectID) ([]Account, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var accounts []Account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *AccountStore) GetOne(ctx context.Context, id string) (*Account, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var account Account
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("GetOne: account not found for id: " + id)
		}
		return nil, err
	}

	return &account, nil
}

func (s *AccountStore) Create(ctx context.Context, account *Account) (*Account, error) {
	result, err := s.collection.InsertOne(ctx, account)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		account.ID = oid
	}

	return account, nil
}

func (s *AccountStore) Update(ctx context.Context, id string, fields bson.M) (*Account, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": fields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Account
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: account not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *AccountStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: account not found for id: " + id)
	}

	return result, nil
}
