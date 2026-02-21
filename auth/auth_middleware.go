package auth

import (
	"context"
	"net/http"
	"time"
)

type ctxKey string

const userIDKey ctxKey = "auth.userId"

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	s, ok := v.(string)
	return s, ok
}

type AuthMiddleware struct {
	authStore     *AuthStore
	pepper        string
	cookieCfg     CookieConfig
	renewIfWithin time.Duration
}

func NewAuthMiddleware(store *AuthStore, pepper string, cookieCfg CookieConfig, renewIfWithin time.Duration) *AuthMiddleware {
	return &AuthMiddleware{
		authStore:     store,
		pepper:        pepper,
		cookieCfg:     cookieCfg,
		renewIfWithin: renewIfWithin,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SidCookieName)
		if err != nil || c.Value == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rawSID := c.Value
		sidHash := HashToken(rawSID, m.pepper)

		sess, err := m.authStore.GetByID(r.Context(), sidHash)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		now := time.Now().UTC()
		if now.After(sess.ExpiresAt) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if sess.ExpiresAt.Sub(now) <= m.renewIfWithin {
			newExp := now.Add(m.cookieCfg.TTL)

			if err := m.authStore.ExtendExpiry(r.Context(), sidHash, newExp); err == nil {
				SetSessionCookie(w, rawSID, m.cookieCfg)
			}
		}

		ctx := context.WithValue(r.Context(), userIDKey, sess.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
