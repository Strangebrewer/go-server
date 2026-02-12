package recipe

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(collection *mongo.Collection) *Store {
	return &Store{collection: collection}
}

func (s *Store) GetAll(ctx context.Context) ([]Recipe, error) {
	cursor, err := s.collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recipes []Recipe
	if err := cursor.All(ctx, &recipes); err != nil {
		return nil, err
	}

	return recipes, nil
}
