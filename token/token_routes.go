package token

import "github.com/go-chi/chi/v5"

func TokenRoutes(svc *TokenService) chi.Router {
	r := chi.NewRouter()
	h := NewTokenHandler(svc)

	// This is what mfe-utils expects by default: POST /token/exchange with refresh token in Authorization header.
	r.Post("/exchange", h.Exchange)

	// Optional but useful.
	r.Post("/revoke", h.Revoke)

	return r
}
