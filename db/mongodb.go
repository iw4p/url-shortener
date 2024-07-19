package db

import (
	"context"
	"time"

	"github.com/iw4p/url-shortener/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
}

func (m *MongoDB) InitMongoDB(collectionName string) (context.Context, *mongo.Collection, error) {
	environmentVariables := utils.ReturnEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(environmentVariables.Mongo_url))
	if err != nil {
		cancel()
		return nil, nil, err
	}

	// Check the connection
	if err = client.Ping(ctx, nil); err != nil {
		cancel()
		return nil, nil, err
	}

	m.client = client

	collection := client.Database(environmentVariables.Db_name).Collection(collectionName)
	return ctx, collection, nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(ctx)
}
