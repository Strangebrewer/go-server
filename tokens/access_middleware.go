package tokens

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ctxKey int

const userIDKey ctxKey = iota

// RequireAccess validates the access JWT from `Authorization: Bearer <jwt>`
// and injects the user id into request context.
func RequireAccess(svc *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwtStr, ok := bearerToken(r)
			if !ok {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			userID, err := svc.VerifyAccessJWT(jwtStr)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (primitive.ObjectID, bool) {
	v := ctx.Value(userIDKey)
	oid, ok := v.(primitive.ObjectID)
	return oid, ok
}
