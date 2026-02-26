package main

import (
	"log"

	"github.com/Strangebrewer/go-server/app"
	"github.com/Strangebrewer/go-server/config"
	"github.com/Strangebrewer/go-server/database"
	"github.com/Strangebrewer/go-server/server"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/users"
)

func main() {
	cfg := config.InitConfig()
	cfg.LoadEnvVariables()

	mongoConnection, err := database.InitMongoConnection(*cfg)
	if err != nil {
		log.Fatalf("mongo init failed: %v", err)
	}

	db := mongoConnection.Client.Database(cfg.MongoDBName)

	tokensCollection := db.Collection("tokens")
	tokenStore := token.NewStore(tokensCollection)
	tokenService, err := token.NewService(tokenStore, cfg.PrivateKeyPEM, cfg.PublicKeyPEM, cfg.RefreshTokenPepper)
	if err != nil {
		log.Fatalf("token service init failed: %v", err)
	}

	thingCollection := db.Collection("thing")
	thingStore := thing.NewStore(thingCollection)

	usersCollection := db.Collection("users")
	usersStore := users.NewStore(usersCollection)

	application := &app.Application{
		ThingStore:   thingStore,
		TokenStore:   tokenStore,
		TokenService: tokenService,
		UserStore:    usersStore,
	}

	s := server.New("localhost:8080", application)

	log.Printf("listening on %s", s.HTTPServer.Addr)
	err = s.HTTPServer.ListenAndServe()
	log.Printf("ListenAndServe returned: %v", err)
}
