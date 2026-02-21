package auth

import (
	"net/http"
	"time"
)

const SidCookieName = "sid"

type CookieConfig struct {
	Domain   string        // usually empty for localhost
	Secure   bool          // false locally, true in prod
	SameSite http.SameSite // Lax locally, likely None in prod for cross-site SPA
	TTL      time.Duration // e.g. 7*24h
}

func SetSessionCookie(w http.ResponseWriter, rawSID string, cfg CookieConfig) {
	exp := time.Now().UTC().Add(cfg.TTL)

	http.SetCookie(w, &http.Cookie{
		Name:     SidCookieName,
		Value:    rawSID,
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Domain:   cfg.Domain,
		Expires:  exp,
		MaxAge:   int(cfg.TTL.Seconds()),
	})
}

func ClearSessionCookie(w http.ResponseWriter, cfg CookieConfig) {
	http.SetCookie(w, &http.Cookie{
		Name:     SidCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Domain:   cfg.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}
