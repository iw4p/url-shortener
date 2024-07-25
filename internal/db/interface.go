package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	Init(collectionName string) (context.Context, *mongo.Collection, error)
	Disconnect(ctx context.Context) error
}
