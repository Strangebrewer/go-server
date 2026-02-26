package token

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenStore struct {
	collection *mongo.Collection
}

func NewTokenStore(tokensCollection *mongo.Collection) *TokenStore {
	return &TokenStore{collection: tokensCollection}
}

func (s *TokenStore) EnsureIndexes(ctx context.Context) error {
	// TTL index on expiresAt so expired tokens are cleaned up automatically
	ttl := mongo.IndexModel{
		Keys: bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(0),
	}

	// Lookup index for refresh tokens by hash (and optionally by userId)
	byHash := mongo.IndexModel{
		Keys: bson.D{{Key: "hash", Value: 1}},
		Options: options.Index().
			SetUnique(true),
	}

	_, err := s.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{ttl, byHash})
	return err
}

func (s *TokenStore) Create(ctx context.Context, t *Token) error {
	now := time.Now().UTC()
	t.CreatedAt = now

	res, err := s.collection.InsertOne(ctx, t)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		t.ID = oid
	}
	return nil
}

func (s *TokenStore) FindActiveRefreshByHash(ctx context.Context, hash string) (*Token, error) {
	var t Token
	err := s.collection.FindOne(ctx, bson.M{
		"type":      TokenTypeRefresh,
		"hash":      hash,
		"revokedAt": bson.M{"$exists": false},
		"expiresAt": bson.M{"$gt": time.Now().UTC()},
	}).Decode(&t)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TokenStore) RevokeByID(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now().UTC()
	_, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": id, "revokedAt": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"revokedAt": now}},
	)
	return err
}

func (s *TokenStore) RevokeAllForUser(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now().UTC()
	_, err := s.collection.UpdateMany(ctx,
		bson.M{"userId": userID, "revokedAt": bson.M{"$exists": false}},
		bson.M{"$set": bson.M{"revokedAt": now}},
	)
	return err
}
