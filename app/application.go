package app

import (
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/users"
)

type Application struct {
	ThingStore   *thing.ThingStore
	TokenStore   *token.TokenStore
	TokenService *token.Service
	UserStore    *users.UserStore
}
