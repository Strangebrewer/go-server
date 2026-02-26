package token

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Type      TokenType          `bson:"type" json:"type"`
	Hash      string             `bson:"hash" json:"-"` // never return
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	RevokedAt *time.Time         `bson:"revokedAt,omitempty" json:"revokedAt,omitempty"`
}
