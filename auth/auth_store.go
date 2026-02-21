package auth

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrSessionNotFound = errors.New("session not found")

type AuthStore struct {
	collection *mongo.Collection
}

func NewAuthStore(collection *mongo.Collection) *AuthStore {
	return &AuthStore{collection: collection}
}

func (s *AuthStore) EnsureIndexes(ctx context.Context) error {
	// TTL index on expiresAt
	_, err := s.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(0),
	})
	return err
}

// Call once on startup
func (s *AuthStore) Create(ctx context.Context, session Session) error {
	_, err := s.collection.InsertOne(ctx, session)
	return err
}

func (s *AuthStore) GetByID(ctx context.Context, sidHash string) (*Session, error) {
	var session Session
	err := s.collection.FindOne(ctx, bson.M{"_id": sidHash}).Decode(&session)
	if err == mongo.ErrNoDocuments {
		return nil, ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	// Safety: treat expired as not found
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionNotFound
	}
	return &session, nil
}

func (s *AuthStore) Delete(ctx context.Context, sidHash string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": sidHash})
	return err
}

func (s *AuthStore) ExtendExpiry(ctx context.Context, sidHash string, newExpiresAt time.Time) error {
	_, err := s.collection.UpdateByID(ctx, sidHash, bson.M{
		"$set": bson.M{"expiresAt": newExpiresAt},
	})
	return err
}
