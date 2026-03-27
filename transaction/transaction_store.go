package transaction

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionStore struct {
	collection *mongo.Collection
}

func NewTransactionStore(collection *mongo.Collection) *TransactionStore {
	return &TransactionStore{collection: collection}
}

func (s *TransactionStore) GetAll(ctx context.Context, userID primitive.ObjectID, f TransactionFilter) ([]Transaction, error) {
	filter := bson.D{{Key: "user_id", Value: userID}}

	if f.Month != "" {
		filter = append(filter, bson.E{Key: "bill_month", Value: f.Month})
	}

	if f.Owner != "" {
		filter = append(filter, bson.E{Key: "owner", Value: f.Owner})
	}

	if f.Shared == "true" {
		filter = append(filter, bson.E{Key: "shared", Value: true})
	} else if f.Shared == "false" {
		filter = append(filter, bson.E{Key: "shared", Value: false})
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []Transaction
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *TransactionStore) GetOne(ctx context.Context, id string) (*Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var txn Transaction
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&txn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("GetOne: transaction not found for id %s: %w", id, mongo.ErrNoDocuments)
		}
		return nil, err
	}

	return &txn, nil
}

func (s *TransactionStore) Create(ctx context.Context, txn *Transaction) (*Transaction, error) {
	result, err := s.collection.InsertOne(ctx, txn)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		txn.ID = oid
	}

	return txn, nil
}

func (s *TransactionStore) Update(ctx context.Context, id string, fields bson.M) (*Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": fields}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Transaction
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Update: transaction not found for id: " + id)
		}
		return nil, err
	}

	return &updated, nil
}

func (s *TransactionStore) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("Delete: transaction not found for id: " + id)
	}

	return result, nil
}
