package app

import (
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/users"
)

type Application struct {
	ThingStore *thing.ThingStore
	UserStore  *users.UserStore
}
