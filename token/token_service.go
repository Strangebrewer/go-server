package token

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrMissingBearerToken = errors.New("missing bearer token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type Service struct {
	store      *TokenStore
	priv       *rsa.PrivateKey
	pub        *rsa.PublicKey
	pepper     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(store *TokenStore, privateKeyPEM string, publicKeyPEM string, refreshPepper string) (*Service, error) {
	priv, err := parseRSAPrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	pub, err := parseRSAPublicKey(publicKeyPEM)
	if err != nil {
		return nil, err
	}

	return &Service{
		store:      store,
		priv:       priv,
		pub:        pub,
		pepper:     refreshPepper,
		accessTTL:  15 * time.Minute,
		refreshTTL: 30 * 24 * time.Hour,
	}, nil
}

type ExchangeResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Service) IssueForUser(ctx context.Context, userID primitive.ObjectID) (*ExchangeResult, error) {
	access, err := s.mintAccessJWT(userID)
	if err != nil {
		return nil, err
	}

	refreshPlain, refreshHash, err := s.mintRefreshToken()
	if err != nil {
		return nil, err
	}

	t := &Token{
		UserID:    userID,
		Type:      TokenTypeRefresh,
		Hash:      refreshHash,
		ExpiresAt: time.Now().UTC().Add(s.refreshTTL),
	}

	if err := s.store.Create(ctx, t); err != nil {
		return nil, err
	}

	return &ExchangeResult{AccessToken: access, RefreshToken: refreshPlain}, nil
}

// Exchange rotates refresh token (revoke old, mint/store new) and returns new pair.
func (s *Service) Exchange(ctx context.Context, refreshPlain string) (*ExchangeResult, error) {
	hash := s.hashRefresh(refreshPlain)

	old, err := s.store.FindActiveRefreshByHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// rotate: revoke old refresh token
	_ = s.store.RevokeByID(ctx, old.ID)

	// mint new pair
	return s.IssueForUser(ctx, old.UserID)
}

func (s *Service) Revoke(ctx context.Context, refreshPlain string) error {
	hash := s.hashRefresh(refreshPlain)

	old, err := s.store.FindActiveRefreshByHash(ctx, hash)
	if err != nil {
		return ErrInvalidToken
	}

	return s.store.RevokeByID(ctx, old.ID)
}

func (s *Service) mintAccessJWT(userID primitive.ObjectID) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub": userID.Hex(),
		"typ": "access",
		"iat": now.Unix(),
		"exp": now.Add(s.accessTTL).Unix(),
		"jti": uuid.NewString(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return tok.SignedString(s.priv)
}

func (s *Service) mintRefreshToken() (plain string, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	plain = base64.RawURLEncoding.EncodeToString(b)
	hash = s.hashRefresh(plain)
	return plain, hash, nil
}

func (s *Service) hashRefresh(refreshPlain string) string {
	sum := sha256.Sum256([]byte(s.pepper + ":" + refreshPlain))
	return hex.EncodeToString(sum[:])
}

func (s *Service) VerifyAccessJWT(jwtStr string) (primitive.ObjectID, error) {
	tok, err := jwt.Parse(jwtStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok || t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, ErrInvalidToken
		}
		return s.pub, nil
	})
	if err != nil || tok == nil || !tok.Valid {
		return primitive.NilObjectID, ErrInvalidToken
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return primitive.NilObjectID, ErrInvalidToken
	}

	// optional but good: ensure token type
	if typ, _ := claims["typ"].(string); typ != "access" {
		return primitive.NilObjectID, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return primitive.NilObjectID, ErrInvalidToken
	}

	uid, err := primitive.ObjectIDFromHex(sub)
	if err != nil {
		return primitive.NilObjectID, ErrInvalidToken
	}
	return uid, nil
}

func parseRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, ErrInvalidToken
	}

	// PKCS#1 ("RSA PRIVATE KEY")
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// PKCS#8 ("PRIVATE KEY")
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, ErrInvalidToken
	}
	key, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidToken
	}
	return key, nil
}

func parseRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, ErrInvalidToken
	}

	// PKIX ("PUBLIC KEY")
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		if pub, ok := pubAny.(*rsa.PublicKey); ok {
			return pub, nil
		}
	}

	// PKCS#1 ("RSA PUBLIC KEY") if you ever use it
	if pub, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return pub, nil
	}

	return nil, ErrInvalidToken
}
