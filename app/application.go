package app

import (
	"github.com/Strangebrewer/go-server/auth"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/users"
)

type Application struct {
	ThingStore *thing.ThingStore
	AuthStore  *auth.AuthStore
	UserStore  *users.UserStore
}
