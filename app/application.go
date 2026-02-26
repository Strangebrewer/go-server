package app

import (
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/tokens"
	"github.com/Strangebrewer/go-server/users"
)

type Application struct {
	ThingStore   *thing.ThingStore
	TokenStore   *tokens.TokenStore
	TokenService *tokens.Service
	UserStore    *users.UserStore
}
