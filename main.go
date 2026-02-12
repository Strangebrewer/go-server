package main

import (
	"log"

	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/config"
	"github.com/Strangebrewer/go-server/database"
	"github.com/Strangebrewer/go-server/recipe"
	"github.com/Strangebrewer/go-server/server"
)

func main() {
	cfg := config.InitConfig()
	cfg.LoadEnvVariables()

	mongoConnection, err := database.InitMongoConnection(*cfg)
	if err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}

	recipeCollection := mongoConnection.Client.Database(cfg.MongoDBName).Collection("recipe")

	recipeStore := recipe.NewStore(recipeCollection)

	application := &app.Application{
		RecipeStore: recipeStore,
	}

	s := server.New("127.0.0.1:8080", application)

	log.Printf("listening on %s", s.HTTPServer.Addr)
	err = s.HTTPServer.ListenAndServe()
	log.Printf("ListenAndServe returned: %v", err)
}
