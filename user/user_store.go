package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user disabled")
)

type UserStore struct {
	collection *mongo.Collection
}

func NewUserStore(usersCollection *mongo.Collection) *UserStore {
	return &UserStore{collection: usersCollection}
}

func (s *UserStore) EnsureIndexes(ctx context.Context) error {
	_, err := s.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_email"),
	})
	if err != nil {
		return fmt.Errorf("create users indexes: %w", err)
	}
	return nil
}

func (s *UserStore) FindByID(ctx context.Context, idHex string) (*User, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	var u User
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) Create(ctx context.Context, email, password string) (User, error) {
	now := time.Now().UTC()
	email = normalizeEmail(email)

	hash, err := HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	u := User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
		Disabled:     false,
	}

	_, err = s.collection.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return User{}, ErrEmailExists
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}

	return u, nil
}

func (s *UserStore) FindByEmail(ctx context.Context, email string) (User, error) {
	email = normalizeEmail(email)

	var u User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

// VerifyCredentials is what auth should call (replaces your stub verifyUser).
// Returns: (userIDHex, ok, err)
func (s *UserStore) VerifyCredentials(ctx context.Context, email, password string) (string, bool, error) {
	u, err := s.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find user: %w", err)
	}

	if u.Disabled {
		// Treat as "not ok" without leaking; auth can decide response wording.
		return "", false, nil
	}

	ok, verr := VerifyPassword(password, u.PasswordHash)
	if verr != nil {
		return "", false, fmt.Errorf("verify password: %w", verr)
	}
	if !ok {
		return "", false, nil
	}

	return u.ID.Hex(), true, nil
}

func normalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
