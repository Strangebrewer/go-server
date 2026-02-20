package server

import (
	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	thingHandler := thing.NewHandler(application.ThingStore)

	r.Get("/things", thingHandler.GetAllThings)
	r.Get("/things/:id", thingHandler.GetOneThing)
	r.Post("/things", thingHandler.CreateThing)
	r.Put("/things/:id", thingHandler.UpdateThing)
	r.Delete("/things/:id", thingHandler.DeleteThing)
}
