package main

import (
	"context"
	"log"

	"github.com/Strangebrewer/go-server/config"
	"github.com/Strangebrewer/go-server/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Narf struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
	Role string `json:"role" bson:"role"`
}

func main() {
	cfg := config.InitConfig()
	cfg.LoadEnvVariables()

	if _, err := database.InitMongoConnection(); err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"id": "abc123"}
	var narf Narf
	err := database.NarfCollection.FindOne(context.TODO(), filter).Decode(&narf)
	if err == mongo.ErrNoDocuments {
		log.Printf("no narfing for id %s", "69856ee2e52d4e39c2ed8d9a")
	}

	log.Printf("narf!! %+v", narf)
}
