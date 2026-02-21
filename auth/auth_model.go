package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type Session struct {
	ID        string    `bson:"_id"` // sidHash (hex/base64url)
	UserID    string    `bson:"userId"`
	CreatedAt time.Time `bson:"createdAt"`
	ExpiresAt time.Time `bson:"expiresAt"`
}

func NewToken(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// URL-safe, no padding (nice for cookies)
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token, pepper string) string {
	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write([]byte(token))
	sum := mac.Sum(nil)
	// base64url is shorter than hex; either is fine as long as consistent
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
