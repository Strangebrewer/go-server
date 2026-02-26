package app

import (
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/user"
)

type Application struct {
	ThingStore   *thing.ThingStore
	TokenStore   *token.TokenStore
	TokenService *token.TokenService
	UserStore    *user.UserStore
}
