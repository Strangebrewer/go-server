package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Strangebrewer/go-server/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var NarfCollection *mongo.Collection
var mongoErrMessage string

func InitMongoConnection() (string, error) {
	cfg := config.GetCurrentCfg()

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s.mongodb.net/",
		cfg.MongoDBUsername,
		cfg.MongoDBPassword,
		cfg.MongoDBCluster,
	)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		mongoErrMessage = "Error opening MongoDB connection: " + err.Error()
		return mongoErrMessage, err
	}

	NarfCollection = mongoClient.Database("derp").Collection("narf")

	return "", nil
}
