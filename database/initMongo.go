package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Strangebrewer/go-server/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

func InitMongoConnection(cfg config.Config) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s.mongodb.net/",
		cfg.MongoDBUsername,
		cfg.MongoDBPassword,
		cfg.MongoDBCluster,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &Mongo{Client: client}, nil
}
