package server

import (
	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/recipe"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router, application *app.Application) {
	recipeHandler := recipe.NewHandler(application.RecipeStore)

	r.Get("/recipes", recipeHandler.GetAll)
}
